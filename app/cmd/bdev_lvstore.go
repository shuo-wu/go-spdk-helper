package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/longhorn/go-spdk-helper/spdk"
)

func BdevLvstoreCmd() cli.Command {
	return cli.Command{
		Name:      "bdev-lvstore",
		ShortName: "lvs",
		Subcommands: []cli.Command{
			BdevLvstoreCreateCmd(),
			BdevLvstoreDeleteCmd(),
			BdevLvstoreGetCmd(),
		},
	}
}

func BdevLvstoreCreateCmd() cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "create a bdev lvstore based on a block device: \"create <BDEV NAME> <LVSTORE NAME>\"",
		Action: func(c *cli.Context) {
			if err := bdevLvstoreCreate(c); err != nil {
				logrus.WithError(err).Fatalf("Error running create bdev lvstore command")
			}
		},
	}
}

func bdevLvstoreCreate(c *cli.Context) error {
	spdkCli, err := spdk.NewClient()
	if err != nil {
		return err
	}

	uuid, err := spdkCli.BdevLvolCreateLvstore(c.Args().First(), c.Args().Get(1))
	if err != nil {
		return err
	}

	bdevLvstoreCreateRespJson, err := json.Marshal(map[string]string{"uuid": uuid})
	if err != nil {
		return err
	}
	fmt.Println(string(bdevLvstoreCreateRespJson))

	return nil
}

func BdevLvstoreDeleteCmd() cli.Command {
	return cli.Command{
		Name: "delete",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "lvs-name",
				Usage: "Optional. Specify this or uuid",
			},
			cli.StringFlag{
				Name:  "uuid",
				Usage: "Optional. Specify this or lvs-name",
			},
		},
		Usage: "delete a bdev lvstore using a block device: \"delete --lvs-name <LVSTORE NAME>\" or \"delete --uuid <UUID>\"",
		Action: func(c *cli.Context) {
			if err := bdevLvstoreDelete(c); err != nil {
				logrus.WithError(err).Fatalf("Error running delete bdev lvstore command")
			}
		},
	}
}

func bdevLvstoreDelete(c *cli.Context) error {
	spdkCli, err := spdk.NewClient()
	if err != nil {
		return err
	}

	deleted, err := spdkCli.BdevLvolDeleteLvstore(c.String("lvs-name"), c.String("uuid"))
	if err != nil {
		return err
	}

	bdevLvstoreDeleteRespJson, err := json.Marshal(deleted)
	if err != nil {
		return err
	}
	fmt.Println(string(bdevLvstoreDeleteRespJson))

	return nil
}

func BdevLvstoreGetCmd() cli.Command {
	return cli.Command{
		Name: "get",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "lvs-name",
				Usage: "Optional. If you want to get one specific Lvstore info, please input this or uuid",
			},
			cli.StringFlag{
				Name:  "uuid",
				Usage: "Optional. If you want to get one specific Lvstore info, please input this or lvs-name",
			},
		},
		Usage: "get all bdev lvstore if the info is not specified: \"get\", or \"get --lvs-name <LVSTORE NAME>\", or \"get --uuid <UUID>\"",
		Action: func(c *cli.Context) {
			if err := bdevLvstoreGet(c); err != nil {
				logrus.WithError(err).Fatalf("Error running get bdev lvstore command")
			}
		},
	}
}

func bdevLvstoreGet(c *cli.Context) error {
	spdkCli, err := spdk.NewClient()
	if err != nil {
		return err
	}

	bdevLvstoreGetResp, err := spdkCli.BdevLvolGetLvstore(c.String("lvs-name"), c.String("uuid"))
	if err != nil {
		return err
	}

	bdevLvstoreGetRespJson, err := json.Marshal(bdevLvstoreGetResp)
	if err != nil {
		return err
	}
	fmt.Println(string(bdevLvstoreGetRespJson))

	return nil
}
