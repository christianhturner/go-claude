# Go-Claude

> Attention:
> The README is not indicative of the current state of the project in the usage and command reference.
> These, at this state are my personal notes presented in a matter apt for a README assumming that
> I implement them in the manner presented. I'll implement a TODO structure by ALPHA for what items work!

Go-Claude is a powerful CLI application built on CobraCLI that allows you to interact with the Claude API. It provides a seamless interface for managing conversations, configuring settings, and interacting with Claude, all while maintaining your privacy and data control.

## Features

- Create and manage conversations with Claude
- Configure global and conversation-specific settings
- Delete conversations and messages
- List and navigate through your conversation history
- Export and import conversations
- Local storage using SQLite database
- Full control over your data

## Installation

[TODO] - Expected at ALPHA
[Instructions for installation]

## Usage

Go-Claude offers a variety of commands to interact with the Claude API:

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

For detailed information on each command and its flags, please refer to the [Command Reference](https://github.com/christianhturner/go-claude/tree/mainline/docs/command-reference.md).

## Data Privacy

Go-Claude prioritizes your privacy. All conversations, messages, and configurations are stored in a local SQLite database that you have full control over. We simply forward your interactions to the Claude API and do not store any data on our end outside of your local environment.

## Usage

Go-Claude offers a variety of commands to interact with the Claude API. Here's an overview of the main commands and their flags:

### Create

Create new conversations or chats:

```
go-claude create conversation
Flags:
  --name string         Name of the conversation
  --model string        Claude model to use
  --temperature float   Temperature setting

go-claude create chat
Flags:
  --name string         Name of the conversation
  --model string        Claude model to use
  --temperature float   Temperature setting
```

### Configure

Configure global settings or specific conversations:

```
go-claude configure global
Flags:
  --api-key string             Set API key
  --default-model string       Set default model
  --default-temperature float  Set default temperature

go-claude configure conversation
Flags:
  --id string           Conversation ID
  --name string         Update conversation name
  --model string        Update model
  --temperature float   Update temperature
```

### Delete

Delete conversations or messages:

```
go-claude delete conversation
Flags:
  --id string   Conversation ID to delete

go-claude delete messages
Flags:
  --conversation-id string   Conversation ID
  --message-ids []string     Message IDs to delete
```

### List and Chat

List conversations and interact with them:

```
go-claude list
Flags:
  --limit int    Number of conversations to list
  --offset int   Offset for pagination

go-claude chat
Flags:
  --id string   Conversation ID to enter

go-claude messages
Flags:
  --id string     Conversation ID to view history
  --limit int     Number of messages to show
  --offset int    Offset for pagination
```

### Export and Import

Export and import conversations:

```
go-claude export conversations
Flags:
  --output string   Output file path

go-claude export template
Flags:
  --output string   Output file path

go-claude import
Flags:
  --file string   Path to the import file
```

### Version

Display version information:

```
go-claude version
```

## Examples

Here are some example commands to get you started:

1. Create a new conversation:

   ```
   go-claude create conversation --name "My First Chat" --model "claude-2" --temperature 0.7
   ```

2. Configure global settings:

   ```
   go-claude configure global --api-key "your-api-key" --default-model "claude-2" --default-temperature 0.5
   ```

3. Start a chat:

   ```
   go-claude chat --id "conversation-id"
   ```

4. Export conversations:

   ```
   go-claude export conversations --format json --output "my_conversations.json"
   ```

For more detailed information on each command and its usage, you can use the `--help` flag with any command:

```
go-claude [command] --help
```

## Roadmap

- Alpha release (September 2024): Basic chat functionality and conversation storage
- Beta release: Implementation of all core features listed above
- Stable release: API stabilization and additional features based on user feedback

Future enhancements may include:

- Search functionality for conversations and messages
- Conversation analysis tools
- Backup and restore capabilities

## Contributing

[TODO]
We welcome contributions! Please see our [Contributing Guidelines](link-to-contributing.md) for more information.

## License

[License information]

---
