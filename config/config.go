package config

type Data struct {
	// Server config
	ServerHost         string `env:"SERVER_HOSTNAME" envDefault:"gofizzbuzz"`
	ServerProtocol     string `env:"SERVER_PROTOCOL" envDefault:"http"`
	ServerPort         int    `env:"SERVER_PORT" envDefault:"8080"`
	// log level mode
	LogLevel           string `env:"LOG_LEVEL" envDefault:"info"`
	// persistent prometheus metrics
	RequestDataPath    string `env:"REQUEST_DATA_PATH" envDefault:"server_data/request_data.txt"`
}
