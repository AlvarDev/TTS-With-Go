package main

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())
	logger.Info().Msg("Request on Health checker")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
