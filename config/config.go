package config

import "github.com/spf13/cobra"

var (
	DataDir           = ""
	CfgFile           = "config.json"
	DbFile            = "data.db"
	LogLevel          = "INFO"
	AnthropicApiKey   = ""
	AnthropicUrl      = "https://api.anthropic.com/"
	AnthropicEndpoint = "v1/messages"
	AnthropicVersion  = "2023-06-01"
	AnthropicBeta     = ""
	MaxTokens         = 2000
	Model             = "claude-3-5-sonnet-20240620"
	Temparature       float64
	TopP              float64
	TopK              float64
)

func AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&DataDir, "dataDir", "dd", "Specifies a data directory for you config, db file, and logs. (Default) `{$HOME}/.conf/go-claude/")
	cmd.Flags().StringVar(&CfgFile, "confName", "conf", "Specifies a name for your config file. (Default) config.json")
	cmd.Flags().StringVar(&DbFile, "dbName", "db", "Specifies a name for your SQLite DB file name. (Default) data.db")
	cmd.Flags().StringVar(&LogLevel, "logLevel", "l", "Specifies the log level that should be set for your application. (Default) INFO.\n(Options)[\"INFO\",\"ERROR\",\"DEBUG\",\"TRACE\"]")
	cmd.Flags().StringVar(&AnthropicApiKey, "APIKey", "api", "(Global) Specifies your Anthroipic API key. When resetting config to Default settings, this value does not change.")
	cmd.Flags().StringVar(&AnthropicBeta, "ABeta", "beta", "(Global) Specifies Anthroipic Beta Options.")
}
