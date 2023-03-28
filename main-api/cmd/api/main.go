package main

import (
	"github.com/yerlanov/go-tour/main-api/internal"
	"github.com/yerlanov/go-tour/main-api/internal/config"
)

func main() {
	// init app
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	err = internal.NewApp(cfg).Run()
	if err != nil {
		panic(err)
	}
}
