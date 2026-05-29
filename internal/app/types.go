package app

type PlatformConfig struct {
	Label        string   `json:"label"`
	Command      string   `json:"command"`
	Args         []string `json:"args"`
	PromptMode   string   `json:"promptMode"`
	Subscription string   `json:"subscription"`
	SetupURL     string   `json:"setupUrl"`
	Enabled      bool     `json:"enabled"`
}

type Config struct {
	Version             int                       `json:"version"`
	DefaultPlatform     string                    `json:"defaultPlatform"`
	PlatformPriority    []string                  `json:"platformPriority"`
	RequireFreshData    bool                      `json:"requireFreshData"`
	AllowTradeExecution bool                      `json:"allowTradeExecution"`
	Platforms           map[string]PlatformConfig `json:"platforms"`
}

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}
