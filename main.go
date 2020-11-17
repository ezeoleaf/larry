package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Tweet Random Repo",
		Usage: "Tweet random repositories",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
