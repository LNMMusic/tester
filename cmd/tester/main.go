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
	// - config: from yaml
	cfg, err := application.NewConfigApplicationDefaultFromYAML(*cfgFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	a := application.NewApplicationDefault(cfg)
	// - run
	if err := a.Run(); err != nil {
		fmt.Println(err)
		return
	}
}