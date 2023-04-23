package gpt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type GptClient struct {
	ApiKey string
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatBody struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatChoice struct {
	Index   int64
	Message ChatMessage
}

type ChatResponose struct {
	Choices []ChatChoice
}

const chatUrl = "https://api.openai.com/v1/chat/completions"

var Handler HttpClient

func init() {
	Handler = &http.Client{}
}

func NewGptClient() *GptClient {
	apiKey := os.Getenv("GPT_API_KEY")
	if apiKey == "" {
		log.Fatal().Msg("required GPT_API_KEY env var is missing")
	}
	return &GptClient{ApiKey: apiKey}
}

func (gpt *GptClient) GenerateResponse(msg string) string {
	req, err := http.NewRequest("POST", chatUrl, strings.NewReader(generateBody(msg)))
	if err != nil {
		log.Err(err).Msgf("[Error] while building request")
	}
	bearerToken := fmt.Sprintf("Bearer %s", gpt.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken)

	log.Printf("%+v", req)

	var chat ChatResponose
	resp, err := Handler.Do(req)
	//DebugResponse(resp)
	if err != nil {
		log.Err(err).Msg("[Error] while querying gpt api")
	}

	err = json.NewDecoder(resp.Body).Decode(&chat)
	if err != nil {
		log.Err(err).Msg("[Error] while decoding the response")
	}

	return chat.Choices[0].Message.Content
}

func generateBody(msg string) string {
	cMessage := ChatMessage{Role: "user", Content: msg}
	cBody := ChatBody{Model: "gpt-3.5-turbo", Messages: []ChatMessage{cMessage}}
	data, err := json.Marshal(cBody)
	if err != nil {
		log.Err(err).Msgf("[Error] while encoding body")
	}
	return string(data)
}

func DebugResponse(resp *http.Response) {
	c := resp
	b, err := io.ReadAll(c.Body)
	if err != nil {
		log.Printf("error while reading body, %s", err)
	}
	log.Printf("got response %s", string(b))
}
