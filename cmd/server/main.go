package main

import (
    "fmt"

    "agnos-test/internal/config"
)

func main() {
    cfg := config.Load()
    fmt.Printf("Agnos backend starting: %s\n", cfg.AppName)
}
