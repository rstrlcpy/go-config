# go-config
Simple helper that populates configuration from environment variables

#### Example

```go
package main

import (
	config "github.com/Ramzec/go-config"
)

type ServiceConfig struct {
	BooleanField         bool   `config:"ENV_VAR1"`
	StringField          string `config:"ENV_VAR2"`
	RequiredStringField  string `config:"ENV_VAR3,required"`
	IntegerField         int    `config:"ENV_VAR4"`
	RequiredIntegerField int    `config:"ENV_VAR5,required"`
}

func main() {
	var serviceConfig ServiceConfig
	config.BuildConfig(&serviceConfig)
}
```

If `required` environment variable is not defined, then hte application will
panic on start.
If non `required` environment variable is not defined, then the corresponding
`string` field will contain "" (empty string), `int` field will contain 0
