# Purches Backend 部署与维护指南 (基于宝塔面板)

本文档旨在为 `purches-backend` Go 应用程序在基于宝塔面板的 Linux 服务器上进行部署、更新和日常维护提供一套标准操作流程。

## 1. 首次部署

首次部署涉及环境配置、代码拉取、编译、服务配置和反向代理设置。

### 1.1 环境准备 (仅需一次)

#### 1.1.1 确认 Go 版本

本项目需要 Go 1.22 或更高版本。宝塔面板默认安装的 Go 版本可能过低。

1.  **检查当前版本**:
    ```bash
    go version
    ```
2.  **如果版本低于 1.22，请升级**:
    ```bash
    # 下载最新的 Go 安装包 (请到 Go 官网 https://go.dev/dl/ 查看最新版)
    wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz

    # 解压到 /usr/local
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz

    # 清理安装包
    rm go1.22.5.linux-amd64.tar.gz
    ```
3.  **配置环境变量**:
    为了让系统默认使用新版 Go，需要将新版 Go 的路径添加到环境变量的**最前面**。
    ```bash
    # 使用 vim 或宝塔的文件管理器编辑 ~/.bashrc 文件
    vim ~/.bashrc

    # 在文件末尾添加以下这行，确保 /usr/local/go/bin 在 $PATH 的最前面
    export PATH=/usr/local/go/bin:$PATH

    # 保存文件后，使其立即生效
    source ~/.bashrc
    ```
4.  **验证新版本**:
    关闭并重新打开终端，或执行 `source` 命令后，再次检查版本，确保输出为新版本。
    ```bash
    go version
    ```

### 1.2 拉取代码与编译

1.  **克隆代码**:
    如果服务器上还没有项目代码，请从 Git 仓库克隆。
    ```bash
    # 找一个合适的目录，例如 /root
    cd /root
    git clone <你的项目Git仓库地址> purches-backend
    cd purches-backend
    ```
2.  **编译应用**:
    ```bash
    # 自动下载依赖
    go mod tidy

    # 编译生成可执行文件，-o 指定输出文件名为 purches-backend
    go build -o purches-backend main.go
    ```
    成功后，目录下会多出一个名为 `purches-backend` 的二进制文件。

### 1.3 配置 Systemd 服务 (后台持久化运行)

为了让应用能在后台稳定运行，并且在服务器重启后能自动启动，我们将其配置为一个 `systemd` 服务。

1.  **创建服务文件**:
    执行以下命令，会自动创建并写入服务配置。
    ```bash
    sudo tee /etc/systemd/system/purches-backend.service > /dev/null <<EOF
    [Unit]
    Description=Purches Backend Service
    After=network.target

    [Service]
    Type=simple
    User=root
    Group=root
    WorkingDirectory=/root/purches-backend
    ExecStart=/root/purches-backend/purches-backend
    Restart=on-failure
    RestartSec=5s

    [Install]
    WantedBy=multi-user.target
    EOF
    ```
    > **注意**: `WorkingDirectory` 和 `ExecStart` 中的路径 `/root/purches-backend` 必须与您项目代码的实际存放路径一致。

2.  **启动并设置开机自启**:
    ```bash
    # 重载 systemd 配置
    sudo systemctl daemon-reload

    # 启用服务 (设置开机自启)
    sudo systemctl enable purches-backend

    # 立即启动服务
    sudo systemctl start purches-backend
    ```

### 1.4 配置 Nginx 反向代理 (通过宝塔面板)

1.  登录宝塔面板，进入“网站”菜单。
2.  确保您的域名 (例如 `www.ency.asia`) 已经添加为网站。
3.  点击该网站右侧的“设置”。
4.  在设置窗口中，选择“反向代理” -> “添加反向代理”。
5.  配置如下：
    *   **代理名称**: `purches-backend` (或任意自定义名称)
    *   **目标URL**: `http://127.0.0.1:8080`
6.  点击“提交”。

### 1.5 验证

此时，访问 `http://<你的域名>/v1/products` (例如 `http://www.ency.asia/v1/products`)，应该能看到 JSON 数据正常返回。首次部署完成！

## 2. 后续更新

当本地代码有更新并推送到 Git 仓库后，在服务器上执行以下步骤即可完成更新。

### 2.1 登录服务器并进入项目目录

```bash
cd /root/purches-backend
```

### 2.2 拉取最新代码

```bash
# 强制覆盖本地所有修改，与远程仓库保持一致
git fetch --all
git reset --hard origin/main
git pull origin main
```
> **警告**: `git reset --hard` 会丢弃服务器上所有未提交的修改。请确保服务器上的任何改动（如配置文件）都已备份或无关紧要。

### 2.3 重新编译并重启服务

```bash
# 重新编译
go mod tidy
go build -o purches-backend main.go

# 重启服务，让新版本生效
sudo systemctl restart purches-backend
```

### 2.4 检查状态

```bash
# 查看服务状态，确保 Active: active (running) 并且没有错误日志
sudo systemctl status purches-backend
```

## 3. 日常维护与故障排查

### 3.1 查看服务状态

```bash
sudo systemctl status purches-backend
```

### 3.2 查看实时日志

```bash
# 查看最新的 100 条日志并持续跟踪新日志
sudo journalctl -u purches-backend -f -n 100
```

### 3.3 启动、停止、重启服务

```bash
sudo systemctl start purches-backend
sudo systemctl stop purches-backend
sudo systemctl restart purches-backend
```

### 3.4 数据库重置

在开发过程中，如果 `models` 结构发生了重大变化，可能导致 `AutoMigrate` 失败（就像我们遇到的 `unbalanced brackets` 错误）。最简单的解决方法是重置数据库。

> **警告**: 此操作会**删除所有数据**，请仅在开发或测试环境中使用！

```bash
# 1. 停止服务
sudo systemctl stop purches-backend

# 2. 进入项目目录
cd /root/purches-backend

# 3. 删除数据库文件 (默认为 purches.db)
rm -f purches.db

# 4. 重新启动服务，它会自动创建新库并导入初始数据
sudo systemctl start purches-backend
``` 