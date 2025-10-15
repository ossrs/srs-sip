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

func TestGetMediaServer(t *testing.T) {
	dbPath := "./test_get_media_server.db"
	defer os.Remove(dbPath)

	db, err := NewMediaServerDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// 添加一个媒体服务器
	err = db.AddMediaServer("TestServer", "ZLM", "192.168.1.200", 8080, "admin", "pass123", "secret123", 0)
	if err != nil {
		t.Fatalf("Failed to add media server: %v", err)
	}

	// 获取服务器列表以获得ID
	servers, err := db.ListMediaServers()
	if err != nil {
		t.Fatalf("Failed to list media servers: %v", err)
	}
	if len(servers) == 0 {
		t.Fatal("No servers found")
	}

	// 通过ID获取服务器
	server, err := db.GetMediaServer(servers[0].ID)
	if err != nil {
		t.Fatalf("Failed to get media server: %v", err)
	}

	// 验证数据
	if server.Name != "TestServer" {
		t.Errorf("Expected name 'TestServer', got '%s'", server.Name)
	}
	if server.Type != "ZLM" {
		t.Errorf("Expected type 'ZLM', got '%s'", server.Type)
	}
	if server.IP != "192.168.1.200" {
		t.Errorf("Expected IP '192.168.1.200', got '%s'", server.IP)
	}
	if server.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", server.Port)
	}
}

func TestGetMediaServerNotFound(t *testing.T) {
	dbPath := "./test_get_not_found.db"
	defer os.Remove(dbPath)

	db, err := NewMediaServerDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// 尝试获取不存在的服务器
	_, err = db.GetMediaServer(999)
	if err == nil {
		t.Error("Expected error when getting non-existent server, got nil")
	}
}

func TestDeleteMediaServer(t *testing.T) {
	dbPath := "./test_delete_media_server.db"
	defer os.Remove(dbPath)

	db, err := NewMediaServerDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// 添加两个服务器
	err = db.AddMediaServer("Server1", "SRS", "192.168.1.1", 1985, "", "", "", 0)
	if err != nil {
		t.Fatalf("Failed to add server1: %v", err)
	}

	err = db.AddMediaServer("Server2", "ZLM", "192.168.1.2", 8080, "", "", "", 0)
	if err != nil {
		t.Fatalf("Failed to add server2: %v", err)
	}

	// 获取服务器列表
	servers, err := db.ListMediaServers()
	if err != nil {
		t.Fatalf("Failed to list servers: %v", err)
	}
	if len(servers) != 2 {
		t.Fatalf("Expected 2 servers, got %d", len(servers))
	}

	// 删除第一个服务器
	err = db.DeleteMediaServer(servers[0].ID)
	if err != nil {
		t.Fatalf("Failed to delete server: %v", err)
	}

	// 验证只剩一个服务器
	servers, err = db.ListMediaServers()
	if err != nil {
		t.Fatalf("Failed to list servers after delete: %v", err)
	}
	if len(servers) != 1 {
		t.Fatalf("Expected 1 server after delete, got %d", len(servers))
	}
}

func TestSetDefaultMediaServer(t *testing.T) {
	dbPath := "./test_set_default.db"
	defer os.Remove(dbPath)

	db, err := NewMediaServerDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// 添加三个服务器
	err = db.AddMediaServer("Server1", "SRS", "192.168.1.1", 1985, "", "", "", 1)
	if err != nil {
		t.Fatalf("Failed to add server1: %v", err)
	}

	err = db.AddMediaServer("Server2", "ZLM", "192.168.1.2", 8080, "", "", "", 0)
	if err != nil {
		t.Fatalf("Failed to add server2: %v", err)
	}

	err = db.AddMediaServer("Server3", "SRS", "192.168.1.3", 1985, "", "", "", 0)
	if err != nil {
		t.Fatalf("Failed to add server3: %v", err)
	}

	// 获取服务器列表
	servers, err := db.ListMediaServers()
	if err != nil {
		t.Fatalf("Failed to list servers: %v", err)
	}

	// 找到 Server2 的 ID
	var server2ID int
	for _, s := range servers {
		if s.Name == "Server2" {
			server2ID = s.ID
			break
		}
	}

	// 设置 Server2 为默认
	err = db.SetDefaultMediaServer(server2ID)
	if err != nil {
		t.Fatalf("Failed to set default server: %v", err)
	}

	// 验证只有 Server2 是默认的
	servers, err = db.ListMediaServers()
	if err != nil {
		t.Fatalf("Failed to list servers: %v", err)
	}

	defaultCount := 0
	for _, s := range servers {
		if s.IsDefault == 1 {
			defaultCount++
			if s.Name != "Server2" {
				t.Errorf("Expected Server2 to be default, got %s", s.Name)
			}
		}
	}

	if defaultCount != 1 {
		t.Errorf("Expected exactly 1 default server, got %d", defaultCount)
	}
}

func TestGetMediaServerByNameAndIP(t *testing.T) {
	dbPath := "./test_get_by_name_ip.db"
	defer os.Remove(dbPath)

	db, err := NewMediaServerDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// 添加服务器
	err = db.AddMediaServer("MyServer", "SRS", "10.0.0.1", 1985, "user", "pass", "secret", 0)
	if err != nil {
		t.Fatalf("Failed to add server: %v", err)
	}

	// 通过名称和IP查询
	server, err := db.GetMediaServerByNameAndIP("MyServer", "10.0.0.1")
	if err != nil {
		t.Fatalf("Failed to get server by name and IP: %v", err)
	}

	if server.Name != "MyServer" || server.IP != "10.0.0.1" {
		t.Errorf("Server data mismatch: %+v", server)
	}

	// 查询不存在的组合
	_, err = db.GetMediaServerByNameAndIP("MyServer", "10.0.0.2")
	if err == nil {
		t.Error("Expected error for non-existent name/IP combination, got nil")
	}
}
