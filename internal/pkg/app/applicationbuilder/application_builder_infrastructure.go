package applicationbuilder

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/config"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/mailing"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/postmark"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/websocket"
)

func (b *ApplicationBuilder) AddInfrastructure() {
	if err := b.Container.Provide(func(cfg *config.Config) *database.Options { return &cfg.GormOptions }); err != nil {
		b.Logger.Fatal(err)
	}

	if err := b.Container.Provide(func(cfg *config.Config) *echoweb.Options { return &cfg.EchoHttpOptions }); err != nil {
		b.Logger.Fatal(err)
	}

	if err := b.Container.Provide(func(cfg *config.Config) *postmark.Options { return &cfg.PostmarkOptions }); err != nil {
		b.Logger.Fatal(err)
	}

	// Keep Gomail as fallback if needed
	if err := b.Container.Provide(func(cfg *config.Config) *mailing.Options { return &cfg.GomailOptions }); err != nil {
		b.Logger.Fatal(err)
	}

	if err := b.Container.Provide(func(cfg *config.Config) *websocket.Options { return &cfg.WebsocketOptions }); err != nil {
		b.Logger.Fatal(err)
	}

	if err := database.AddGorm(b.Container); err != nil {
		b.Logger.Fatal(err)
	}

	if err := echoweb.AddEcho(b.Container); err != nil {
		b.Logger.Fatal(err)
	}

	if err := websocket.AddWebsocket(b.Container); err != nil {
		b.Logger.Fatal(err)
	}
}
