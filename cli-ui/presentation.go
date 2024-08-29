package cliui

import (
	"fmt"

	"github.com/christianhturner/go-claude/chat"
	"github.com/christianhturner/go-claude/logger"
)

func PresentHistoricMessagePairs(displayAmount int, messagePairs []chat.MessagePair) {
	if displayAmount == 0 {
		fmt.Println("Nothing to print")
		return
	}
	maxPairs := len(messagePairs)
	if displayAmount > maxPairs {
		logger.Error(fmt.Sprintf("Display amount provided, %d, is greater than the number of message pairs, %d.\n", displayAmount, maxPairs))
	}
	fmt.Println("\nConversationHistory:")
	for i := 0; i < displayAmount && i < maxPairs; i++ {
		pair := messagePairs[i]
		fmt.Printf("\nUser: %s\n", pair.UserMessage.Content)
		fmt.Printf("\nClaude: %s\n", pair.AssistantMessage.Content)
	}
	fmt.Println("\n")
}
