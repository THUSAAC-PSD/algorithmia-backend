package applicationbuilder

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/config"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/gomail"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb"
)

func (b *ApplicationBuilder) AddInfrastructure() {
	err := config.AddAppConfig(b.Container)
	if err != nil {
		b.Logger.Fatal(err)
	}

	err = database.AddGorm(b.Container)
	if err != nil {
		b.Logger.Fatal(err)
	}

	err = echoweb.AddEcho(b.Container)
	if err != nil {
		b.Logger.Fatal(err)
	}

	err = gomail.AddGomail(b.Container)
	if err != nil {
		b.Logger.Fatal(err)
	}
}
