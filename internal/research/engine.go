package research

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

//go:embed swot_prompt.txt
var swotPromptTemplate string

type Engine struct {
	llm    *openai.LLM
	tavily TavilyTool
}

func NewEngine(openaiKey, tavilyKey, baseURL, model string) (*Engine, error) {
	opts := []openai.Option{
		openai.WithToken(openaiKey),
		openai.WithModel(model),
	}

	if baseURL != "" {
		opts = append(opts, openai.WithBaseURL(baseURL))
	}

	llm, err := openai.New(opts...)
	if err != nil {
		return nil, err
	}

	return &Engine{
		llm: llm,
		tavily: TavilyTool{
			APIKey: tavilyKey,
		},
	}, nil
}

func (e *Engine) Analyze(ctx context.Context, companyName string) (string, error) {
	query := fmt.Sprintf("company %s overview strengths weaknesses opportunities threats news 2024 2025", companyName)
	searchResult, err := e.tavily.Call(ctx, query)
	if err != nil {
		return "", fmt.Errorf("search error: %w", err)
	}

	fmt.Printf("Debug: Tavily result length for %s: %d\n", companyName, len(searchResult))

	prompt := prompts.NewPromptTemplate(
		swotPromptTemplate,
		[]string{"company", "info"},
	)

	chain := chains.NewLLMChain(e.llm, prompt)

	result, err := chains.Predict(ctx, chain, map[string]any{
		"company": companyName,
		"info":    searchResult,
	})
	if err != nil {
		return "", fmt.Errorf("analysis error: %w", err)
	}

	return strings.TrimSpace(result), nil
}
