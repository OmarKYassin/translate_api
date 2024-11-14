package openai

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/OmarKYassin/translate_api/pkg/types"
	"github.com/openai/openai-go"
)

type IndexSentence struct {
	Index    int    `json:"index" jsonschema_description:"The index of the sentence"`
	Sentence string `json:"sentence" jsonschema_description:"The sentence"`
}

type OpenAITranslator struct {
	Transcript   types.Transcript
	OpenAICaller Caller
}

var (
	clientOnce sync.Once
	client     *openai.Client
)

func (t *OpenAITranslator) Translate() error {
	prompts, containsArabic := t.buildPrompts()
	if !containsArabic {
		return nil
	}
	var totalTranslatedSentences []IndexSentence

	for _, prompt := range prompts {
		translatedSentences, err := t.OpenAICaller.requestTranslation(prompt)
		if err != nil {
			return fmt.Errorf("failed to translate chunk: %w", err)
		}

		totalTranslatedSentences = append(
			totalTranslatedSentences,
			translatedSentences.Sentences...,
		)
	}

	// Update original transcript with translated sentences
	for _, indexSentencePair := range totalTranslatedSentences {
		t.Transcript[indexSentencePair.Index].Sentence = indexSentencePair.Sentence
	}

	return nil
}

// Builds prompt and checks if there's any Arabic content
func (t *OpenAITranslator) buildPrompts() ([]string, bool) {
	containsArabic := false
	var prompts []string
	var currentPrompt strings.Builder
	currentPrompt.WriteString("Translate the following sentences to English:\n")
	currentLength := currentPrompt.Len()

	for idx, entry := range t.Transcript {
		if !hasArabicLetters(entry.Sentence) {
			continue
		}
		containsArabic = true
		entryText := fmt.Sprintf("%d: %s\n", idx, entry.Sentence)
		entryLength := len(entryText)

		// If adding the entry exceeds the limit, finalize the current prompt and start a new one
		if currentLength+entryLength > 8000 { // Adjust limit to stay within API token constraints
			prompts = append(prompts, currentPrompt.String())
			currentPrompt.Reset()
			currentPrompt.WriteString("Translate the following sentences to English:\n")
			currentLength = currentPrompt.Len()
		}

		// Add entry to the current prompt
		currentPrompt.WriteString(entryText)
		currentLength += entryLength
	}

	if !containsArabic {
		return []string{}, containsArabic
	}

	// Add the last prompt if there’s remaining content
	if currentPrompt.Len() > 0 {
		prompts = append(prompts, currentPrompt.String())
	}

	return prompts, containsArabic
}

// Helper function to check for Arabic letters
func hasArabicLetters(s string) bool {
	arabicLetterRegex := regexp.MustCompile(`\p{Arabic}`)
	return arabicLetterRegex.MatchString(s)
}
