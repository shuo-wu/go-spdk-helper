package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/longhorn/go-spdk-helper/spdk"
)

func BdevAioCmd() cli.Command {
	return cli.Command{
		Name:      "bdev-aio",
		ShortName: "aio",
		Subcommands: []cli.Command{
			BdevAioCreateCmd(),
			BdevAioDeleteCmd(),
			BdevAioGetCmd(),
		},
	}
}

func BdevAioCreateCmd() cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "create a bdev aio based on a block device: create <BLOCK DEVICE PATH> <BDEV NAME> --block-size <BLOCK SIZE>",
		Flags: []cli.Flag{
			cli.Uint64Flag{
				Name:  "block-size, b",
				Usage: "The block size of created bdev aio",
				Value: 4096,
			},
		},
		Action: func(c *cli.Context) {
			if err := bdevAioCreate(c); err != nil {
				logrus.WithError(err).Fatalf("Error running create bdev aio command")
			}
		},
	}
}

func bdevAioCreate(c *cli.Context) error {
	spdkCli, err := spdk.NewClient()
	if err != nil {
		return err
	}

	bdevName, err := spdkCli.BdevAioCreate(c.Args().First(), c.Args().Get(1), c.Uint64("block-size"))
	if err != nil {
		return err
	}

	bdevAioCreateRespJson, err := json.Marshal(map[string]string{"bdev_name": bdevName})
	if err != nil {
		return err
	}
	fmt.Println(string(bdevAioCreateRespJson))

	return nil
}

func BdevAioDeleteCmd() cli.Command {
	return cli.Command{
		Name:  "delete",
		Usage: "delete a bdev aio using a block device: delete <BDEV NAME>",
		Action: func(c *cli.Context) {
			if err := bdevAioDelete(c); err != nil {
				logrus.WithError(err).Fatalf("Error running delete bdev aio command")
			}
		},
	}
}

func bdevAioDelete(c *cli.Context) error {
	spdkCli, err := spdk.NewClient()
	if err != nil {
		return err
	}

	deleted, err := spdkCli.BdevAioDelete(c.Args().First())
	if err != nil {
		return err
	}

	bdevAioDeleteRespJson, err := json.Marshal(deleted)
	if err != nil {
		return err
	}
	fmt.Println(string(bdevAioDeleteRespJson))

	return nil
}

func BdevAioGetCmd() cli.Command {
	return cli.Command{
		Name: "get",
		Flags: []cli.Flag{
			cli.Uint64Flag{
				Name:  "timeout, t",
				Usage: "Optional. Determine the timeout of the execution",
				Value: 0,
			},
		},
		Usage: "get all bdev aio if a bdev name is not specified: get <BDEV NAME>",
		Action: func(c *cli.Context) {
			if err := bdevAioGet(c); err != nil {
				logrus.WithError(err).Fatalf("Error running get bdev aio command")
			}
		},
	}
}

func bdevAioGet(c *cli.Context) error {
	spdkCli, err := spdk.NewClient()
	if err != nil {
		return err
	}

	bdevAioGetResp, err := spdkCli.BdevAioGet(c.Args().First(), 0)
	if err != nil {
		return err
	}

	bdevAioGetRespJson, err := json.Marshal(bdevAioGetResp)
	if err != nil {
		return err
	}
	fmt.Println(string(bdevAioGetRespJson))

	return nil
}
