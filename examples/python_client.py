#!/usr/bin/env python3
"""
Stock-A-Future API Python客户端示例
使用前请确保API服务器已启动
"""

import requests
import json
from datetime import datetime, timedelta
from typing import Dict, List, Optional


class StockAFutureClient:
    """Stock-A-Future API客户端"""
    
    def __init__(self, base_url: str = "http://localhost:8080"):
        self.base_url = base_url.rstrip('/')
        self.session = requests.Session()
        self.session.headers.update({
            'Content-Type': 'application/json',
            'User-Agent': 'StockAFuture-Python-Client/1.0'
        })
    
    def health_check(self) -> Dict:
        """健康检查"""
        response = self.session.get(f"{self.base_url}/api/v1/health")
        return response.json()
    
    def get_daily_data(self, stock_code: str, start_date: Optional[str] = None, 
                      end_date: Optional[str] = None) -> Dict:
        """
        获取股票日线数据
        
        Args:
            stock_code: 股票代码，如 '000001.SZ'
            start_date: 开始日期，格式 'YYYYMMDD'
            end_date: 结束日期，格式 'YYYYMMDD'
        """
        url = f"{self.base_url}/api/v1/stocks/{stock_code}/daily"
        params = {}
        
        if start_date:
            params['start_date'] = start_date
        if end_date:
            params['end_date'] = end_date
            
        response = self.session.get(url, params=params)
        return response.json()
    
    def get_indicators(self, stock_code: str) -> Dict:
        """获取技术指标"""
        url = f"{self.base_url}/api/v1/stocks/{stock_code}/indicators"
        response = self.session.get(url)
        return response.json()
    
    def get_predictions(self, stock_code: str) -> Dict:
        """获取买卖点预测"""
        url = f"{self.base_url}/api/v1/stocks/{stock_code}/predictions"
        response = self.session.get(url)
        return response.json()


def format_date(days_ago: int = 0) -> str:
    """格式化日期为YYYYMMDD"""
    date = datetime.now() - timedelta(days=days_ago)
    return date.strftime('%Y%m%d')


def print_json(data: Dict, title: str = ""):
    """美化打印JSON数据"""
    if title:
        print(f"\n=== {title} ===")
    print(json.dumps(data, indent=2, ensure_ascii=False))


def main():
    """主函数 - 演示API使用"""
    client = StockAFutureClient()
    
    # 测试股票代码列表
    test_stocks = ['000001.SZ', '600000.SH', '000002.SZ']
    
    print("Stock-A-Future API Python客户端示例")
    print("=" * 50)
    
    # 1. 健康检查
    try:
        health = client.health_check()
        print_json(health, "健康检查")
        
        if not health.get('success', False):
            print("❌ API服务不可用，请检查服务器状态")
            return
            
    except Exception as e:
        print(f"❌ 连接失败: {e}")
        print("请确保API服务器已启动 (make dev)")
        return
    
    # 2. 测试股票数据
    for stock_code in test_stocks[:1]:  # 只测试第一个股票，避免API限制
        print(f"\n{'='*20} 测试股票: {stock_code} {'='*20}")
        
        try:
            # 获取最近10天的日线数据
            start_date = format_date(10)
            end_date = format_date(0)
            
            daily_data = client.get_daily_data(stock_code, start_date, end_date)
            
            if daily_data.get('success'):
                data_list = daily_data['data']
                print(f"✅ 获取到 {len(data_list)} 条日线数据")
                if data_list:
                    latest = data_list[-1]
                    print(f"   最新数据: {latest['trade_date']}")
                    print(f"   收盘价: {latest['close']}")
                    print(f"   涨跌幅: {latest['pct_chg']}%")
            else:
                print(f"❌ 获取日线数据失败: {daily_data.get('error', 'Unknown error')}")
                continue
                
        except Exception as e:
            print(f"❌ 获取日线数据异常: {e}")
            continue
        
        try:
            # 获取技术指标
            indicators = client.get_indicators(stock_code)
            
            if indicators.get('success'):
                data = indicators['data']
                print("✅ 技术指标:")
                
                if data.get('macd'):
                    macd = data['macd']
                    print(f"   MACD: DIF={macd['dif']}, DEA={macd['dea']}, 信号={macd['signal']}")
                
                if data.get('rsi'):
                    rsi = data['rsi']
                    print(f"   RSI: RSI12={rsi['rsi12']}, 信号={rsi['signal']}")
                    
                if data.get('boll'):
                    boll = data['boll']
                    print(f"   布林带: 上轨={boll['upper']}, 中轨={boll['middle']}, 下轨={boll['lower']}")
                    
            else:
                print(f"❌ 获取技术指标失败: {indicators.get('error', 'Unknown error')}")
                
        except Exception as e:
            print(f"❌ 获取技术指标异常: {e}")
        
        try:
            # 获取买卖点预测
            predictions = client.get_predictions(stock_code)
            
            if predictions.get('success'):
                data = predictions['data']
                pred_list = data.get('predictions', [])
                confidence = data.get('confidence', 0)
                
                print(f"✅ 买卖点预测 (置信度: {float(confidence):.2%}):")
                
                for pred in pred_list:
                    print(f"   {pred['type']} - 价格: {pred['price']}, "
                          f"概率: {float(pred['probability']):.1%}, "
                          f"理由: {pred['reason']}")
                          
                if not pred_list:
                    print("   暂无明确的买卖信号")
                    
            else:
                print(f"❌ 获取预测失败: {predictions.get('error', 'Unknown error')}")
                
        except Exception as e:
            print(f"❌ 获取预测异常: {e}")
    
    print(f"\n{'='*50}")
    print("示例运行完成!")
    print("\n💡 提示:")
    print("- 确保已配置有效的Tushare Token")
    print("- 免费账号有API调用限制，请适度使用")
    print("- 预测结果仅供参考，不构成投资建议")


if __name__ == "__main__":
    main()
