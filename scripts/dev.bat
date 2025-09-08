@echo off
chcp 65001 >nul

echo 🚀 启动本地开发环境...

REM 切换到项目根目录
cd /d "%~dp0\.."

REM 检查Go环境
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go 环境未安装，请先安装 Go
    echo 下载地址: https://golang.org/dl/
    pause
    exit /b 1
)

REM 检查项目文件
if not exist "main.go" (
    echo ❌ 找不到main.go文件
    pause
    exit /b 1
)

echo 📚 安装依赖...
go mod tidy

echo.
echo ✅ 启动开发服务器 (http://localhost:8080)
echo 按 Ctrl+C 停止服务
echo ==================================
echo.

REM 启动并实时显示日志
go run main.go 