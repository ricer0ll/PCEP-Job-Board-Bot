package utils

import "testing"

func TestGetDiscordChannelID(t *testing.T) {
	t.Setenv("CHANNEL_ID", "12345")

	token := GetDiscordChannelID()

	if token == "" {
		t.Errorf("Token is not supposed to be empty")
	}

	tokenSecond := GetDiscordChannelID()

	if tokenSecond == "" {
		t.Errorf("Token is not supposed to be empty")
	}
}

func TestGetDiscordChannelIDFail(t *testing.T) {
	token := GetDiscordChannelID()

	if token != "" {
		t.Errorf("Token is supposed to be empty")
	}
}
