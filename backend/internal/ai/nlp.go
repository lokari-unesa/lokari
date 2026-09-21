package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/sashabaranov/go-openai"
)

// SummarizeAlert uses any OpenAI-compatible API (OpenRouter, Groq, DeepSeek)
func SummarizeAlert(ctx context.Context, rawData string) (string, error) {
	apiKey := os.Getenv("AI_API_KEY")
	baseURL := os.Getenv("AI_BASE_URL")
	modelName := os.Getenv("AI_MODEL")

	if apiKey == "" || baseURL == "" || modelName == "" {
		return "", fmt.Errorf("konfigurasi AI (API_KEY, BASE_URL, MODEL) belum diatur di .env")
	}

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL
	client := openai.NewClientWithConfig(config)

	prompt := fmt.Sprintf(`Anda adalah asisten mitigasi bencana untuk warga Desa Jarak, Kediri. 
Berikut adalah data mentah satelit aktivitas gunung berapi dari NASA EONET:
%s

Tugas Anda:
1. Rangkum data di atas dalam 1-2 paragraf pendek berbahasa Indonesia yang santai tapi jelas.
2. Jelaskan apakah ada ancaman letusan yang dekat dengan "Gunung Kelud" atau area Jawa Timur.
3. Berikan himbauan mitigasi singkat untuk warga Desa Jarak.
Jangan gunakan bahasa teknis berlebihan.`, rawData)

	req := openai.ChatCompletionRequest{
		Model: modelName,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "Anda adalah asisten mitigasi bencana yang berbicara dengan ramah, berempati, dan sangat mudah dipahami warga.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		if apiErr, ok := err.(*openai.APIError); ok {
			if apiErr.HTTPStatusCode == 429 || apiErr.HTTPStatusCode == 402 {
				return "", fmt.Errorf("KUOTA_OPENROUTER_HABIS")
			}
		}
		return "", fmt.Errorf("AI Provider error: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}

// NewsItem represents the structured JSON returned by AI
type NewsItem struct {
	Category string `json:"kategori"`
	Title    string `json:"judul"`
	Summary  string `json:"ringkasan"`
}

// GenerateNewsSummary uses Cohere API to summarize JSON data into a NewsItem
func GenerateNewsSummary(ctx context.Context, rawData string, source string) (*NewsItem, error) {
	apiKey := os.Getenv("AI_API_KEY")
	baseURL := os.Getenv("AI_BASE_URL")
	modelName := os.Getenv("AI_MODEL")

	if apiKey == "" || baseURL == "" || modelName == "" {
		return nil, fmt.Errorf("konfigurasi AI (API_KEY, BASE_URL, MODEL) belum diatur di .env")
	}

	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL
	client := openai.NewClientWithConfig(config)

	prompt := fmt.Sprintf(`Anda adalah asisten cerdas pemantau kebencanaan untuk Desa Jarak, Kediri.
Berikut adalah data mentah dari %s:
%s

Tugas Anda:
Pahami data tersebut dan buatlah sebuah ringkasan berita singkat berbahasa Indonesia yang mudah dipahami warga awam. 
Fokus pada kewaspadaan, mitigasi, atau peringatan dini yang relevan (jika ada ancaman letusan/gempa, jelaskan bahayanya misal "jauhi bantaran sungai dari lahar dingin", atau "pantau terus info resmi").
JANGAN PANIK jika tidak ada bahaya.
Kategorikan ke salah satu dari: volcano, lahar, evac, weather, warning.

Output HARUS berformat JSON persis seperti ini (tanpa markdown blok):
{
  "kategori": "volcano",
  "judul": "Judul Berita Singkat",
  "ringkasan": "Satu paragraf penjelasan untuk warga."
}`, source, rawData)

	req := openai.ChatCompletionRequest{
		Model: modelName,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("AI Provider error: %v", err)
	}

	var news NewsItem
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &news); err != nil {
		return nil, fmt.Errorf("Gagal parsing JSON dari AI: %v\nResponse: %s", err, resp.Choices[0].Message.Content)
	}

	return &news, nil
}

// GenerateEmbedding uses Cohere free API to generate 1024-dimensional vectors
func GenerateEmbedding(text string) ([]float32, error) {
	apiKey := os.Getenv("COHERE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("COHERE_API_KEY is not set in .env")
	}

	url := "https://api.cohere.com/v1/embed"

	payload := map[string]interface{}{
		"texts":           []string{text},
		"model":           "embed-multilingual-v3.0",
		"input_type":      "search_query",
		"embedding_types": []string{"float"},
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Cohere Request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("KUOTA_COHERE_HABIS")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Cohere API returned status: %d", resp.StatusCode)
	}

	var result struct {
		Embeddings struct {
			Float [][]float32 `json:"float"`
		} `json:"embeddings"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("Cohere Parse error: %v", err)
	}
	
	if len(result.Embeddings.Float) > 0 {
		return result.Embeddings.Float[0], nil
	}
	
	return nil, fmt.Errorf("empty embedding result")
}
