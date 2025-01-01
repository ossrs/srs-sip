package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type MediaServerDB struct {
	db *sql.DB
}

func NewMediaServerDB(dbPath string) (*MediaServerDB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// 创建媒体服务器表
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS media_servers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			ip TEXT NOT NULL,
			port INTEGER NOT NULL,
			username TEXT,
			password TEXT,
			type TEXT NOT NULL,
			status INTEGER NOT NULL DEFAULT 0,
			is_default INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, err
	}

	return &MediaServerDB{db: db}, nil
}

func (m *MediaServerDB) AddMediaServer(name, ip string, port int, username, password, serverType string, isDefault int) error {
	_, err := m.db.Exec(`
		INSERT INTO media_servers (name, ip, port, username, password, type, status, is_default)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, name, ip, port, username, password, serverType, 0, isDefault)
	return err
}

func (m *MediaServerDB) DeleteMediaServer(id int) error {
	_, err := m.db.Exec("DELETE FROM media_servers WHERE id = ?", id)
	return err
}

func (m *MediaServerDB) GetMediaServer(id int) (*MediaServer, error) {
	var ms MediaServer
	err := m.db.QueryRow(`
		SELECT id, name, ip, port, username, password, type, status, is_default, created_at
		FROM media_servers WHERE id = ?
	`, id).Scan(&ms.ID, &ms.Name, &ms.IP, &ms.Port, &ms.Username, &ms.Password, &ms.Type, &ms.Status, &ms.IsDefault, &ms.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &ms, nil
}

func (m *MediaServerDB) ListMediaServers() ([]MediaServer, error) {
	rows, err := m.db.Query(`
		SELECT id, name, ip, port, username, password, type, status, is_default, created_at
		FROM media_servers ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []MediaServer
	for rows.Next() {
		var ms MediaServer
		err := rows.Scan(&ms.ID, &ms.Name, &ms.IP, &ms.Port, &ms.Username, &ms.Password, &ms.Type, &ms.Status, &ms.IsDefault, &ms.CreatedAt)
		if err != nil {
			return nil, err
		}
		servers = append(servers, ms)
	}
	return servers, nil
}

func (m *MediaServerDB) UpdateMediaServerStatus(id int, status int) error {
	_, err := m.db.Exec("UPDATE media_servers SET status = ? WHERE id = ?", status, id)
	return err
}

func (db *MediaServerDB) SetDefaultMediaServer(id int) error {
	// 先将所有服务器设置为非默认
	if _, err := db.db.Exec("UPDATE media_servers SET is_default = 0"); err != nil {
		return err
	}

	// 将指定ID的服务器设置为默认
	_, err := db.db.Exec("UPDATE media_servers SET is_default = 1 WHERE id = ?", id)
	return err
}

type MediaServer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Type      string `json:"type"`
	Status    int    `json:"status"`
	IsDefault int    `json:"is_default"`
	CreatedAt string `json:"created_at"`
}

func (m *MediaServerDB) Close() error {
	return m.db.Close()
}
