package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/christianhturner/go-claude/terminal"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ConfigItem struct {
	Flag      string
	ConfigKey string
	Value     interface{}
}

var (
	DataDir           = ""
	CfgFile           = ""
	DbFile            = ""
	LogLevel          = "INFO"
	AnthropicApiKey   = ""
	AnthropicUrl      = "https://api.anthropic.com/"
	AnthropicEndpoint = "v1/messages"
	AnthropicVersion  = "2023-06-01"
	AnthropicBeta     = ""
	MaxTokens         = 2000
	Model             = "claude-3-5-sonnet-20240620"
	Temperature       float64 // might should be a string in order to support an empty value?
	TopP              float64 // might should be a string in order to support an empty value?
	TopK              float64 // might should be a string in order to support an empty value?

	DataDirKey           = "data_dir"
	CfgFileKey           = "cfg_file"
	DbFileKey            = "db_file"
	LogLevelKey          = "log_level"
	AnthropicApiKeyKey   = "anthropic_api_key"
	AnthropicUrlKey      = "anthropic_url"
	AnthropicEndpointKey = "anthropic_endpoint"
	AnthropicVersionKey  = "anthropic_version"
	AnthripicBetaKey     = "anthropic_beta"
	MaxTokensKey         = "max_tokens"
	ModelKey             = "model_key"
	TemperatureKey       = "temperature_key"
	TopPKey              = "top_p"
	TopKKey              = "top_k"
)

var ConfigItems = []ConfigItem{
	{Flag: "max-tokens", ConfigKey: MaxTokensKey, Value: &MaxTokens},
	{Flag: "data-dir", ConfigKey: DataDirKey, Value: &DataDir},
	{Flag: "cfg-file", ConfigKey: CfgFileKey, Value: &CfgFile},
	{Flag: "db-file", ConfigKey: DbFileKey, Value: &DbFile},
	{Flag: "log-level", ConfigKey: LogLevelKey, Value: &LogLevel},
	{Flag: "anthropic-url", ConfigKey: AnthropicUrlKey, Value: &AnthropicUrl},
	{Flag: "anthropic-endpoint", ConfigKey: AnthropicEndpointKey, Value: &AnthropicEndpoint},
	{Flag: "anthropic-version", ConfigKey: AnthropicVersionKey, Value: &AnthropicVersion},
	{Flag: "anthropic-beta", ConfigKey: AnthripicBetaKey, Value: &AnthropicBeta},
	{Flag: "max-tokens", ConfigKey: MaxTokensKey, Value: &MaxTokens},
	{Flag: "model", ConfigKey: ModelKey, Value: &Model},
	{Flag: "temperature", ConfigKey: TemperatureKey, Value: &Temperature},
	{Flag: "top-p", ConfigKey: TopPKey, Value: &TopP},
	{Flag: "top-k", ConfigKey: TopKKey, Value: &TopK},
}

