package db

import (
	"database/sql"
	"sync"

	"github.com/ossrs/srs-sip/pkg/models"
	_ "modernc.org/sqlite"
)

var (
	instance *MediaServerDB
	once     sync.Once
)

type MediaServerDB struct {
	models.MediaServerResponse
	db *sql.DB
}

// GetInstance 返回 MediaServerDB 的单例实例
func GetInstance(dbPath string) (*MediaServerDB, error) {
	var err error
	once.Do(func() {
		instance, err = NewMediaServerDB(dbPath)
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
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
			type TEXT NOT NULL,
			name TEXT NOT NULL,
			ip TEXT NOT NULL,
			port INTEGER NOT NULL,
			username TEXT,
			password TEXT,
			secret TEXT,
			is_default INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, err
	}

	return &MediaServerDB{db: db}, nil
}

func (m *MediaServerDB) AddMediaServer(name, serverType, ip string, port int, username, password, secret string, isDefault int) error {
	_, err := m.db.Exec(`
		INSERT INTO media_servers (name, type, ip, port, username, password, secret, is_default)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, name, serverType, ip, port, username, password, secret, isDefault)
	return err
}

func (m *MediaServerDB) DeleteMediaServer(id int) error {
	_, err := m.db.Exec("DELETE FROM media_servers WHERE id = ?", id)
	return err
}

func (m *MediaServerDB) GetMediaServer(id int) (*models.MediaServerResponse, error) {
	var ms models.MediaServerResponse
	err := m.db.QueryRow(`
		SELECT id, name, type, ip, port, username, password, secret, is_default, created_at
		FROM media_servers WHERE id = ?
	`, id).Scan(&ms.ID, &ms.Name, &ms.Type, &ms.IP, &ms.Port, &ms.Username, &ms.Password, &ms.Secret, &ms.IsDefault, &ms.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &ms, nil
}

func (m *MediaServerDB) ListMediaServers() ([]models.MediaServerResponse, error) {
	rows, err := m.db.Query(`
		SELECT id, name, type, ip, port, username, password, secret, is_default, created_at
		FROM media_servers ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []models.MediaServerResponse
	for rows.Next() {
		var ms models.MediaServerResponse
		err := rows.Scan(&ms.ID, &ms.Name, &ms.Type, &ms.IP, &ms.Port, &ms.Username, &ms.Password, &ms.Secret, &ms.IsDefault, &ms.CreatedAt)
		if err != nil {
			return nil, err
		}
		servers = append(servers, ms)
	}
	return servers, nil
}

func (m *MediaServerDB) SetDefaultMediaServer(id int) error {
	// 先将所有服务器设置为非默认
	if _, err := m.db.Exec("UPDATE media_servers SET is_default = 0"); err != nil {
		return err
	}

	// 将指定ID的服务器设置为默认
	_, err := m.db.Exec("UPDATE media_servers SET is_default = 1 WHERE id = ?", id)
	return err
}

func (m *MediaServerDB) Close() error {
	return m.db.Close()
}
