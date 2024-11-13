package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/OmarKYassin/translate_api/pkg/types"
	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type IndexSentence struct {
	Index    int    `json:"index" jsonschema_description:"The index of the sentence"`
	Sentence string `json:"sentence" jsonschema_description:"The sentence"`
}

type translation struct {
	Sentences []IndexSentence `json:"translations"`
}

type OpenAITranslator struct {
	Transcript types.Transcript
}

var (
	clientOnce                sync.Once
	client                    *openai.Client
	translationResponseSchema = generateSchema[translation]()
)

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

func (t *OpenAITranslator) Translate() error {
	prompt, containsArabic := t.buildPrompt()
	if !containsArabic {
		return nil
	}

	translatedSentences, err := t.requestTranslationFromOpenAI(prompt)
	if err != nil {
		return err
	}

	// Update original transcript with translated sentences
	for _, indexSentencePair := range translatedSentences.Sentences {
		t.Transcript[indexSentencePair.Index].Sentence = indexSentencePair.Sentence
	}

	return nil
}

// Builds prompt and checks if there's any Arabic content
func (t *OpenAITranslator) buildPrompt() (string, bool) {
	var prompt strings.Builder
	containsArabic := false
	prompt.WriteString("Translate the following sentences to English:\n")

	for idx, entry := range t.Transcript {
		if hasArabicLetters(entry.Sentence) {
			containsArabic = true
			prompt.WriteString(fmt.Sprintf("%d: %s\n", idx, entry.Sentence))
		}
	}
	return prompt.String(), containsArabic
}

// Makes the OpenAI API call and returns the translation response
func (t *OpenAITranslator) requestTranslationFromOpenAI(prompt string) (translation, error) {
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

// Helper function to check for Arabic letters
func hasArabicLetters(s string) bool {
	arabicLetterRegex := regexp.MustCompile(`\p{Arabic}`)
	return arabicLetterRegex.MatchString(s)
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
