package ai

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"google.golang.org/genai"
)

type GeminiService struct {
	client *genai.Client
	model  string
}

func NewGeminiService(ctx context.Context, apiKey string) (*GeminiService, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiService{
		client: client,
		model:  "gemini-2.5-flash-lite",
	}, nil
}

func (s *GeminiService) AnalyzeImage(ctx context.Context, file multipart.File) (bool, string, error) {
	if _, err := file.Seek(0, 0); err != nil {
		return false, "", fmt.Errorf("failed to rewind file: %w", err)
	}

	imgData, err := io.ReadAll(file)
	if err != nil {
		return false, "", fmt.Errorf("failed to read image bytes: %w", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return false, "", fmt.Errorf("failed to reset file pointer: %w", err)
	}

	promptText := `Analyze this image. Does it clearly contain a cup, bottle, or container with water or a transparent liquid being consumed or ready for consumption? 
	Answer EXACTLY with the word "YES" or "NO". If the image is dark, unclear, or shows something else, answer "NO". 
	On the next line, briefly explain why.`
	parts := []*genai.Part{
		{Text: promptText},
		{
			InlineData: &genai.Blob{
				MIMEType: "image/jpeg",
				Data:     imgData,
			},
		},
	}
	contents := []*genai.Content{
		{Parts: parts},
	}

	resp, err := s.client.Models.GenerateContent(ctx, s.model, contents, nil)
	if err != nil {
		return false, "", fmt.Errorf("GenAI API error: %w", err)
	}

	if resp == nil || len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return false, "no response from AI", nil
	}

	part := resp.Candidates[0].Content.Parts[0]
	fullText := part.Text
	if fullText == "" {
		fullText = fmt.Sprintf("%v", part)
	}

	lines := strings.Split(fullText, "\n")
	firstLine := strings.TrimSpace(strings.ToUpper(lines[0]))

	explanation := "No explanation provided"
	if len(lines) > 1 {
		explanation = strings.Join(lines[1:], " ")
	}

	isWater := strings.Contains(firstLine, "YES")

	return isWater, explanation, nil
}
