package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ca17/teamsacs/assets"
	"github.com/ca17/teamsacs/teamsctl/settings"
	"github.com/urfave/cli/v2"
	_ "github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{Name: "build-version", Aliases: []string{"bv"}, Usage: "build version info"},
	}
	app.Name = "TeamsCtl"
	app.Usage = "TeamsCtl Usage"
	app.Description = "teamscli is teamsacs command tools" // 描述
	app.Version = "1.0.1"                                  // 版本
	app.Action = func(c *cli.Context) error {
		fmt.Println(assets.BuildInfo)
		return cli.ShowAppHelp(c)
	}

	app.Commands = append(app.Commands, settings.Commands...)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
