package db

import (
	"testing"
	"time"
)

func TestListConversations(t *testing.T) {
	err := InitDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer Close()

	testConversations := []struct {
		title string
		date  time.Time
	}{
		{"Test Conversation 1", time.Now().Add(-2 * time.Hour)},
		{"Test Conversation 2", time.Now().Add(-1 * time.Hour)},
		{"Test Conversation 3", time.Now()},
	}

	for _, tc := range testConversations {
		_, err := db.Exec("INSERT INTO conversations(title, created_at, updated_at) VALUES (?, ?, ?)",
			tc.title, tc.date, tc.date)
		if err != nil {
			t.Fatalf("Failed to insert test conversation: %v", err)
		}
	}

	conversations, err := ListConversations()
	if err != nil {
		t.Fatalf("ListConversations returned an error: %v", err)
	}

	if len(conversations) != 3 {
		t.Errorf("Expected 3 conversations, got %d", len(conversations))
	}

	expectedTitles := []string{"Test Conversation 3", "Test Conversation 2", "Test Conversation 1"}
	for i, conv := range conversations {
		if conv.Title != expectedTitles[i] {
			t.Errorf("Expected conversation %d to have title %s, got %s", i, expectedTitles[i], conv.Title)
		}

		if conv.ID == 0 {
			t.Errorf("Expected non-zero ID for conversation %d", i)
		}

		expectedTime := testConversations[2-i].date
		if conv.CreatedAt.Sub(expectedTime) > time.Second || expectedTime.Sub(conv.CreatedAt) > time.Second {
			t.Errorf("CreatedAt time for conversation %d is off by more than 1 second", i)
		}

		if conv.UpdatedAt.Sub(expectedTime) > time.Second || expectedTime.Sub(conv.UpdatedAt) > time.Second {
			t.Errorf("UpdatedAt time for conversation %d is off by more than 1 second", i)
		}
	}
}

func TestCreateConversation(t *testing.T) {
	err := InitDatabase(":memory:")
	if err != nil {
		t.Fatalf("FAiled to initialize database: %v", err)
	}
	defer Close()

	title := "Test Conversation"

	id, err := CreateConversation(title)
	if err != nil {
		t.Fatalf("CreateConversation returned an error: %v", err)
	}

	if id == 0 {
		t.Errorf("Expected non-zero ID< got 0")
	}

	var storedTitle string
	var createdAt, UpdatedAt time.Time
	err = db.QueryRow("SELECT title, created_at, updated_at FROM conversations WHERE id = ?", id).
		Scan(&storedTitle, &createdAt, &UpdatedAt)
	if err != nil {
		t.Fatalf("Failed to retrieve created conversation: %v", err)
	}

	if storedTitle != title {
		t.Errorf("Expected title %s, got %s", title, storedTitle)
	}

	now := time.Now()
	if now.Sub(createdAt) > time.Minute {
		t.Errorf("created_at is not recent: %v", err)
	}
	if now.Sub(UpdatedAt) > time.Minute {
		t.Errorf("updated_at is not recent: %v", err)
	}

	id, err = CreateConversation("")
	if err != nil {
		t.Fatalf("CreateConversation with empty title returned an error: %v", err)
	}
	if id == 0 {
		t.Errorf("Expected non-zero ID for empty title conversation, got 0")
	}
	err = db.QueryRow("SELECT title FROM conversations WHERE id = ?", id).Scan(&storedTitle)
	if err != nil {
		t.Fatalf("Failed to retrieve created conversation with empty title: %v", err)
	}
	if storedTitle != "" {
		t.Errorf("Expected empty title, got %s", storedTitle)
	}
}
