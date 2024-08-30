package delete

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/christianhturner/go-claude/db"
	"github.com/christianhturner/go-claude/logger"
)

func ParseMessageIDs(input string) ([]int64, error) {
	var result []int64
	parts := strings.Split(input, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}
			start, err := strconv.ParseInt(strings.TrimSpace(rangeParts[0]), 10, 64)
			if err != nil {
				return nil, err
			}
			end, err := strconv.ParseInt(strings.TrimSpace(rangeParts[1]), 10, 64)
			if err != nil {
				return nil, err
			}
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
		} else {
			id, err := strconv.ParseInt(part, 10, 64)
			if err != nil {
				return nil, err
			}
			result = append(result, id)
		}
	}
	return result, nil
}

func DeleteMessages(messageId int64) error {
	db.DeleteMessage(messageId)
	return nil
}

func DeleteConversation(conversationId int64) error {
	conversations, err := db.ListConversations()
	if err != nil {
		logger.PanicError(err, "Error listing conversations")
	}
	convExist := false
	for _, conv := range conversations {
		if conv.ID == conversationId {
			convExist = true
			break
		} else {
			continue
		}
	}
	if !convExist {
		return err
	}
	db.DeleteConversation(conversationId)
	return nil
}
