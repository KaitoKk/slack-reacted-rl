package rl

import (
	"slack-reacted-rl/rl/clients"
	"slack-reacted-rl/rl/models"
	"slack-reacted-rl/rl/scraper"
)

func RL(url string, apiKey string, databaseId string) error {
	title, err := scraper.FetchTitle(url)
	if err != nil {
		return err
	}

	c := clients.NewNotionClient(
		apiKey,
		databaseId,
	)

	article := models.Article{
		Title: title,
		Link:  url,
	}

	err = c.PostArticle(article)
	if err != nil {
		return err
	}
	return nil
}
