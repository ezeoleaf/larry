package main

import (
	"fmt"
	"log"

	"github.com/ezeoleaf/GobotTweet/config"
)

func tweetContent(cfg config.Config, content string) bool {
	if cfg.SafeMode {
		log.Print("Running in Safe Mode")
		log.Print(content)
		return true
	}

	client := cfg.AccessCfg.GetTwitterClient()
	_, _, err := client.Statuses.Update(content, nil)

	if err != nil {
		log.Print(err)
		return false
	}

	fmt.Println("Content Published")
	return true
}
