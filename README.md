# Go-Claude

> Attention:
> The README is not indicative of the current state of the project in the usage and command reference.
> These, at this state are my personal notes presented in a matter apt for a README assumming that
> I implement them in the manner presented.

Go-Claude is a powerful CLI application built on CobraCLI that allows you to interact with the Claude API. It provides a seamless interface for managing conversations, configuring settings, and interacting with Claude, all while maintaining your privacy and data control.

## Features

- Create and manage conversations with Claude
- Configure global and conversation-specific settings
- Delete conversations and messages
- List and navigate through your conversation history
- Export and import conversations
- Local storage using SQLite database
- Full control over your data

## Roadmap to Beta

I'm excited to announced that this project works, and I'll be classifying it as alpha. I'll work to ensure that the mainline branch
remains stable and not push changes that will break the current functionality. In it's alpha state, I'll have two additional branches
apart from the mainline, development and bubbletea. The development branch will be were I'm working to add additional core features
presented in the below roadmap. Bubbletea is were I'll be working to implement the TUI user interface. While the application will work
during the alpha to beta roadmap, there may be lots of changes to commands and functionality, but I will ensure to maintain that anything
published to mainline will continue to be able to perform the features that have been stated as functional.

The application is 80 percent of the way to being feature complete as compared to using other popular interfaces for interacting with LLMs.
You can create "conversations" (a message thread to categorize your topics), basic configuration (global), list conversations and messages,
and delete conversations and messages.

Configuration, at the moment is limited. I'd primarily stick to adding your API key and adjusting your token usage. I have settings for temperature,
topP, topK, changing default data directories, etc, but I honestly haven't tested these when using alternative locations. For some things, config file,
log file, and db, it should work if you change it. If you already have a db file, log file, or config file, and change those values, it would recreate them,
it will not migrate those. For now, keeping it at the default is recommended. I have not created a lot of unit test, so this is something I'll be wanting to
complete before adding a lot of new features.

The application, as I stated is pretty feature complete, you can create new conversations, you can chat and use those messages as context, and you can delete
messages and conversations to manage your data all within the application. With that being said, you have to be cautious when deleting messages, in particular.
You can safely delete conversations, but you have to be cautious with deleting messages as you could break a conversation from being able to work. To understand
this limitation, you have to understand how anthropic expects chat messages. Every message, except for the first message, is expected in pairs. If you send a message,
claude will reply. For your third message, if you want to have context of the messages included, anthropic expects a user message and assistant (claude) message to be
sent in pairs. For example, if you have a message thread with 10 messages you'd have messages with the following roles:
`[user, assistant, user, assistant, user, assistant, user, assistant, user, assistant]`
If you wanted to delete messages, you have to ensure that you delete both the user and assistant message for that "pair". Otherwise you're messages after this will
receive an error when you attempt to chat again. You'll then need to manually fix this problem by either deleting the conversation and recreating it (losing the messages)
or you'll have to go in and delete messages to ensure that the pairs are provided.

Due to this limitation, you have to be careful with the delete message functionality. Most of the time I do not recommend deleting messages via flags, as it's more likely
that you'll make a mistake and corrupting this conversation. I have a multi-select prompt that is safer, because you can see the messages as you select them for deletion.
Even still, this will not stop you from deleting single messages and thereby breaking the pairs. I have a struct that tethers together messages as pairs, but I was hesitant
to force this, as someone may want to use the delete feature to selectively delete between pairs while still ensuring that pairs remain together. I have some other ideas but
they require implementation and some thought.

I appreciate everyone who tries this out, feel free to open issues and suggest features. For the latest development, if you're not scared of breaking changes, check out the
development branch. If you want to start using the TUI interface once that is up and running, check out the bubbletea branch.

- [x] create
  - [x] conversation
- [ ] configure
  - [x] global
  - [ ] conversation
- [x] list
  - [x] Conversation
  - [x] Messages
- [x] chat
  - [x] without flags (Need to just add the chat client to chat.)
  - [x] with flags (Technically done; Just have to implement the logic to differentiate; defaulted to the interactive)
- [x] delete
- [ ] messages
- [ ] import
- [ ] export

## Installation

Unfortunately the installation isn't super easy at this point, and I do not have a release channel at the moment.
For installation, I recommend cloning the project so it's more easy to update with new changes. You'll also,
for the moment, will have to install go in order to compile and insall the project easily. You can install
Go from their website [here](https://go.dev/dl/). I haven't done a lot of testing regarding the oldest Version
of go that will be supported by the project. I used Generics and `any`, so I think you'll need to atleast
have go 1.18+, as stated [here](https://go.dev/doc/tutorial/generics). Once you have installed Golang, set up your paths, run the following commands:

```shell
cd ~/path/to/go-claude/
go build
go install
```

If you already set your paths up correctly, you should just be able to run go-claude and it works.
By default, Go installs packages to the $GOBIN environment variable, which defaults to $GOPATH/bin or
$HOME/go/bin if the GOPATH environment variables are not set. For me, this is at `$HOME/go/bin/`. You
can print your go environment variables by running`go env`.

This should not require cross compiling. I intentionally chose to use very few packages, and the ones I did,
to not use CGO, to help ensure that cross platform support was as easy as possible. I imagine that their will
likely be issues as far as terminals are concerned. I am primarily developing this on a MAC, but I work on
Windows and Linux as well, so I should be able to text any issues that you may fine. There are some places
that I'm using unicode encoding, for additional icons. If your terminal does not support that, and only supports
ASCII please report those findings. I can always go back and try to add additional implementations to try to
achieve better support.

## Adding Completions

### OH MY ZSH

Run the following:

```
go-claude completions zsh > ~/.oh-my-zsh/custom/go-claude_completions.zsh
```

For now, you'll need to edit this file, and delete the top line, as it will include a println statement that will
break the completion. Use you're editor, go to the file and delete that top line.
Then:

```
chmod +x ~/.oh-my-zsh/custom/go-claude_completions.zsh
echo 'source $ZSH_CUSTOM/go-claude_completions.zsh' >> ~/.zshrc
source ~/.zshrc
```

More Coming Soon...

## Usage

Go-Claude offers a variety of commands to interact with the Claude API:

```
claude
├── create
│   ├── conversations
│   └── chat
├── configure
│   ├── global
│   └── conversations
├── delete
│   ├── conversations
│   └── messages
├── list
│   ├── conversations
│   └── messages
├── chat
├── messages
├── export
│   ├── conversations
│   └── template
└── import
```

For detailed information on each command and its flags, please refer to the [Command Reference](https://github.com/christianhturner/go-claude/tree/mainline/docs/command-reference.md).

## Data Privacy

Go-Claude prioritizes your privacy. All conversations, messages, and configurations are stored in a local SQLite database that you have full control over. We simply forward your interactions to the Claude API and do not store any data on our end outside of your local environment.

## Usage

> The usage here is largely going to be not correct, please refer to the usage according to the help menu of the application.
I will spend some time going through the usage here as things are more solidified. If you'd like to explore functionality here
on github, got to the /cmd/ directory and each command can be found there. It's pretty easy to reason about I think.

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
