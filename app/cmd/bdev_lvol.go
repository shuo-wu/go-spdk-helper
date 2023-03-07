package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/longhorn/go-spdk-helper/spdk"
)

func BdevLvolCmd() cli.Command {
	return cli.Command{
		Name:      "bdev-lvol",
		ShortName: "lvol",
		Subcommands: []cli.Command{
			BdevLvolCreateCmd(),
			BdevLvolDeleteCmd(),
			BdevLvolGetCmd(),
		},
	}
}

func BdevLvolCreateCmd() cli.Command {
	return cli.Command{
		Name: "create",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "lvs-name",
			},
			cli.StringFlag{
				Name: "lvol-name",
			},
			cli.Uint64Flag{
				Name: "size",
			},
		},
		Usage: "create a bdev lvol in a lvstore: \"create --lvs-name <LVSTORE NAME> --lvol-name <LVOL NAME> --size <LVOL SIZE>\"",
		Action: func(c *cli.Context) {
			if err := bdevLvolCreate(c); err != nil {
				logrus.WithError(err).Fatalf("Error running create bdev lvol command")
			}
		},
	}
}

func bdevLvolCreate(c *cli.Context) error {
	spdkCli, err := spdk.NewClient()
	if err != nil {
		return err
	}

	uuid, err := spdkCli.BdevLvolCreate(c.String("lvs-name"), c.String("lvol-name"), c.Uint64("size"))
	if err != nil {
		return err
	}

	bdevLvolCreateRespJson, err := json.Marshal(map[string]string{"uuid": uuid})
	if err != nil {
		return err
	}
	fmt.Println(string(bdevLvolCreateRespJson))

	return nil
}

func BdevLvolDeleteCmd() cli.Command {
	return cli.Command{
		Name: "delete",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "alias-name",
				Usage: "Optional. The alias of a lvol is <LVSTORE NAME>/<LVOL NAME>. Specify this or uuid.",
			},
			cli.StringFlag{
				Name:  "uuid",
				Usage: "Optional. Specify this or alias-name",
			},
		},
		Usage: "delete a bdev lvol using a block device: \"delete --alias-name <LVSTORE NAME>/<LVOL NAME>\" or \"delete --uuid <UUID>\"",
		Action: func(c *cli.Context) {
			if err := bdevLvolDelete(c); err != nil {
				logrus.WithError(err).Fatalf("Error running delete bdev lvol command")
			}
		},
	}
}

func bdevLvolDelete(c *cli.Context) error {
	spdkCli, err := spdk.NewClient()
	if err != nil {
		return err
	}

	name := c.String("alias-name")
	if name == "" {
		name = c.String("uuid")
	}

	deleted, err := spdkCli.BdevLvolDelete(name)
	if err != nil {
		return err
	}

	bdevLvolDeleteRespJson, err := json.Marshal(deleted)
	if err != nil {
		return err
	}
	fmt.Println(string(bdevLvolDeleteRespJson))

	return nil
}

func BdevLvolGetCmd() cli.Command {
	return cli.Command{
		Name: "get",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "alias-name",
				Usage: "Optional. The alias of a lvol is <LVSTORE NAME>/<LVOL NAME>. If you want to get one specific Lvol info, please input this or uuid.",
			},
			cli.StringFlag{
				Name:  "uuid",
				Usage: "Optional. If you want to get one specific Lvol info, please input this or alias-name",
			},
			cli.Uint64Flag{
				Name:  "timeout, t",
				Usage: "Optional. Determine the timeout of the execution",
				Value: 0,
			},
		},
		Usage: "get all bdev lvol if the info is not specified: \"get\", or \"get --lvs-name <LVOL NAME>\", or \"get --uuid <UUID>\"",
		Action: func(c *cli.Context) {
			if err := bdevLvolGet(c); err != nil {
				logrus.WithError(err).Fatalf("Error running get bdev lvol command")
			}
		},
	}
}

func bdevLvolGet(c *cli.Context) error {
	spdkCli, err := spdk.NewClient()
	if err != nil {
		return err
	}

	name := c.String("alias-name")
	if name == "" {
		name = c.String("uuid")
	}

	bdevLvolGetResp, err := spdkCli.BdevLvolGet(name, c.Uint64("timeout"))
	if err != nil {
		return err
	}

	bdevLvolGetRespJson, err := json.Marshal(bdevLvolGetResp)
	if err != nil {
		return err
	}
	fmt.Println(string(bdevLvolGetRespJson))

	return nil
}
