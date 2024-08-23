package db

// ListConversations: Retrieves all conversations from the database.
// CreateConversation: Creates a new conversation.
// ConfigureConversation: Sets or updates an option for a specific conversation.
// GetConversationOptions: Retrieves all options for a specific conversation.
// DeleteConversation: Deletes a conversation and all its associated messages and options.
// AddMessage: Adds a new message to a conversation.
// GetMessages: Retrieves all messages for a specific conversation.
// DeleteMessage: Deletes a specific message from a conversation.
// EditMessage: Updates the content of a specific message.
// SetGlobalOption: Sets or updates a global option (using conversation_id 0).
// GetGlobalOptions: Retrieves all global options.
// GetConversation: Retrieves a specific conversation by ID.
// UpdateConversationTitle: Updates the title of a conversation.
// BeginTransaction: Starts a new database transaction for more complex operations.

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Conversation represents a conversation in the database
type Conversation struct {
	ID        int64
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Message represents a message in the database
type Message struct {
	ID             int64
	ConversationID int64
	Role           string
	Content        string
	CreatedAt      time.Time
}

// ConversationOption represents an option for a conversation
type ConversationOption struct {
	ID             int64
	ConversationID int64
	OptionName     string
	OptionValue    string
}

// ListConversations retrieves all conversations from the database
func ListConversations() ([]Conversation, error) {
	rows, err := db.Query("SELECT id, title, created_at, updated_at FROM conversations ORDER BY updated_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var c Conversation
		err := rows.Scan(&c.ID, &c.Title, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, c)
	}
	return conversations, nil
}

// CreateConversation creates a new conversation in the database
func CreateConversation(title string) (int64, error) {
	result, err := db.Exec("INSERT INTO conversations (title) VALUES (?)", title)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// ConfigureConversation sets or updates an option for a conversation
func ConfigureConversation(conversationID int64, optionName, optionValue string) error {
	_, err := db.Exec(`
		INSERT INTO conversation_options (conversation_id, option_name, option_value)
		VALUES (?, ?, ?)
		ON CONFLICT(conversation_id, option_name) DO UPDATE SET option_value = ?`,
		conversationID, optionName, optionValue, optionValue)
	return err
}

// GetConversationOptions retrieves all options for a conversation
func GetConversationOptions(conversationID int64) ([]ConversationOption, error) {
	rows, err := db.Query("SELECT id, conversation_id, option_name, option_value FROM conversation_options WHERE conversation_id = ?", conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []ConversationOption
	for rows.Next() {
		var o ConversationOption
		err := rows.Scan(&o.ID, &o.ConversationID, &o.OptionName, &o.OptionValue)
		if err != nil {
			return nil, err
		}
		options = append(options, o)
	}
	return options, nil
}

// DeleteConversation deletes a conversation and all its messages and options
func DeleteConversation(conversationID int64) error {
	_, err := db.Exec("DELETE FROM conversations WHERE id = ?", conversationID)
	return err
}

// AddMessage adds a new message to a conversation
func AddMessage(conversationID int64, role, content string) error {
	_, err := db.Exec("INSERT INTO messages (conversation_id, role, content) VALUES (?, ?, ?)",
		conversationID, role, content)
	return err
}

// GetMessages retrieves all messages for a conversation
func GetMessages(conversationID int64) ([]Message, error) {
	rows, err := db.Query("SELECT id, conversation_id, role, content, created_at FROM messages WHERE conversation_id = ? ORDER BY created_at ASC", conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		err := rows.Scan(&m.ID, &m.ConversationID, &m.Role, &m.Content, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

// DeleteMessage deletes a specific message from a conversation
func DeleteMessage(messageID int64) error {
	_, err := db.Exec("DELETE FROM messages WHERE id = ?", messageID)
	return err
}

// EditMessage updates the content of a specific message
func EditMessage(messageID int64, newContent string) error {
	_, err := db.Exec("UPDATE messages SET content = ? WHERE id = ?", newContent, messageID)
	return err
}

// SetGlobalOption sets or updates a global option
func SetGlobalOption(optionName, optionValue string) error {
	_, err := db.Exec(`
		INSERT INTO conversation_options (conversation_id, option_name, option_value)
		VALUES (0, ?, ?)
		ON CONFLICT(conversation_id, option_name) DO UPDATE SET option_value = ?`,
		optionName, optionValue, optionValue)
	return err
}

// GetGlobalOptions retrieves all global options
func GetGlobalOptions() ([]ConversationOption, error) {
	return GetConversationOptions(0)
}

// GetConversation retrieves a specific conversation by ID
func GetConversation(conversationID int64) (*Conversation, error) {
	var c Conversation
	err := db.QueryRow("SELECT id, title, created_at, updated_at FROM conversations WHERE id = ?", conversationID).
		Scan(&c.ID, &c.Title, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Conversation not found
		}
		return nil, err
	}
	return &c, nil
}

// UpdateConversationTitle updates the title of a conversation
func UpdateConversationTitle(conversationID int64, newTitle string) error {
	_, err := db.Exec("UPDATE conversations SET title = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", newTitle, conversationID)
	return err
}

// BeginTransaction starts a new database transaction
// Examples of use:
// // Start a new transaction
// tx, err := db.BeginTransaction()
//
//	if err != nil {
//	    // Handle error
//	    return err
//	}
//
// // Defer a rollback in case anything fails
// defer tx.Rollback()
//
// // Perform your database operations using the transaction
// // Instead of db.Exec or db.Query, use tx.Exec or tx.Query
//
// // If all operations are successful, commit the transaction
// err = tx.Commit()
//
//	if err != nil {
//	    // Handle commit error
//	    return err
//	}
//	Example use case:
//
// You want to create a new conversation and immediately add a message to it.
// IN this instance, you want these operations to be atomic (i.e., both succeed
// or both fail.)
func BeginTransaction() (*sql.Tx, error) {
	return db.BeginTx(context.Background(), nil)
}
