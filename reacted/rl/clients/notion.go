package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slack-reacted-rl/rl/models"
)

const API_BASE_URL = "https://api.notion.com/v1/"

type NotionConfig struct {
	apiKey     string
	databaseId string
}

type NotionClient struct {
	config NotionConfig
	client *http.Client
}

func NewNotionClient(apiKey string, databaseId string) *NotionClient {
	return &NotionClient{
		config: NotionConfig{
			apiKey:     apiKey,
			databaseId: databaseId,
		},
		client: &http.Client{},
	}
}

func (c NotionClient) buildRequest(method string, path string, body io.Reader) *http.Request {
	uri := API_BASE_URL + path

	req, _ := http.NewRequest(method, uri, body)
	req.Header.Set("Authorization", "Bearer "+c.config.apiKey)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	return req
}

/*
 * MEMO: 使ってないけど、とりあえず残しておく
 * 初回とかに疎通確認に使えると良いかも
 */
func (c NotionClient) GetDatabase() {
	path := "databases/" + c.config.databaseId
	req := c.buildRequest("GET", path, nil)

	res, err := c.client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	fmt.Println(string(body))
}

type PostArticleRequest struct {
	Parent     Parent     `json:"parent"`
	Properties Properties `json:"properties"`
}

type Properties struct {
	Title Title `json:"title,omitempty"`
	Link  Link  `json:"link,omitempty"`
}

type Title struct {
	Title []Text `json:"title"`
}

type Link struct {
	Url string `json:"url"`
}

type Text struct {
	Content Content `json:"text"`
}

type Content struct {
	Content string `json:"content"`
}

type Parent struct {
	DatabaseId string `json:"database_id"`
}

func (c NotionClient) buildPostArticleBody(title string, link string) PostArticleRequest {
	return PostArticleRequest{
		Parent: Parent{
			DatabaseId: c.config.databaseId,
		},
		Properties: Properties{
			Title: Title{
				[]Text{
					{Content: Content{Content: title}},
				},
			},
			Link: Link{
				Url: link,
			},
		},
	}
}

func (c NotionClient) PostArticle(article models.Article) error {
	reqBody := c.buildPostArticleBody(article.Title, article.Link)
	jsonBody, _ := json.Marshal(reqBody)

	path := "pages"
	req := c.buildRequest("POST", path, bytes.NewBuffer(jsonBody))

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	log.Println("status code: ", res.StatusCode)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	log.Println(string(body))
	return nil
}
