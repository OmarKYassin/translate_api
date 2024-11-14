package openai

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/OmarKYassin/translate_api/pkg/types"
	"github.com/stretchr/testify/assert"
)

type mockCaller struct{}

func (m mockCaller) requestTranslation(prompt string) (translation, error) {
	if strings.Contains(prompt, "فشل") {
		return translation{}, fmt.Errorf("Test error")
	}

	return translation{
		Sentences: []IndexSentence{{
			Index:    0,
			Sentence: "Hello",
		},
		},
	}, nil
}

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
		Transcript:   tran,
		OpenAICaller: mockCaller{},
	}

	prompt, containsArabic := translator.buildPrompts()

	assert.Contains(t, prompt[0], "Translate the following sentences to English:\n")
	assert.Contains(t, prompt[0], "0: مرحبا\n")
	assert.NotContains(t, prompt[0], "1: hello\n")
	assert.True(t, containsArabic)
}

func TestBuildPrompt10k(t *testing.T) {
	var tran types.Transcript
	data, err := os.ReadFile("../../tests/fixtures/realistic_arabic_conversation.json")
	if err != nil {
		t.Logf("Failed in reading the fixture file, error: %+v", err)
		t.Fail()
		return
	}
	err = json.Unmarshal(data, &tran)
	if err != nil {
		t.Logf("Failed in marhsaling the fixture data, error: %+v", err)
		t.Fail()
		return
	}

	translator := &OpenAITranslator{
		Transcript:   tran,
		OpenAICaller: mockCaller{},
	}

	prompts, containsArabic := translator.buildPrompts()

	assert.True(t, containsArabic)
	assert.Contains(t, prompts[0], "Translate the following sentences to English:\n")
	assert.Contains(t, prompts[0], "0: انت قولتلي إننا نشتغل سوا، وما نخونش بعض.")
	assert.Contains(t, prompts[0], "75: لو رجعت دلوقتي، مش هعرف أبص في وشي في المراية. مشكلتي مش معاك يا عشري، مشكلتي هنا.")
	assert.Contains(t, prompts[1], "Translate the following sentences to English:\n")
	assert.Contains(t, prompts[1], "76: قرارك واضح يا إبراهيم، خلاص ما فيش كلام تاني.")
	assert.Contains(t, prompts[1], "97: خلاص، يا إبراهيم. شكلك فعلاً اخترت طريقك، واللي فيها فيها. ربنا يسترنا من اللي جاي.")
	assert.Equal(t, len(prompts), 2)
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
				Transcript:   test.input,
				OpenAICaller: mockCaller{},
			}
			translator.Translate()
			assert.Equal(t, test.expected, translator.Transcript)
		})
	}
}

func TestTranslateFailure(t *testing.T) {
	t.Run("Caller returns an error", func(t *testing.T) {
		translator := &OpenAITranslator{
			Transcript: types.Transcript{
				{Speaker: "A", Time: "00:01", Sentence: "فشل"},
			},
			OpenAICaller: mockCaller{},
		}
		err := translator.Translate()
		assert.Equal(t, err.Error(), "failed to translate chunk: Test error")
	})
}
