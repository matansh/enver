# A Go(lang) configuration library developed by @matansh.

## Usage
```go
import github.com/matansh/enver

type Config struct {
	LogLevel          string `env:"LOG_LEVEL"`
	Port              int64  `env:"PORT"`
	FeatureFlag       bool   `env:"TURN_ON_FEATURE"`
	ExplicitlyIgnored string `env:"-"` // passing "-" instructs the lib not to populate this field
	ImplicitlyIgnored string           // untagged struct fields will be ignored
}

var cfg Config

errs := config.LoadEnv(&cfg)
if len(errs) != 0 {
    // failed to load config
}
```

## Background
This library is intended to help projects implement the twelve-factor app methodology - https://12factor.net/

### footnote
This library is intentionally dependency-less in order to minimize the dependency trees of its importers, you are welcome ;)
