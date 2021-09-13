package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
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

	httpHandler := &ochttp.Handler{
		Propagation: &propagation.HTTPFormat{},
		Handler:     handler,
	}

	log.Info().Msg("Starting server...")

	if err := http.ListenAndServe(":8080", httpHandler); err != nil {
		log.Fatal().Err(err).Msg("Can't start server")
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	logger := log.Ctx(r.Context())
	logger.Info().Msg("Request on Health checker")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func ttsHandler(w http.ResponseWriter, r *http.Request) {

	// logger := log.Ctx(r.Context())
	fmt.Println(r.Body)

	response := make(map[string]interface{})
	response["name"] = "AlvarDev"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(response)

}
