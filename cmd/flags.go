package cmd

var (
	confirm           bool   // false, "--yes", "-y"
	userMessage       string // "", "--message", "-m"
	showHistory       bool   // true, "--history", "-H"
	conversationTitle string // "", "--title", "-t"
	conversationId    int64  // 0, "--id"
	messageId         int64  // 0, "--messId"
	messageIds        string // "", "--messIds"
)

func chatCmdFlags() {
	chatCmd.Flags().StringVarP(&userMessage, "message", "m", "", "Send a message to Claude")
	chatCmd.Flags().Int64Var(&conversationId, "id", 0, "Specify a Conversation by it's ID")
	chatCmd.Flags().BoolVarP(&showHistory, "history", "H", true, "Specify whether you want to see your last messages")
}

func configureCmdFlags() {
}

func createCmdFlags() {
}

func deleteCmdFlags() {
	deleteCmd.AddCommand(deleteConversation)
	deleteConversation.Flags().Int64Var(&conversationId, "id", 0, "Specify a Conversation by it's ID")
	deleteConversation.Flags().BoolP("yes", "y", false, "Automatically confirm without prompts.")
	deleteMessages.Flags().Int64Var(&messageId, "messId", 0, "Specify a message by it's ID")
	deleteMessages.Flags().StringVar(&messageIds, "messIds", "", "Specify a range of ID's. (i.e. '1, 2, 5-7' ")
}

func exportCmdFlags() {
}

func importCmdFlags() {
}

func listCmdFlags() {
}

func rootCmdFlags() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
