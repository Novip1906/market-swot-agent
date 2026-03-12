package research

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tmc/langchaingo/tools"
)

type TavilyTool struct {
	APIKey string
}

var _ tools.Tool = TavilyTool{}

func (t TavilyTool) Name() string {
	return "tavily_search"
}

func (t TavilyTool) Description() string {
	return "A search engine optimized for LLMs. Use this to get the latest news and information about companies and products. Input should be a search query."
}

type tavilyRequest struct {
	Query       string `json:"query"`
	APIKey      string `json:"api_key"`
	SearchDepth string `json:"search_depth"`
	MaxResults  int    `json:"max_results"`
}

type tavilyResponse struct {
	Results []struct {
		Title   string `json:"title"`
		URL     string `json:"url"`
		Content string `json:"content"`
	} `json:"results"`
}

func (t TavilyTool) Call(ctx context.Context, input string) (string, error) {
	reqBody, _ := json.Marshal(tavilyRequest{
		Query:       input,
		APIKey:      t.APIKey,
		SearchDepth: "advanced",
		MaxResults:  5,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.tavily.com/search", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res tavilyResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	var output string
	for _, r := range res.Results {
		output += fmt.Sprintf("Title: %s\nURL: %s\nContent: %s\n\n", r.Title, r.URL, r.Content)
	}

	if output == "" {
		return "No results found.", nil
	}

	return output, nil
}
