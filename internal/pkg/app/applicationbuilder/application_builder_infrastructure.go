package applicationbuilder

import (
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/database"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/http/echoweb"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/mailing"
	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/websocket"
)

func (b *ApplicationBuilder) AddInfrastructure() {
	if err := database.AddGorm(b.Container); err != nil {
		b.Logger.Fatal(err)
	}

	if err := echoweb.AddEcho(b.Container); err != nil {
		b.Logger.Fatal(err)
	}

	if err := websocket.AddWebsocket(b.Container); err != nil {
		b.Logger.Fatal(err)
	}

	if err := mailing.AddGomail(b.Container); err != nil {
		b.Logger.Fatal(err)
	}
}
