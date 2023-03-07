package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/longhorn/go-spdk-helper/app/cmd"
)

func main() {
	a := cli.NewApp()

	a.Before = func(c *cli.Context) error {
		if c.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}
	a.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "debug",
		},
	}
	a.Commands = []cli.Command{
		cmd.BdevAioCmd(),
		cmd.BdevLvstoreCmd(),
		cmd.BdevLvolCmd(),
		// cmd.BdevNvmeCmd(),
		// cmd.BdevRaidCmd(),
		// cmd.NvmfCmd(),
	}
	if err := a.Run(os.Args); err != nil {
		logrus.WithError(err).Fatal("Error when executing command")
	}
}
