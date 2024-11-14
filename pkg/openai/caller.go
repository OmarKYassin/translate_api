package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Caller interface {
	requestTranslation(string) (translation, error)
}

type translation struct {
	Sentences []IndexSentence `json:"translations"`
}

type OpenAICaller struct{}

var translationResponseSchema = generateSchema[translation]()

func Client() *openai.Client {
	clientOnce.Do(initClient)
	return client
}

func initClient() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("OPENAI_API_KEY environment variable is required")
	}
	client = openai.NewClient(option.WithAPIKey(apiKey))
	return
}

// Makes the OpenAI API call and returns the translation response
func (c OpenAICaller) requestTranslation(prompt string) (translation, error) {
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:   openai.F("IndexSentenceSchema"),
		Schema: openai.F(translationResponseSchema),
		Strict: openai.Bool(true),
	}

	responseFormat := openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
		openai.ResponseFormatJSONSchemaParam{
			Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
			JSONSchema: openai.F(schemaParam),
		},
	)

	chat, err := Client().Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			}),
			ResponseFormat: responseFormat,
			Model:          openai.F(openai.ChatModelGPT4o2024_08_06),
		},
	)
	if err != nil {
		return translation{}, fmt.Errorf("Failed to create chat completion: %w", err)
	}

	var translatedSentences translation
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &translatedSentences)
	if err != nil {
		return translation{}, fmt.Errorf("Failed to parse OpenAI response: %w", err)
	}
	return translatedSentences, nil
}

// JSON Schema generator for structured output
func generateSchema[T any]() interface{} {
	var v T
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	return reflector.Reflect(v)
}
