package main

import (
	"flag"
	"github.com/cacos-group/cacos/internal/conf"
	"github.com/cacos-group/cacos/internal/di"
)

func main() {
	flag.Parse()

	config := conf.Init()

	app, cf, err := di.InitApp(config)
	if err != nil {
		panic(err)
	}

	defer cf()

	err = app.Start()
	if err != nil {
		panic(err)
	}
}
