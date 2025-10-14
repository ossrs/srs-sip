# GitHub Actions 工作流说明

本项目使用 GitHub Actions 实现自动化 CI/CD 流程。

## 工作流列表

### 1. CI (ci.yml)
**触发条件**: 推送到 main/develop 分支或创建 PR

**功能**:
- ✅ 多平台构建测试 (Linux, macOS, Windows)
- ✅ 运行单元测试并生成覆盖率报告
- ✅ 代码质量检查 (golangci-lint)
- ✅ 前端构建 (Vue.js)
- ✅ 完整打包构建

**任务**:
- `backend`: Go 后端构建和测试
- `frontend`: Vue 前端构建
- `lint`: 代码质量检查
- `full-build`: 完整构建并打包

### 2. Release (release.yml)
**触发条件**: 推送版本标签 (例如 v1.0.0)

**功能**:
- 🚀 自动创建 GitHub Release
- 📦 构建多平台二进制文件
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64)
- 📤 上传发布资产到 Release

**使用方法**:
```bash
# 创建并推送标签
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

### 3. Docker Build (docker.yml)
**触发条件**: 推送到 main 分支、推送标签或创建 PR

**功能**:
- 🐳 构建 Docker 镜像
- 📤 推送到 GitHub Container Registry (ghcr.io)
- 🏗️ 支持多架构 (amd64, arm64)
- 💾 使用 GitHub Actions 缓存加速构建

**镜像地址**:
```
ghcr.io/ossrs/srs-sip:main
ghcr.io/ossrs/srs-sip:v1.0.0
ghcr.io/ossrs/srs-sip:sha-abc123
```

### 4. CodeQL Security Scan (codeql.yml)
**触发条件**: 推送到 main/develop 分支、PR 或每周一定时运行

**功能**:
- 🔒 安全漏洞扫描
- 🔍 代码质量分析
- 📊 生成安全报告

## 状态徽章

在项目 README.md 中添加以下徽章来显示构建状态：

```markdown
[![CI](https://github.com/ossrs/srs-sip/actions/workflows/ci.yml/badge.svg)](https://github.com/ossrs/srs-sip/actions/workflows/ci.yml)
[![Docker Build](https://github.com/ossrs/srs-sip/actions/workflows/docker.yml/badge.svg)](https://github.com/ossrs/srs-sip/actions/workflows/docker.yml)
[![CodeQL](https://github.com/ossrs/srs-sip/actions/workflows/codeql.yml/badge.svg)](https://github.com/ossrs/srs-sip/actions/workflows/codeql.yml)
[![codecov](https://codecov.io/gh/ossrs/srs-sip/branch/main/graph/badge.svg)](https://codecov.io/gh/ossrs/srs-sip)
```

## 本地测试

### 运行测试
```bash
# 运行所有测试
go test -v ./...

# 运行测试并生成覆盖率报告
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# 查看覆盖率
go tool cover -html=coverage.out
```

### 代码质量检查
```bash
# 安装 golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 运行 lint
golangci-lint run

# 自动修复部分问题
golangci-lint run --fix
```

### 构建
```bash
# 后端构建
make build

# 前端构建
cd html/NextGB
npm install
npm run build

# 完整构建
make all
```

## 配置文件

- `.golangci.yml`: golangci-lint 配置
- `.github/workflows/*.yml`: GitHub Actions 工作流定义

## 注意事项

1. **首次运行**: 首次推送代码后，工作流会自动运行
2. **权限**: Docker 推送需要 `packages: write` 权限
3. **密钥**: Release 和 Docker 推送使用 `GITHUB_TOKEN`，无需额外配置
4. **缓存**: 使用 Go modules 和 npm 缓存加速构建
5. **并行**: 多个任务并行运行，提高效率

## 故障排查

### 构建失败
1. 检查 Actions 标签页查看详细日志
2. 本地运行相同的命令进行调试
3. 确保所有依赖都在 go.mod 和 package.json 中

### Docker 推送失败
1. 确保仓库设置中启用了 GitHub Packages
2. 检查 GITHUB_TOKEN 权限

### 测试失败
1. 本地运行 `go test -v ./...` 复现问题
2. 检查测试日志中的错误信息
3. 确保测试环境一致

## 扩展

可以根据需要添加更多工作流：
- 性能测试
- 集成测试
- 自动部署
- 依赖更新检查
- 文档生成

