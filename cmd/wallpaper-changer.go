package main

import (
	"log"
	"os"

	wp "github.com/Quik95/wallpaper-changer"
	"github.com/urfave/cli/v2"
)

func main() {
	type Options struct {
		Categories string
		Purity     string
		Sortign    string
		Order      string
		TopRange   string
		Atleast    string
		Resolution string
		Ratios     string
		Page       uint32
		Seed       string
		Query      string
	}

	app := &cli.App{
		Name:  "Wallpaper Changer",
		Usage: "Set your desktop wallpaper to one of many amazing wallpapers delivered by wallhaven.cc",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "categories",
				Value:   "111",
				Usage:   "Select which categories to fetch from General/Anime/People",
				Aliases: []string{"c"},
			},
			&cli.StringFlag{
				Name:    "purity",
				Value:   "100",
				Usage:   "Select which categories to fetch from General/Anime/People",
				Aliases: []string{"p"},
			},
			&cli.StringFlag{
				Name:    "sorting",
				Value:   "date_added",
				Usage:   "Select sorting method",
				Aliases: []string{"s"},
			},
			&cli.StringFlag{
				Name:  "order",
				Value: "desc",
				Usage: "Select sorting order",
			},
			&cli.StringFlag{
				Name:    "top-range",
				Value:   "1M",
				Usage:   "Select time range for toplist sorting",
				Aliases: []string{"r"},
			},
			&cli.StringFlag{
				Name:    "atleast",
				Value:   "",
				Usage:   "Minimal wallpaper resolution",
				Aliases: []string{"a"},
			},
			&cli.StringFlag{
				Name:  "resolutions",
				Value: "",
				Usage: "List of resolutions to fetch",
			},
			&cli.StringFlag{
				Name:  "ratios",
				Value: "",
				Usage: "List of wallpaper ratios to fetch",
			},
			&cli.IntFlag{
				Name:  "pages",
				Value: 1,
				Usage: "Number of pages to fetch",
			},
			&cli.StringFlag{
				Name:  "seed",
				Value: "",
				Usage: "Seed used to generate random results",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Value: false,
				Usage: "Don't set desktop background. Only download wallpaper",
			},
			&cli.StringFlag{
				Name:    "output",
				Value:   "./",
				Usage:   "Where to save wallpaper",
				Aliases: []string{"o"},
			},
			&cli.StringFlag{
				Name:    "api-key",
				Value:   "",
				Usage:   "Wallhaven api key",
				Aliases: []string{"k"},
			},
		},
		Action: func(c *cli.Context) error {
			if err := wp.ValidateArgs(c); err != nil {
				return err
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
