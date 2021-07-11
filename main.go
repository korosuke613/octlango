package main

import (
	"encoding/json"
	"fmt"
	"github.com/korosuke613/octlango/core"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"runtime/debug"
)

func action(c *cli.Context) error {
	octclient := core.NewOctclient(c.String("username"), c.String("token"))
	query, err := octclient.GetRepositoriesContributedTo(
		c.Context,
		c.Bool("sort-by-size"),
		c.Bool("reverse-order"),
	)

	if err != nil {
		return err
	}

	result, err := json.MarshalIndent(query, "", "  ")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", result)
	return nil
}

var (
	version string
)

func Version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "(devel)"
	}
	return info.Main.Version
}

func main() {
	if version == "" {
		version = Version()
	}

	isGlobalRequire := true
	switch os.Args[1] {
	case "version":
	case "help":
		isGlobalRequire = false
	}

	app := &cli.App{
		Name:    "octlango",
		Usage:   "CLI to get statistics on languages used on GitHub.",
		Action:  action,
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "username",
				Aliases:  []string{"u"},
				Usage:    "your `GITHUB_USERNAME`",
				EnvVars:  []string{"OCTLANGO_GH_USERNAME"},
				Required: isGlobalRequire,
			},
			&cli.StringFlag{
				Name:     "token",
				Aliases:  []string{"t"},
				Usage:    "your `GITHUB_TOKEN`",
				EnvVars:  []string{"OCTLANGO_GH_TOKEN", "GITHUB_TOKEN"},
				Required: isGlobalRequire,
			},
			&cli.BoolFlag{
				Name:    "sort-by-size",
				Aliases: []string{"s"},
				Usage:   "if true, the order is by size.",
				Value:   true,
			},
			&cli.BoolFlag{
				Name:    "reverse-order",
				Aliases: []string{"r"},
				Usage:   "if true, reverse the result.",
				Value:   false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "version",
				Usage:   "Print octlango version",
				Action: func(c *cli.Context) error {
					fmt.Printf("octlango version %s\n", version)
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
