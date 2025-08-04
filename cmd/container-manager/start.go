package main

import (
	"context"

	"container-manager/internal/config"
	"container-manager/internal/logger"
	"container-manager/internal/manager"

	"github.com/urfave/cli/v3"
)

func startCMD(l logger.Logger) *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "start the server",
		Flags: []cli.Flag{cfgPathsFlag()},
		Action: func(ctx context.Context, cl *cli.Command) error {
			conf, err := config.Load(cl.String("config"))
			if err != nil {
				return err
			}

			srv, err := manager.New(l, conf)
			if err != nil {
				return err
			}

			return srv.Start(ctx)
		},
	}
}
