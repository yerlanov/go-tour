package main

import (
	"github.com/go-tour/internal"
	"github.com/go-tour/internal/config"
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
