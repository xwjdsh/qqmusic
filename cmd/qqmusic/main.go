package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/xwjdsh/qqmusic"
)

func main() {
	cookieName := "QQMUSIC_COOKIE"
	cookie := os.Getenv(cookieName)
	if cookie == "" {
		log.Fatalf("%s env is empty", cookieName)
	}
	h := qqmusic.New(cookie)
	app := &cli.App{
		Name:  "qqmusic",
		Usage: "command-line qqmusic tool",
		Commands: []*cli.Command{
			{
				Name:  "singer",
				Usage: "Show singer info",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "order",
						Usage: "Order by [favor|comment]",
						Value: "favor",
					},
					&cli.IntFlag{
						Name:  "count",
						Usage: "Song count",
						Value: 10,
					},
					&cli.BoolFlag{
						Name:  "solo",
						Usage: "Solo",
					},
				},
				Action: func(c *cli.Context) error {
					return singerAction(h, c)
				},
			},
		},
		Action: func(c *cli.Context) error {
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
