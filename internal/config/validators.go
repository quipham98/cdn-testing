package config

type RootFolder struct {
	AllowedFiles []string `mapstructure:"allowed_files,omitempty"`
	SkipFiles    []string `mapstructure:"skip_files,omitempty"`
	SkipDirs     []string `mapstructure:"skip_dirs,omitempty"`
}

type ChainFolder struct {
	AllowedFiles []string `mapstructure:"allowed_files,omitempty"`
}

type AssetFolder struct {
	AllowedFiles []string `mapstructure:"allowed_files,omitempty"`
}

type ChainInfoFolder struct {
	HasFiles []string `mapstructure:"has_files,omitempty"`
}
