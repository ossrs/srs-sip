# GitHub Actions CI/CD 设置指南

本文档说明如何为 SRS-SIP 项目设置和使用 GitHub Actions CI/CD 流水线。

## 📋 目录

- [快速开始](#快速开始)
- [工作流说明](#工作流说明)
- [本地测试](#本地测试)
- [发布流程](#发布流程)
- [故障排查](#故障排查)

## 🚀 快速开始

### 1. 推送代码触发 CI

当你推送代码到 `main` 或 `develop` 分支，或创建 Pull Request 时，CI 会自动运行：

```bash
git add .
git commit -m "feat: add new feature"
git push origin main
```

### 2. 查看 CI 状态

访问 GitHub 仓库的 **Actions** 标签页查看工作流运行状态。

### 3. 在 README 中显示状态徽章

已在 README.md 中添加了以下徽章：

- ✅ CI 构建状态
- 🐳 Docker 构建状态
- 🔒 CodeQL 安全扫描状态

## 📊 工作流说明

### CI 工作流 (ci.yml)

**触发条件**:
- 推送到 `main` 或 `develop` 分支
- 创建或更新 Pull Request

**包含的任务**:

1. **Backend Build & Test** (多平台)
   - Ubuntu, macOS, Windows
   - Go 1.23
   - 运行单元测试
   - 生成代码覆盖率报告
   - 上传到 Codecov

2. **Frontend Build**
   - 构建 Vue.js 前端
   - 生成静态资源

3. **Code Quality**
   - golangci-lint 代码检查
   - 检查代码规范和潜在问题

4. **Full Build**
   - 完整打包构建
   - 生成发布包

### Release 工作流 (release.yml)

**触发条件**:
- 推送版本标签 (例如 `v1.0.0`)

**功能**:
- 自动创建 GitHub Release
- 构建多平台二进制文件
- 上传发布资产

**使用方法**:

```bash
# 1. 创建标签
git tag -a v1.0.0 -m "Release version 1.0.0"

# 2. 推送标签
git push origin v1.0.0

# 3. GitHub Actions 会自动：
#    - 创建 Release
#    - 构建所有平台的二进制文件
#    - 上传到 Release 页面
```

**支持的平台**:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### Docker 工作流 (docker.yml)

**触发条件**:
- 推送到 `main` 分支
- 推送版本标签
- Pull Request

**功能**:
- 构建 Docker 镜像
- 推送到 GitHub Container Registry (ghcr.io)
- 支持多架构 (amd64, arm64)

**镜像标签**:
```
ghcr.io/ossrs/srs-sip:main          # main 分支最新版本
ghcr.io/ossrs/srs-sip:v1.0.0        # 版本标签
ghcr.io/ossrs/srs-sip:sha-abc123    # Git commit SHA
```

**使用镜像**:
```bash
docker pull ghcr.io/ossrs/srs-sip:main
docker run -d -p 5060:5060 ghcr.io/ossrs/srs-sip:main
```

### CodeQL 安全扫描 (codeql.yml)

**触发条件**:
- 推送到 `main` 或 `develop` 分支
- Pull Request
- 每周一定时运行

**功能**:
- 扫描 Go 和 JavaScript 代码
- 检测安全漏洞
- 生成安全报告

## 🧪 本地测试

在推送代码前，建议先在本地运行测试：

### 运行所有测试

```bash
make test
```

### 生成覆盖率报告

```bash
make test-coverage
# 会生成 coverage.html，用浏览器打开查看
```

### 代码质量检查

```bash
# 安装 golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 运行检查
make lint

# 自动修复部分问题
make lint-fix
```

### 格式化代码

```bash
make fmt
```

### 完整 CI 流程

```bash
# 运行与 CI 相同的流程
make ci
```

## 📦 发布流程

### 1. 准备发布

```bash
# 确保所有测试通过
make test

# 确保代码质量检查通过
make lint

# 更新版本号和 CHANGELOG
```

### 2. 创建标签

```bash
# 创建带注释的标签
git tag -a v1.0.0 -m "Release version 1.0.0

主要更新:
- 新功能 A
- Bug 修复 B
- 性能优化 C
"

# 推送标签
git push origin v1.0.0
```

### 3. 等待 CI 完成

- GitHub Actions 会自动运行 Release 工作流
- 构建所有平台的二进制文件
- 创建 GitHub Release
- 上传发布资产

### 4. 编辑 Release 说明

访问 GitHub Release 页面，编辑自动生成的 Release，添加详细的更新说明。

## 🔧 配置文件说明

### .golangci.yml

golangci-lint 配置文件，定义了代码质量检查规则。

**启用的检查器**:
- errcheck: 检查未处理的错误
- gosimple: 简化代码建议
- govet: 检查可疑的构造
- staticcheck: 静态分析
- unused: 检查未使用的代码
- gofmt: 代码格式检查
- revive: 代码规范检查

### .github/dependabot.yml

Dependabot 配置，自动检查依赖更新。

**监控的依赖**:
- Go modules
- npm packages
- GitHub Actions
- Docker

**更新频率**: 每周一

## 🐛 故障排查

### CI 构建失败

1. **查看详细日志**
   - 访问 Actions 标签页
   - 点击失败的工作流
   - 查看具体失败的步骤

2. **本地复现**
   ```bash
   # 运行相同的命令
   make ci
   ```

3. **常见问题**
   - 依赖问题: `go mod tidy`
   - 格式问题: `make fmt`
   - 测试失败: `make test`

### Docker 推送失败

1. **检查权限**
   - 确保仓库启用了 GitHub Packages
   - 检查 GITHUB_TOKEN 权限

2. **本地测试**
   ```bash
   docker build -t srs-sip:test .
   docker run -it srs-sip:test
   ```

### 测试失败

1. **本地运行测试**
   ```bash
   go test -v ./pkg/db/
   ```

2. **查看详细错误**
   ```bash
   go test -v -race ./...
   ```

3. **清理缓存**
   ```bash
   go clean -testcache
   ```

## 📈 最佳实践

### 1. 提交前检查

```bash
# 运行完整的 CI 流程
make ci

# 或者分步运行
make fmt      # 格式化代码
make lint     # 代码检查
make test     # 运行测试
make build    # 构建
```

### 2. 编写好的提交信息

使用约定式提交 (Conventional Commits):

```
feat: 添加新功能
fix: 修复 bug
docs: 更新文档
style: 代码格式调整
refactor: 重构代码
test: 添加测试
chore: 构建/工具相关
```

### 3. 创建有意义的 PR

- 使用 PR 模板
- 填写完整的描述
- 关联相关 Issue
- 确保 CI 通过

### 4. 定期更新依赖

- 关注 Dependabot 的 PR
- 及时更新安全补丁
- 测试依赖更新的影响

## 🔗 相关链接

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [golangci-lint 文档](https://golangci-lint.run/)
- [Codecov 文档](https://docs.codecov.com/)
- [Dependabot 文档](https://docs.github.com/en/code-security/dependabot)

## 💡 提示

- 所有工作流配置在 `.github/workflows/` 目录
- 可以在 Actions 标签页手动触发工作流
- 使用 GitHub Secrets 存储敏感信息
- 工作流日志保留 90 天

