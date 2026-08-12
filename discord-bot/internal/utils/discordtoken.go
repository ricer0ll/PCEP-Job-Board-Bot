package utils

import (
	"os"
)

func GetDiscordToken() (*string, error) {
	var token string

	// load env vars
	token = os.Getenv("DISCORD_TOKEN")
	if len(token) != 0 {
		return &token, nil
	}

	return &token, nil
}
