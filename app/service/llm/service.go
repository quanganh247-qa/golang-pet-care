package llm

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type OllamaRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float64 `json:"temperature"`
	Stream      bool    `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

type SOAPService struct {
	ollamaURL string
}

func NewSOAPService(ollamaURL string) *SOAPService {
	return &SOAPService{
		ollamaURL: ollamaURL,
	}
}

func CallOllamaAPI(req *OllamaRequest) (string, error) {
	ollamaURL := "http://localhost:11434/api/generate" // Configure as needed
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(ollamaURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", err
	}

	return ollamaResp.Response, nil
}
