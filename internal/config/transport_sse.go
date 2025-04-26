package config

type SSETransport struct {
	Address string `default:"localhost:8080" envconfig:"ADDRESS"`
}
