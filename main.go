package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yfuruyama/crzerolog"
	"go.opencensus.io/plugin/ochttp"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler).Methods("GET")
	r.HandleFunc("/tts", ttsHandler).Methods("POST")

	rootLogger := zerolog.New(os.Stdout)
	middleware := crzerolog.InjectLogger(&rootLogger)

	handler := cors.Default().Handler(r)
	handler = middleware(handler)
	handler = authMiddleware(handler)

	httpHandler := &ochttp.Handler{
		Propagation: &propagation.HTTPFormat{},
		Handler:     handler,
	}

	log.Info().Msg("Starting server...")

	if err := http.ListenAndServe(":8080", httpHandler); err != nil {
		log.Fatal().Err(err).Msg("Can't start server")
	}

}

func initializeAppDefault() *firebase.App {
	ctx := context.Background()
	config := &firebase.Config{ProjectID: os.Getenv("PROJECT_ID")}
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		fmt.Printf("error initializing app: %v\n", err)
	}
	return app
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeader) != 2 {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		idToken := authHeader[1]
		app := initializeAppDefault()
		ctx := context.Background()
		client, err := app.Auth(ctx)
		if err != nil {
			http.Error(w, "FirebaseError", http.StatusInternalServerError)
			return
		}

		token, err := client.VerifyIDToken(ctx, idToken)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "User not found", http.StatusForbidden)
			return
		}

		if token.Audience != os.Getenv("PROJECT_ID") {
			http.Error(w, "Unauthorized user", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	logger := log.Ctx(r.Context())
	logger.Info().Msg("Request on Health checker")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

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
