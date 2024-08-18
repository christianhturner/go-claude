Guideline for myself for how the commands and flags may be structured. Subject to change of course.

```
claude
├── create
│   ├── conversation
│   └── chat
├── configure
│   ├── global
│   └── conversation
├── delete
│   ├── conversation
│   └── messages
├── list
├── chat
├── messages
├── export
│   ├── conversations
│   └── template
├── import
└── version
```

Flags

```
claude create conversation
Flags:

--name (string): Name of the conversation
--model (string): Claude model to use
--temperature (float): Temperature setting


claude create chat
Flags:

--name (string): Name of the conversation
--model (string): Claude model to use
--temperature (float): Temperature setting


claude configure global
Flags:

--api-key (string): Set API key
--default-model (string): Set default model
--default-temperature (float): Set default temperature


claude configure conversation
Flags:

--id (string): Conversation ID
--name (string): Update conversation name
--model (string): Update model
--temperature (float): Update temperature


claude delete conversation
Flags:

--id (string): Conversation ID to delete


claude delete messages
Flags:

--conversation-id (string): Conversation ID
--message-ids ([]string): Message IDs to delete


claude list
Flags:

--limit (int): Number of conversations to list
--offset (int): Offset for pagination


claude chat
Flags:

--id (string): Conversation ID to enter


claude history
Flags:

--id (string): Conversation ID to view history
--limit (int): Number of messages to show
--offset (int): Offset for pagination


claude export conversations
Flags:

--format (string): Export format (e.g., JSON, CSV)
--output (string): Output file path


claude export template
Flags:

--format (string): Template format (e.g., JSON, CSV)
--output (string): Output file path


claude import
Flags:

--file (string): Path to the import file


claude version
(No flags, displays version information)

Stretch Goals:

claude search

Search through conversations and messages
Flags:
--query (string): Search query
--conversation-id (string): Limit search to a specific conversation


claude analyze

Analyze conversation statistics
Flags:
--id (string): Conversation ID to analyze
--metric (string): Specific metric to analyze (e.g., sentiment, topic)


claude backup

Create a backup of all conversations
Flags:
--output (string): Output directory for backup


claude restore

Restore conversations from a backup
Flags:
--input (string): Input directory or file for restore
```
