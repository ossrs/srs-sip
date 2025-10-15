package main

import (
	"path"
	"path/filepath"
	"strings"
	"testing"
)

// TestPathTraversalPrevention 测试路径遍历防护
func TestPathTraversalPrevention(t *testing.T) {
	baseDir := "/var/www/html"

	tests := []struct {
		name       string
		inputPath  string
		shouldFail bool
		reason     string
	}{
		{
			name:       "Normal file",
			inputPath:  "/index.html",
			shouldFail: false,
			reason:     "Normal file access should be allowed",
		},
		{
			name:       "Subdirectory file",
			inputPath:  "/css/style.css",
			shouldFail: false,
			reason:     "Subdirectory access should be allowed",
		},
		{
			name:       "Deep subdirectory",
			inputPath:  "/js/lib/jquery.min.js",
			shouldFail: false,
			reason:     "Deep subdirectory access should be allowed",
		},
		{
			name:       "Parent directory traversal",
			inputPath:  "/../etc/passwd",
			shouldFail: true,
			reason:     "Parent directory traversal should be blocked",
		},
		{
			name:       "Double parent traversal",
			inputPath:  "/../../etc/passwd",
			shouldFail: true,
			reason:     "Double parent traversal should be blocked",
		},
		{
			name:       "Multiple parent traversal",
			inputPath:  "/../../../etc/passwd",
			shouldFail: true,
			reason:     "Multiple parent traversal should be blocked",
		},
		{
			name:       "Mixed path with parent",
			inputPath:  "/css/../../etc/passwd",
			shouldFail: true,
			reason:     "Mixed path with parent should be blocked",
		},
		{
			name:       "Dot slash path",
			inputPath:  "/./index.html",
			shouldFail: false,
			reason:     "Dot slash should be cleaned but allowed",
		},
		{
			name:       "Complex traversal",
			inputPath:  "/css/../js/../../../etc/passwd",
			shouldFail: true,
			reason:     "Complex traversal should be blocked",
		},
		{
			name:       "Root path",
			inputPath:  "/",
			shouldFail: false,
			reason:     "Root path should be allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 模拟修复后的路径验证逻辑
			// 首先检查原始路径是否包含 ".."
			containsDoubleDotInOriginal := strings.Contains(tt.inputPath, "..")

			// 如果原始路径包含 ".."，直接阻止
			if containsDoubleDotInOriginal {
				if !tt.shouldFail {
					t.Errorf("%s: Path contains '..' but should be allowed: %s", tt.reason, tt.inputPath)
				}
				t.Logf("Input: %s, Contains '..': true, Blocked: true (early check)", tt.inputPath)
				return
			}

			// 清理路径
			cleanPath := path.Clean(tt.inputPath)

			// 构建文件路径
			filePath := filepath.Join(baseDir, cleanPath)

			// 获取绝对路径
			absDir, err := filepath.Abs(baseDir)
			if err != nil {
				t.Fatalf("Failed to get absolute path of base dir: %v", err)
			}

			absFilePath, err := filepath.Abs(filePath)
			if err != nil {
				t.Fatalf("Failed to get absolute path of file: %v", err)
			}

			// 验证路径是否在允许的目录内
			isOutsideBaseDir := !strings.HasPrefix(absFilePath, absDir)

			// 判断是否应该被阻止
			shouldBlock := isOutsideBaseDir

			if tt.shouldFail && !shouldBlock {
				t.Errorf("%s: Expected path to be blocked, but it was allowed. Path: %s, Clean: %s, Abs: %s",
					tt.reason, tt.inputPath, cleanPath, absFilePath)
			}

			if !tt.shouldFail && shouldBlock {
				t.Errorf("%s: Expected path to be allowed, but it was blocked. Path: %s, Clean: %s, Abs: %s",
					tt.reason, tt.inputPath, cleanPath, absFilePath)
			}

			// 额外的日志信息用于调试
			t.Logf("Input: %s, Clean: %s, Outside base: %v, Blocked: %v",
				tt.inputPath, cleanPath, isOutsideBaseDir, shouldBlock)
		})
	}
}

// TestPathCleanBehavior 测试 path.Clean 的行为
func TestPathCleanBehavior(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/index.html", "/index.html"},
		{"/../etc/passwd", "/etc/passwd"},
		{"/./index.html", "/index.html"},
		{"/css/../index.html", "/index.html"},
		{"//double//slash", "/double/slash"},
		{"/trailing/slash/", "/trailing/slash"},
		{"/./././index.html", "/index.html"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := path.Clean(tt.input)
			if result != tt.expected {
				t.Errorf("path.Clean(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestAbsolutePathValidation 测试绝对路径验证
func TestAbsolutePathValidation(t *testing.T) {
	// 使用临时目录进行测试
	baseDir := t.TempDir()

	tests := []struct {
		name       string
		path       string
		shouldFail bool
	}{
		{
			name:       "File in base directory",
			path:       "index.html",
			shouldFail: false,
		},
		{
			name:       "File in subdirectory",
			path:       "css/style.css",
			shouldFail: false,
		},
		{
			name:       "Attempt to escape with parent",
			path:       "../outside.txt",
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanPath := path.Clean(tt.path)
			filePath := filepath.Join(baseDir, cleanPath)

			absDir, err := filepath.Abs(baseDir)
			if err != nil {
				t.Fatalf("Failed to get absolute path: %v", err)
			}

			absFilePath, err := filepath.Abs(filePath)
			if err != nil {
				t.Fatalf("Failed to get absolute file path: %v", err)
			}

			isOutside := !strings.HasPrefix(absFilePath, absDir)

			if tt.shouldFail && !isOutside {
				t.Errorf("Expected path to be outside base dir, but it wasn't: %s", absFilePath)
			}

			if !tt.shouldFail && isOutside {
				t.Errorf("Expected path to be inside base dir, but it wasn't: %s", absFilePath)
			}
		})
	}
}

// BenchmarkPathValidation 性能基准测试
func BenchmarkPathValidation(b *testing.B) {
	baseDir := "/var/www/html"
	testPath := "/css/style.css"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cleanPath := path.Clean(testPath)
		_ = strings.Contains(cleanPath, "..")
		filePath := filepath.Join(baseDir, cleanPath)
		absDir, _ := filepath.Abs(baseDir)
		absFilePath, _ := filepath.Abs(filePath)
		_ = strings.HasPrefix(absFilePath, absDir)
	}
}

// BenchmarkPathValidationMalicious 恶意路径的性能测试
func BenchmarkPathValidationMalicious(b *testing.B) {
	baseDir := "/var/www/html"
	testPath := "/../../../etc/passwd"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cleanPath := path.Clean(testPath)
		_ = strings.Contains(cleanPath, "..")
		filePath := filepath.Join(baseDir, cleanPath)
		absDir, _ := filepath.Abs(baseDir)
		absFilePath, _ := filepath.Abs(filePath)
		_ = strings.HasPrefix(absFilePath, absDir)
	}
}
