package memory

import (
	"database/sql"
	"encoding/json"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Store struct {
	db  *sql.DB
	mu  sync.Mutex
	max int
}

func New(path string, maxHistory int) (*Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS history (
			user_id INTEGER NOT NULL,
			messages TEXT NOT NULL,
			updated_at INTEGER NOT NULL,
			PRIMARY KEY (user_id)
		);
	`); err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db, max: maxHistory}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Get(userID int64) ([]Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var raw string
	err := s.db.QueryRow(`SELECT messages FROM history WHERE user_id = ?`, userID).Scan(&raw)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var msgs []Message
	if err := json.Unmarshal([]byte(raw), &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

func (s *Store) Append(userID int64, role, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	msgs, _ := s.getUnlocked(userID)
	msgs = append(msgs, Message{Role: role, Content: content})
	if len(msgs) > s.max {
		msgs = msgs[len(msgs)-s.max:]
	}
	raw, err := json.Marshal(msgs)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`
		INSERT INTO history (user_id, messages, updated_at) VALUES (?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET messages = excluded.messages, updated_at = excluded.updated_at
	`, userID, string(raw), time.Now().Unix())
	return err
}

func (s *Store) Clear(userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec(`DELETE FROM history WHERE user_id = ?`, userID)
	return err
}

func (s *Store) getUnlocked(userID int64) ([]Message, error) {
	var raw string
	err := s.db.QueryRow(`SELECT messages FROM history WHERE user_id = ?`, userID).Scan(&raw)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var msgs []Message
	_ = json.Unmarshal([]byte(raw), &msgs)
	return msgs, nil
}
