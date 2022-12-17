package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sha65536/voronoi"
	"github.com/urfave/cli/v2"
)

func main() {
	var numPoints int
	var show bool
	app := &cli.App{
		Name:  "voronoi",
		Usage: "Generates voronoi animation in the terminal!",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "points",
				Aliases:     []string{"p"},
				Value:       8,
				Usage:       "Number of points",
				Destination: &numPoints,
				Action: func(ctx *cli.Context, i int) error {
					if i < 2 {
						return fmt.Errorf("flag points value %v lower than 2", i)
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:        "show",
				Aliases:     []string{"s", "v"},
				Value:       false,
				Usage:       "Show the location of points",
				Destination: &show,
			},
		},
		Action: func(*cli.Context) error {
			anim, err := voronoi.MakeAnimation(numPoints, show)
			if err != nil {
				return err
			}
			anim.Start()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
