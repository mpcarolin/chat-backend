package mock

import (
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"chat-backend/internal/chat"
)

type QAPair struct {
	Question string
	Answer   string
}

type MockChatProvider struct {
	qaData []QAPair
}

func NewMockChatProvider() *MockChatProvider {
	provider := &MockChatProvider{}
	provider.loadTSVData()
	return provider
}

func (m *MockChatProvider) loadTSVData() {
	file, err := os.Open("sample-data.tsv")
	if err != nil {
		slog.Warn("Failed to load sample-data.tsv, using fallback responses", "error", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'

	records, err := reader.ReadAll()
	if err != nil {
		slog.Warn("Failed to parse sample-data.tsv", "error", err)
		return
	}

	// Skip header row
	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) >= 2 {
			m.qaData = append(m.qaData, QAPair{
				Question: record[0],
				Answer:   record[1],
			})
		}
	}

	slog.Info("Loaded TSV data", "pairs", len(m.qaData))
}

func (m *MockChatProvider) Chat(ctx context.Context, req *chat.ChatRequest) (*chat.ChatResponse, error) {
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("no messages provided")
	}

	question := req.Messages[len(req.Messages)-1].Content

	if question == "" {
		return &chat.ChatResponse{
			Content: "Please provide a question.",
		}, nil
	}

	if len(question) < 3 {
		return &chat.ChatResponse{
			Content: "Your question seems too short. Could you provide more details?",
		}, nil
	}

	// If no TSV data loaded, use fallback
	if len(m.qaData) == 0 {
		return &chat.ChatResponse{
			Content: fmt.Sprintf("This is a mock response to your question: '%s'. In a real implementation, this would come from a chat service.", question),
		}, nil
	}

	// Find best match using longest common substring
	bestMatch := m.findBestMatch(question)
	if bestMatch != nil {
		return &chat.ChatResponse{
			Content: bestMatch.Answer,
		}, nil
	}

	return &chat.ChatResponse{
		Content: "I don't have an answer for that question in my knowledge base.",
	}, nil
}

func (m *MockChatProvider) findBestMatch(userQuestion string) *QAPair {
	var bestMatch *QAPair
	maxMatchLength := 0

	userQuestionLower := strings.ToLower(userQuestion)

	for i := range m.qaData {
		qa := &m.qaData[i]
		qaQuestionLower := strings.ToLower(qa.Question)

		matchLength := longestCommonSubstring(userQuestionLower, qaQuestionLower)

		if matchLength > maxMatchLength {
			maxMatchLength = matchLength
			bestMatch = qa
		}
	}

	// Only return match if it's substantial enough (at least 3 characters)
	if maxMatchLength >= 3 {
		return bestMatch
	}

	return nil
}

func longestCommonSubstring(s1, s2 string) int {
	maxLen := 0

	for i := range len(s1) {
		for j := i + 1; j <= len(s1); j++ {
			substring := s1[i:j]
			if strings.Contains(s2, substring) {
				if len(substring) > maxLen {
					maxLen = len(substring)
				}
			}
		}
	}

	return maxLen
}