func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&DataDir, "data-dir", DataDir, "Specifies a data directory for your config, db file, and logs. (Default: $HOME/.conf/go-claude/)")

	cmd.PersistentFlags().StringVar(&CfgFile, "config-file", CfgFile, "Specifies a name for your config file. (Default: config.json)")

	cmd.PersistentFlags().StringVar(&DbFile, "database-file", DbFile, "Specifies a name for your SQLite DB file. (Default: data.db)")

	cmd.PersistentFlags().StringVar(&LogLevel, "log-level", LogLevel, "Specifies the log level. (Default: INFO, Options: INFO, ERROR, DEBUG, TRACE)")

	cmd.PersistentFlags().StringVar(&AnthropicApiKey, "api-key", AnthropicApiKey, "Specifies your Anthropic API key. (Global, persists across resets)")

	cmd.PersistentFlags().StringVar(&AnthropicUrl, "anthropic-url", AnthropicUrl, "Specifies the Anthropic API URL. (Global, Default: https://api.anthropic.com/)")

	cmd.PersistentFlags().StringVar(&AnthropicEndpoint, "anthropic-endpoint", AnthropicEndpoint, "Specifies the Anthropic API endpoint. (Global, Default: v1/messages)")

	cmd.PersistentFlags().StringVar(&AnthropicVersion, "anthropic-version", AnthropicVersion, "Specifies the Anthropic API version. (Global, Default: 2023-06-01)")

	cmd.PersistentFlags().StringVar(&AnthropicBeta, "beta-options", AnthropicBeta, "Specifies Anthropic Beta options. (Global)")

	cmd.PersistentFlags().IntVar(&MaxTokens, "max-tokens", MaxTokens, "Specifies the maximum number of tokens for the conversation. (Global)")

	cmd.PersistentFlags().StringVar(&Model, "model", Model, "Specifies the Claude model to use. (Global, Default: claude-3-5-sonnet-20240620, Options: claude-3-5-sonnet, claude-3-sonnet, claude-3-haiku, claude-3-opus)")

	cmd.PersistentFlags().Float64Var(&Temperature, "temperature", Temperature, "Specifies the temperature for response generation. (Global)")

	cmd.PersistentFlags().Float64Var(&TopP, "top-p", TopP, "Specifies the top-p value for response generation. (Global)")

	cmd.PersistentFlags().Float64Var(&TopK, "top-k", TopK, "Specifies the top-k value for response generation. (Global)")
}

func InitConfig() {
	setDefaults()
	if CfgFile != "" {
		viper.SetConfigName(CfgFile)
	} else {
		CfgFile = "config"
	}
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	DataDir = filepath.Join(home, ".config", "go-claude")
	err = os.MkdirAll(DataDir, 0755)
	cobra.CheckErr(err)

	viper.AddConfigPath(DataDir)
	viper.SetConfigName(CfgFile)
	viper.SetConfigType("json")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// Config file not found; ignore error if desired
		viper.SafeWriteConfig()
	} else {
		// Config file was found but another error was produced
		fmt.Printf("Error reading config file: %s\n", err)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", e.Name)
	})
	viper.WatchConfig()

	checkApiKey := viper.GetString("Anthropic_API_Key")
	if checkApiKey == "" {
		term := terminal.New()
		userInput, err := term.Prompt("Please provide your Anthroipic API Key:\n")
		if err != nil {
			fmt.Errorf("Error requesting user input for API key: %v", err)
		}
		viper.Set("Anthropic_API_Key", userInput)
		viper.WriteConfig()
	}
}

func setDefaults() {
	for _, item := range ConfigItems {
		if item.ConfigKey == AnthropicApiKeyKey {
			apiConfigValue := viper.GetString(AnthropicApiKeyKey)
			viper.Set(AnthropicApiKeyKey, apiConfigValue)
		} else {
			switch v := item.Value.(type) {
			case *string:
				viper.SetDefault(item.ConfigKey, *v)
			case *int:
				viper.SetDefault(item.ConfigKey, *v)
			case *float64:
				viper.SetDefault(item.ConfigKey, *v)
			default:
				// Handle other types or log an error
				fmt.Printf("Unsupported type for config key: %s\n", item.ConfigKey)
			}
		}
	}
}

func ResetToDefaults() {
	for _, item := range ConfigItems {
		if item.ConfigKey == AnthropicApiKeyKey {
			apiConfigValue := viper.GetString(AnthropicApiKeyKey)
			viper.Set(AnthropicApiKeyKey, apiConfigValue)
		} else {
			switch v := item.Value.(type) {
			case *string:
				viper.Set(item.ConfigKey, *v)
			case *int:
				viper.Set(item.ConfigKey, *v)
			case *float64:
				viper.SetDefault(item.ConfigKey, *v)
			}
		}
	}
	viper.WriteConfig()
}

func UpdateConfig(cmd *cobra.Command) {
	cmd.Flags().Visit(func(f *pflag.Flag) {
		for _, item := range ConfigItems {
			if f.Name == item.Flag {
				viper.Set(item.ConfigKey, f.Value.String())
				break
			}
		}
	})
	viper.WriteConfig()
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}
