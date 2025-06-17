package database

type Options struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	User        string `mapstructure:"user"`
	DBName      string `mapstructure:"dbName"`
	SSLMode     string `mapstructure:"sslMode"`
	Password    string `mapstructure:"password"`
	UseInMemory bool   `mapstructure:"useInMemory"`
}
