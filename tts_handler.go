package main

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func ttsHandler(w http.ResponseWriter, r *http.Request) {

	logger := log.Ctx(r.Context())
	logger.Info().Msg("Request on TTS")

	ms := MessageSpeech{}
	err := json.NewDecoder(r.Body).Decode(&ms)
	if err != nil {
		logger.Error().Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Check body att

	audioB64, err := synthesizeSpeechRequest(
		ms.Message,
		ms.LanguageCode,
		ms.VoiceName,
		ms.Gender,
	)

	if err != nil {
		logger.Error().Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response
	response := make(map[string]interface{})
	response["audioB64"] = audioB64

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}

// Request API
type MessageSpeech struct {
	Gender       string `json:"gender"`
	LanguageCode string `json:"languageCode"`
	Message      string `json:"message"`
	VoiceName    string `json:"voiceName"`
}
