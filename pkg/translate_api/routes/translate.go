package routes

import (
	"net/http"

	"github.com/OmarKYassin/translate_api/pkg/openai"
	"github.com/OmarKYassin/translate_api/pkg/types"
	"github.com/gin-gonic/gin"
)

func translate(c *gin.Context) {
	var tran types.Transcript
	if err := c.ShouldBindJSON(&tran); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parsing parameters", "details": err.Error()})
		return
	}

	translator := &openai.OpenAITranslator{
		Transcript: tran,
	}

	if err := translator.Translate(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Translation failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, translator.Transcript)
}
