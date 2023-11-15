package main

import (
	"fmt"
	"github.com/slack-go/slack"
)

func getMessage(channel string, ts string) (string, error) {
	api := slack.New(slackBotToken)
	history, err := api.GetConversationHistory(&slack.GetConversationHistoryParameters{
		ChannelID: channel,
		Latest:    ts,
		Inclusive: true,
	})
	if err != nil {
		return "", err
	}

	if len(history.Messages) == 0 {
		return "", error(fmt.Errorf("message not found"))
	}

	return history.Messages[0].Text, nil
}
