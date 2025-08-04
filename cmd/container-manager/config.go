package main

import "github.com/urfave/cli/v3"

func cfgPathsFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Value:   "config.yaml",
		Usage:   "allows you to use your own paths to configuration files. by default it uses config.yaml",
	}
}
