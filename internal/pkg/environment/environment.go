package environment

import "github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/constant"

type Environment string

func (env Environment) IsDevelopment() bool {
	return env == Development
}

func (env Environment) IsProduction() bool {
	return env == Production
}

func (env Environment) GetEnvironmentName() string {
	return string(env)
}

var (
	Development = Environment(constant.Dev)
	Test        = Environment(constant.Test)
	Production  = Environment(constant.Production)
)
