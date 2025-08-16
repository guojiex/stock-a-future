package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Config 存储命令行配置
type Config struct {
	Method          string
	Headers         map[string]string
	Data            string
	User            string
	Password        string
	Timeout         time.Duration
	Insecure        bool
	Verbose         bool
	FollowRedirects bool
	MaxRedirects    int
}

// Response 存储HTTP响应信息
type Response struct {
	StatusCode   int
	Status       string
	Headers      http.Header
	Body         []byte
	ResponseTime time.Duration
}

func main() {
	config := parseFlags()

	if len(flag.Args()) == 0 {
		fmt.Println("用法: curl [选项] <URL>")
		fmt.Println("选项:")
		fmt.Println("  -X, --request <method>    指定HTTP方法 (默认: GET)")
		fmt.Println("  -H, --header <header>     添加请求头 (格式: Key:Value)")
		fmt.Println("  -d, --data <data>         发送POST数据")
		fmt.Println("  -u, --user <user:pass>    基本认证")
		fmt.Println("  --timeout <seconds>       请求超时时间 (默认: 30秒)")
		fmt.Println("  -k, --insecure           跳过SSL证书验证")
		fmt.Println("  -v, --verbose            详细输出")
		fmt.Println("  -L, --location           跟随重定向")
		fmt.Println("  --max-redirects <num>    最大重定向次数 (默认: 10)")
		fmt.Println("  -h, --help               显示帮助信息")
		os.Exit(1)
	}

	url := flag.Args()[0]

	// 创建HTTP客户端
	client := createHTTPClient(config)

	// 创建请求
	req, err := createRequest(url, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "创建请求失败: %v\n", err)
		os.Exit(1)
	}

	// 发送请求
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "请求失败: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取响应体失败: %v\n", err)
		os.Exit(1)
	}

	responseTime := time.Since(start)
	response := &Response{
		StatusCode:   resp.StatusCode,
		Status:       resp.Status,
		Headers:      resp.Header,
		Body:         body,
		ResponseTime: responseTime,
	}

	// 输出结果
	printResponse(response, config)
}

// parseFlags 解析命令行参数
func parseFlags() *Config {
	config := &Config{
		Method:          "GET",
		Headers:         make(map[string]string),
		Timeout:         30 * time.Second,
		Insecure:        false,
		Verbose:         false,
		FollowRedirects: false,
		MaxRedirects:    10,
	}

	var method, data, user, timeout string
	var insecure, verbose, location bool
	var maxRedirects int

	flag.StringVar(&method, "X", "GET", "HTTP方法")
	flag.StringVar(&method, "request", "GET", "HTTP方法")
	flag.StringVar(&data, "d", "", "POST数据")
	flag.StringVar(&data, "data", "", "POST数据")
	flag.StringVar(&user, "u", "", "基本认证 (user:pass)")
	flag.StringVar(&user, "user", "", "基本认证 (user:pass)")
	flag.StringVar(&timeout, "timeout", "30", "超时时间(秒)")
	flag.BoolVar(&insecure, "k", false, "跳过SSL验证")
	flag.BoolVar(&insecure, "insecure", false, "跳过SSL验证")
	flag.BoolVar(&verbose, "v", false, "详细输出")
	flag.BoolVar(&verbose, "verbose", false, "详细输出")
	flag.BoolVar(&location, "L", false, "跟随重定向")
	flag.BoolVar(&location, "location", false, "跟随重定向")
	flag.IntVar(&maxRedirects, "max-redirects", 10, "最大重定向次数")

	// 自定义头部解析
	var headers string
	flag.StringVar(&headers, "H", "", "请求头 (格式: Key:Value)")
	flag.StringVar(&headers, "header", "", "请求头 (格式: Key:Value)")

	flag.Parse()

	// 设置配置
	config.Method = strings.ToUpper(method)
	config.Data = data
	config.Insecure = insecure
	config.Verbose = verbose
	config.FollowRedirects = location
	config.MaxRedirects = maxRedirects

	// 解析超时
	if timeout != "" {
		if t, err := time.ParseDuration(timeout + "s"); err == nil {
			config.Timeout = t
		}
	}

	// 解析用户认证
	if user != "" {
		parts := strings.SplitN(user, ":", 2)
		if len(parts) == 2 {
			config.User = parts[0]
			config.Password = parts[1]
		}
	}

	// 解析请求头
	if headers != "" {
		headerPairs := strings.Split(headers, ",")
		for _, pair := range headerPairs {
			parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				config.Headers[key] = value
			}
		}
	}

	return config
}

// createHTTPClient 创建HTTP客户端
func createHTTPClient(config *Config) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.Insecure,
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}

	// 处理重定向
	if !config.FollowRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) >= config.MaxRedirects {
				return fmt.Errorf("重定向次数超过限制: %d", config.MaxRedirects)
			}
			return nil
		}
	}

	return client
}

// createRequest 创建HTTP请求
func createRequest(urlStr string, config *Config) (*http.Request, error) {
	var req *http.Request
	var err error

	if config.Data != "" && (config.Method == "POST" || config.Method == "PUT" || config.Method == "PATCH") {
		// 检查是否为JSON数据
		if strings.TrimSpace(config.Data)[0] == '{' {
			// JSON数据
			req, err = http.NewRequest(config.Method, urlStr, bytes.NewBufferString(config.Data))
			if err == nil {
				req.Header.Set("Content-Type", "application/json")
			}
		} else {
			// 表单数据
			req, err = http.NewRequest(config.Method, urlStr, bytes.NewBufferString(config.Data))
			if err == nil {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
		}
	} else {
		req, err = http.NewRequest(config.Method, urlStr, nil)
	}

	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// 设置基本认证
	if config.User != "" {
		req.SetBasicAuth(config.User, config.Password)
	}

	// 设置User-Agent
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Go-Curl/1.0")
	}

	return req, nil
}

// printResponse 输出响应结果
func printResponse(response *Response, config *Config) {
	if config.Verbose {
		fmt.Printf("响应时间: %v\n", response.ResponseTime)
		fmt.Printf("状态码: %d %s\n", response.StatusCode, response.Status)
		fmt.Println("\n响应头:")
		for key, values := range response.Headers {
			for _, value := range values {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}
		fmt.Println()
	}

	// 输出响应体
	if len(response.Body) > 0 {
		contentType := response.Headers.Get("Content-Type")

		if strings.Contains(contentType, "application/json") {
			// 美化JSON输出
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, response.Body, "", "  "); err == nil {
				fmt.Println(prettyJSON.String())
			} else {
				fmt.Println(string(response.Body))
			}
		} else if strings.Contains(contentType, "text/html") {
			// HTML内容，只显示前500字符
			bodyStr := string(response.Body)
			if len(bodyStr) > 500 {
				fmt.Printf("%s...\n(HTML内容已截断，总长度: %d 字符)", bodyStr[:500], len(bodyStr))
			} else {
				fmt.Println(bodyStr)
			}
		} else {
			// 其他类型内容
			fmt.Println(string(response.Body))
		}
	}
}
