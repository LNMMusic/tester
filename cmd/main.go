package main

import (
	"flag"
	"fmt"

	"github.com/LNMMusic/tester/internal/application"
)

func main() {
	// env
	// ...

	// cmd
	// - flag: config file path
	cfgFile := flag.String("config", "config.yaml", "config file path in yaml format")
	flag.Parse()
	
	// application
	// - config from yaml file
	cfg, err := application.NewConfigApplicationDefaultFromYAML(*cfgFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	// - new application
	a := application.NewApplicationDefault(cfg)
	// - tear down
	defer a.TearDown()
	// - set up
	if err := a.SetUp(); err != nil {
		fmt.Println(err)
		return
	}
	// - run
	if err := a.Run(); err != nil {
		fmt.Println(err)
		return
	}
}