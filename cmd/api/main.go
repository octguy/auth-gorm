package main

import (
	"fmt"

	"github.com/octguy/auth-gorm/config"
)

func main() {
	cfg := config.Load()

	fmt.Println("hello")
	fmt.Println(cfg.Port)
	fmt.Println(cfg.JWTSecret)
	fmt.Println(cfg.GinMode)
}
