package rl

import (
	"errors"
	"os"
	"slack-reacted-rl/rl/clients"
	"slack-reacted-rl/rl/models"
	"slack-reacted-rl/rl/scraper"
)

var (
	apiKey     = os.Getenv("NOTION_API_KEY")
	databaseId = os.Getenv("NOTION_DATABASE_ID")
)

func rl(url string) error {
	if err := checkVar(); err != nil {
		return err
	}

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

func checkVar() error {
	if apiKey == "" {
		return errors.New("環境変数NOTION_API_KEYがありません")
	}
	if databaseId == "" {
		return errors.New("環境変数NOTION_DATABASE_IDがありません")
	}
	return nil
}
