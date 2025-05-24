package zaipra

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/tmc/langchaingo/llms"
)

type Info struct {
	Title       string
	Description string
	Content     any
}

// Answer is the main entry point of Zaipra
// It performs Info Routing and Prompt Diffusion
func Answer(userQusetion, systemPrompt string, infos []Info, llm llms.Model, options ...llms.CallOption) (string, error) {
	// Step 1: Info Routing
	usedIndexes, err := classify(userQusetion, infos, llm)
	if err != nil {
		return "", err
	}
	selectedInfos := make([]Info, 0, len(usedIndexes))
	for _, i := range usedIndexes {
		if i >= 0 && i < len(infos) {
			selectedInfos = append(selectedInfos, infos[i])
		}
	}

	// Step 2: Prompt Diffusion (simple combine for now)
	return generateAnswer(userQusetion, systemPrompt, selectedInfos, llm, options...)
}

// classify simulates selecting relevant info indexes
func classify(userQusetion string, infos []Info, llm llms.Model) ([]int, error) {
	infoOptions := ""
	for i, info := range infos {
		infoName := info.Title
		if info.Description != "" {
			infoName += fmt.Sprintf("（%s）", info.Description)
		}
		infoOptions += fmt.Sprintf("- %s：代號%d\n", infoName, i)
	}
	prompt := fmt.Sprintf("問題： %s\n\n請判斷回答以上問題需要哪些資訊：%s\n\n只回答代號，不回答其它文字。\n如果需要多項資訊，用半形逗號隔開。\n如果不需要任何資訊，或是現有資料不足以回答問題，請回傳-1。", userQusetion, infoOptions)
	llmResponse, err := llms.GenerateFromSinglePrompt(context.Background(), llm, prompt)
	if err != nil {
		return nil, err
	}
	// Parse the response to get indexes
	indexes := strings.Split(strings.TrimSpace(llmResponse), ",")
	usedIndexes := make([]int, 0, len(indexes))
	for _, index := range indexes {
		if index == "-1" {
			return []int{}, nil
		}
		i, err := strconv.Atoi(strings.TrimSpace(index))
		if err != nil {
			return nil, err
		}
		usedIndexes = append(usedIndexes, i)
	}
	return usedIndexes, nil
}

// generateAnswer simulates generating the final answer
func generateAnswer(userQusetion, systemPrompt string, infos []Info, llm llms.Model, options ...llms.CallOption) (string, error) {
	prompt := `user:` + userQusetion + "\n\n以下是相關資訊：\n"
	for i, info := range infos {
		prompt += fmt.Sprintf("資訊 %d：%s\n", i+1, info.Title)
		if info.Description != "" {
			prompt += fmt.Sprintf("描述：%s\n", info.Description)
		}
		if info.Content != nil {
			prompt += fmt.Sprintf("內容：%v\n", info.Content)
		}
	}
	prompt += `
system:` + systemPrompt + `
	
請根據以上資訊以及system指示操作。`
	llmResponse, err := llms.GenerateFromSinglePrompt(context.Background(), llm, prompt, options...)
	return llmResponse, err
}
