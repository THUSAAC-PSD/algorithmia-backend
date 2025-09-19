package main

import (
    "flag"
    "fmt"
    "log"

    infra "github.com/THUSAAC-PSD/algorithmia-backend/internal/user/infrastructure"
)

func main() {
    var password string
    flag.StringVar(&password, "p", "Test@2025!", "password to hash")
    flag.Parse()

    hasher := infra.NewArgonPasswordHasher()
    hash, err := hasher.Hash(password)
    if err != nil {
        log.Fatalf("failed to hash: %v", err)
    }
    fmt.Println(hash)
}
