# Python 虚拟环境管理指南

本指南介绍如何在 Stock-A-Future 项目中设置和管理 Python 虚拟环境，以及如何安装和运行 AKTools 服务。

## 为什么使用虚拟环境？

Python 虚拟环境可以为每个项目创建独立的 Python 环境，避免全局包冲突，确保项目依赖的一致性和可移植性。

## 目录

1. [环境要求](#环境要求)
2. [Linux/macOS 使用指南](#linuxmacos-使用指南)
3. [Windows 使用指南](#windows-使用指南)
4. [AKTools 使用说明](#aktools-使用说明)
5. [常见问题](#常见问题)

## 环境要求

- Python 3.7 或更高版本
- pip（Python 包管理器）
- 对于 Windows：PowerShell 或命令提示符
- 对于 Linux/macOS：Bash 或兼容的 shell

## Linux/macOS 使用指南

### 创建和激活虚拟环境

```bash
# 添加执行权限
chmod +x scripts/venv/setup_venv.sh
chmod +x scripts/venv/deactivate_venv.sh
chmod +x scripts/venv/install_aktools.sh
chmod +x scripts/venv/start_aktools.sh
chmod +x scripts/venv/test_aktools_connection.sh

# 创建并激活虚拟环境
./scripts/venv/setup_venv.sh
```

### 停用虚拟环境

```bash
./scripts/venv/deactivate_venv.sh
# 或直接使用
deactivate
```

### 安装 AKTools

```bash
./scripts/venv/install_aktools.sh
```

### 启动 AKTools 服务

```bash
# 默认端口 8080
./scripts/venv/start_aktools.sh

# 指定端口
./scripts/venv/start_aktools.sh 8081
```

### 测试 AKTools 连接

```bash
# 测试默认地址 http://127.0.0.1:8080
./scripts/venv/test_aktools_connection.sh

# 测试指定地址
./scripts/venv/test_aktools_connection.sh http://127.0.0.1:8081
```

## Windows 使用指南

### 命令提示符 (CMD)

#### 创建和激活虚拟环境

```cmd
scripts\venv\setup_venv.bat
```

#### 停用虚拟环境

```cmd
scripts\venv\deactivate_venv.bat
```

#### 安装 AKTools

```cmd
scripts\venv\install_aktools.bat
```

#### 启动 AKTools 服务

```cmd
REM 默认端口 8080
scripts\venv\start_aktools.bat

REM 指定端口
scripts\venv\start_aktools.bat 8081
```

#### 测试 AKTools 连接

```cmd
REM 测试默认地址
scripts\venv\test_aktools_connection.bat

REM 测试指定地址
scripts\venv\test_aktools_connection.bat http://127.0.0.1:8081
```

### PowerShell

#### 创建和激活虚拟环境

```powershell
.\scripts\venv\setup_venv.ps1
```

#### 停用虚拟环境

```powershell
.\scripts\venv\deactivate_venv.ps1
```

#### 安装 AKTools

```powershell
.\scripts\venv\install_aktools.ps1
```

#### 启动 AKTools 服务

```powershell
# 默认端口 8080
.\scripts\venv\start_aktools.ps1

# 指定端口
.\scripts\venv\start_aktools.ps1 8081
```

#### 测试 AKTools 连接

```powershell
# 测试默认地址
.\scripts\venv\test_aktools_connection.ps1

# 测试指定地址
.\scripts\venv\test_aktools_connection.ps1 http://127.0.0.1:8081
```

## AKTools 使用说明

AKTools 是一个基于 AKShare 的免费开源财经数据接口，提供 A 股、港股、美股等多种市场数据。

### 配置 Stock-A-Future 使用 AKTools

1. 确保 AKTools 服务已启动并正常运行
2. 编辑项目根目录的 `config.env` 文件：

```
# 数据源配置 - 使用AKTools
DATA_SOURCE_TYPE=aktools

# AKTools配置 - 运行在本地8080端口
AKTOOLS_BASE_URL=http://127.0.0.1:8080

# 服务器配置 - 使用不同的端口避免冲突
SERVER_PORT=8081
SERVER_HOST=localhost
```

3. 启动 Stock-A-Future 服务器：

```bash
# Linux/macOS
go run cmd/server/main.go

# Windows
go run cmd\server\main.go
```

### 支持的数据接口

- 股票日线数据
- 股票基本信息
- 股票列表
- 更多接口详情请参考 AKTools 文档

## 常见问题

### 1. 找不到 Python 命令

确保 Python 已安装并添加到系统 PATH 环境变量中。

### 2. 虚拟环境创建失败

- 检查 Python 版本是否为 3.7+
- 确保有足够的磁盘空间
- 尝试手动运行 `python -m venv .venv` 命令查看详细错误

### 3. AKTools 安装失败

- 确保虚拟环境已激活
- 检查网络连接
- 尝试更新 pip：`pip install --upgrade pip`

### 4. AKTools 服务启动失败

- 检查端口是否被占用
- 确保 AKTools 已正确安装
- 查看详细错误信息

### 5. 连接测试失败

- 确保 AKTools 服务正在运行
- 验证服务地址和端口是否正确
- 检查网络连接和防火墙设置