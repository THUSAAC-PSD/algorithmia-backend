package database

import (
	"fmt"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/config"
)

type Options struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	User        string `mapstructure:"user"`
	DBName      string `mapstructure:"dbName"`
	SSLMode     bool   `mapstructure:"sslMode"`
	Password    string `mapstructure:"password"`
	UseInMemory bool   `mapstructure:"useInMemory"`
	UseSQLLite  bool   `mapstructure:"useSqlLite"`
}

func (h *Options) DNS() string {
	datasource := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		h.User,
		h.Password,
		h.Host,
		h.Port,
		"postgres",
	)

	return datasource
}

func ProvideConfig() (*Options, error) {
	return config.BindConfigKey[*Options]("gormOptions")
}
