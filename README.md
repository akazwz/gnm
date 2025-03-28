# GNM - Go Node Manager

GNM 是一个用 Go 语言编写的简单、快速的 Node.js 版本管理工具，类似于 fnm、nvm 和 n。

## 特点

- **速度快**：使用 Go 语言编写，启动和执行速度比其他版本管理器更快
- **跨平台**：支持 macOS、Linux 和 Windows
- **简单易用**：命令简洁，易于记忆和使用
- **多版本管理**：轻松安装、切换和管理多个 Node.js 版本

## 安装

确保你已经安装了 Go（1.16 或更高版本）：

```bash
# 克隆仓库
git clone https://github.com/akazwz/gnm.git
cd gnm

# 构建并安装
go build -o gnm ./cmd/gnm
sudo mv gnm /usr/local/bin/
```

## 使用方法

### 安装 Node.js 版本

```bash
# 安装特定版本
gnm install 18.12.0

# 安装最新的 LTS 版本
gnm install lts
```

### 切换 Node.js 版本

```bash
# 切换到特定版本
gnm use 18.12.0

# 切换到已安装的最新 LTS 版本
gnm use lts
```

### 列出已安装的版本

```bash
gnm list
# 或
gnm ls
```

### 列出可用的远程版本

```bash
# 列出所有可用版本
gnm ls-remote

# 只列出 LTS 版本
gnm ls-remote --lts
gnm ls-remote -l

# 列出所有版本（不限制数量）
gnm ls-remote --all
gnm ls-remote -a
```

### 卸载 Node.js 版本

```bash
gnm uninstall 18.12.0
# 或
gnm remove 18.12.0
```

## 配置

在使用 GNM 管理的 Node.js 版本时，请将 GNM 的 bin 目录添加到你的 PATH 中：

```bash
# 在 ~/.bashrc, ~/.zshrc 或其他 shell 配置文件中添加
export PATH="$HOME/.gnm/bin:$PATH"
```

## 贡献

欢迎提交 issues 和 pull requests！

## 许可证

MIT 