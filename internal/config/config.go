package config

type Config struct {
    AppName string
}

func Load() Config {
    return Config{AppName: "agnos-backend"}
}
