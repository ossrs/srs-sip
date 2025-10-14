package db

import (
	"os"
	"testing"
)

func TestAddOrUpdateMediaServer(t *testing.T) {
	// 创建临时数据库文件
	dbPath := "./test_media_servers.db"
	defer os.Remove(dbPath)

	// 创建数据库实例
	db, err := NewMediaServerDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// 测试第一次添加
	err = db.AddOrUpdateMediaServer("Default", "SRS", "192.168.1.100", 1985, "", "", "", 1)
	if err != nil {
		t.Fatalf("Failed to add media server: %v", err)
	}

	// 验证添加成功
	servers, err := db.ListMediaServers()
	if err != nil {
		t.Fatalf("Failed to list media servers: %v", err)
	}
	if len(servers) != 1 {
		t.Fatalf("Expected 1 server, got %d", len(servers))
	}
	if servers[0].Name != "Default" || servers[0].IP != "192.168.1.100" {
		t.Fatalf("Server data mismatch: %+v", servers[0])
	}

	// 测试重复添加（应该更新而不是插入新记录）
	err = db.AddOrUpdateMediaServer("Default", "SRS", "192.168.1.100", 1985, "admin", "password", "secret", 1)
	if err != nil {
		t.Fatalf("Failed to update media server: %v", err)
	}

	// 验证没有重复记录
	servers, err = db.ListMediaServers()
	if err != nil {
		t.Fatalf("Failed to list media servers: %v", err)
	}
	if len(servers) != 1 {
		t.Fatalf("Expected 1 server after update, got %d", len(servers))
	}
	if servers[0].Username != "admin" || servers[0].Password != "password" {
		t.Fatalf("Server update failed: %+v", servers[0])
	}

	// 测试多次调用（模拟容器重启）
	for i := 0; i < 5; i++ {
		err = db.AddOrUpdateMediaServer("Default", "SRS", "192.168.1.100", 1985, "", "", "", 1)
		if err != nil {
			t.Fatalf("Failed on iteration %d: %v", i, err)
		}
	}

	// 验证仍然只有一条记录
	servers, err = db.ListMediaServers()
	if err != nil {
		t.Fatalf("Failed to list media servers: %v", err)
	}
	if len(servers) != 1 {
		t.Fatalf("Expected 1 server after multiple restarts, got %d", len(servers))
	}

	t.Log("Test passed: No duplicate servers created on restart")
}

func TestAddMediaServerDuplicates(t *testing.T) {
	// 创建临时数据库文件
	dbPath := "./test_media_servers_dup.db"
	defer os.Remove(dbPath)

	// 创建数据库实例
	db, err := NewMediaServerDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// 使用旧的 AddMediaServer 方法测试重复添加问题
	for i := 0; i < 3; i++ {
		err = db.AddMediaServer("Default", "SRS", "192.168.1.100", 1985, "", "", "", 1)
		if err != nil {
			t.Fatalf("Failed to add media server: %v", err)
		}
	}

	// 验证会产生重复记录
	servers, err := db.ListMediaServers()
	if err != nil {
		t.Fatalf("Failed to list media servers: %v", err)
	}
	if len(servers) != 3 {
		t.Fatalf("Expected 3 duplicate servers with old method, got %d", len(servers))
	}

	t.Log("Test confirmed: Old AddMediaServer method creates duplicates")
}

