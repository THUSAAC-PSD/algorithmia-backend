package tzinit

import (
	"os"
)

func init() {
	err := os.Setenv("TZ", "UTC")
	if err != nil {
		panic("Failed to set timezone to UTC: " + err.Error())
	}
}
