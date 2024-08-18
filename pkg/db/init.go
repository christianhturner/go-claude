package db

import (
	"context"
	"database/sql"

	"github.com/christianhturner/go-claude/pkg/log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDatabase(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	err = createConversationsTable()
	if err != nil {
		return err
	}

	err = createMessagesTable()
	if err != nil {
		return err
	}

	err = createConversationOptionsTable()
	if err != nil {
		return err
	}

	return nil
}

func Close() {
	err := db.Close()
	log.WarnError(err, "Error closing database")
}

func createConversationsTable() error {
	_, err := db.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS conversations(
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
	)
	log.LogError(err, "Failed to create conversations table")
	return err
}

func createMessagesTable() error {
	_, err := db.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    conversation_id INTEGER,
    role TEXT CHECK(role IN ('user', 'assistant')),
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
)`,
	)
	log.LogError(err, "Failed to create messages table")
	return err
}

func createConversationOptionsTable() error {
	_, err := db.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS conversation_options(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        conversation_id INTEGER,
        option_name TEXT,
        option_value TEXT,
        FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
        )`,
	)
	log.LogError(err, "Failed to create conversation options table")
	return err
}
