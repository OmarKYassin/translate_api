package openai

import (
	"testing"

	"github.com/OmarKYassin/translate_api/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestHasArabicLetters(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"مرحبا", true},
		{"hello", false},
		{"مرحبا hello", true},
		{"", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := hasArabicLetters(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestBuildPrompt(t *testing.T) {
	tran := types.Transcript{
		{Speaker: "A", Time: "00:01", Sentence: "مرحبا"},
		{Speaker: "B", Time: "00:02", Sentence: "hello"},
	}

	translator := &OpenAITranslator{
		Transcript: tran,
	}

	prompt, containsArabic := translator.buildPrompt()

	assert.Contains(t, prompt, "Translate the following sentences to English:\n")
	assert.Contains(t, prompt, "0: مرحبا\n")
	assert.NotContains(t, prompt, "1: hello\n")
	assert.True(t, containsArabic)
}

func TestTranslate(t *testing.T) {
	tests := []struct {
		name     string
		input    types.Transcript
		expected types.Transcript
	}{
		{
			name: "Contain an arabic sentence",
			input: types.Transcript{
				{Speaker: "A", Time: "00:01", Sentence: "مرحبا"},
				{Speaker: "B", Time: "00:02", Sentence: "hello"},
			},
			expected: types.Transcript{
				{Speaker: "A", Time: "00:01", Sentence: "Hello"},
				{Speaker: "B", Time: "00:02", Sentence: "hello"},
			},
		},
		{
			name: "Doesn't an arabic sentence",
			input: types.Transcript{
				{Speaker: "A", Time: "00:01", Sentence: "Hi"},
				{Speaker: "B", Time: "00:02", Sentence: "hello"},
			},
			expected: types.Transcript{
				{Speaker: "A", Time: "00:01", Sentence: "Hi"},
				{Speaker: "B", Time: "00:02", Sentence: "hello"},
			},
		},
		{
			name: "Has empty strings",
			input: types.Transcript{
				{Speaker: "A", Time: "00:01", Sentence: "مرحبا"},
				{Speaker: "B", Time: "00:02", Sentence: ""},
			},
			expected: types.Transcript{
				{Speaker: "A", Time: "00:01", Sentence: "Hello"},
				{Speaker: "B", Time: "00:02", Sentence: ""},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			translator := &OpenAITranslator{
				Transcript: test.input,
			}
			translator.Translate()
			assert.Equal(t, test.expected, translator.Transcript)
		})
	}
}
