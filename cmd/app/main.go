package main

import (
	"os"

	"github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/app"

	"github.com/spf13/cobra"

	_ "github.com/THUSAAC-PSD/algorithmia-backend/internal/pkg/tzinit"
)

var rootCmd = &cobra.Command{
	Use:              "algorithmia-backend",
	TraverseChildren: true,
	Run: func(_ *cobra.Command, _ []string) {
		app.NewApp().Run()
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
