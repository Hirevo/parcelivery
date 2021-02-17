package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Parcelivery",
		Usage: "parcel delivery simulation software",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "interactive",
				Aliases:     []string{"i"},
				Usage:       "Enable interactive mode (step-by-step)",
				Required:    false,
				Value:       false,
			},
			&cli.BoolFlag{
				Name:        "help",
				Aliases:     []string{"h"},
				Usage:       "Display this help message",
				Required:    false,
				Value:       false,
			},
		},
		HideHelp: true,
		ArgsUsage: "INPUT_FILE",
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 || c.Bool("help") {
				cli.ShowAppHelp(c)
				return nil
			}

			env, err := parseFromFile(c.Args().First())
			if err != nil {
				return err
			}

			var state = ESContinue
			var count = 1
			var truckEvent Event
			var transportEvents []Event
			if c.Bool("interactive") {
				fmt.Println("Press Enter to step through the turns of the simulation")
				for state == ESContinue {
					fmt.Scanln()
					truckEvent, transportEvents, state = env.NextTurn()
					if state != ESContinue {
						break
					}
					env.DisplayFancy(count, truckEvent, transportEvents)
					fmt.Println()
					count++
				}
			} else {
				for state == ESContinue {
					truckEvent, transportEvents, state = env.NextTurn()
					if state != ESContinue {
						break
					}
					env.DisplayStandard(count, truckEvent, transportEvents)
					fmt.Println()
					count++
				}
			}

			switch state {
			case ESCompleted:
				fmt.Println("😎")
			case ESTurnCountExpired:
				fmt.Println("🙂")
			case ESContinue:
				fmt.Println("😱")
			case ESCantMakeProgress:
				fmt.Println("😱")
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
