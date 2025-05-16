package config

type App struct {
	Name    string `mapstructure:"APP_NAME"`
	Prefork bool   `mapstructure:"APP_PREFORK"`
}

type Config struct {
	App App `mapstructure:",squash"`
}
