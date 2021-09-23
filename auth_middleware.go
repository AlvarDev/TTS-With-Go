package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go"
)

func initializeAppDefault() (*firebase.App, error) {
	ctx := context.Background()
	config := &firebase.Config{ProjectID: os.Getenv("PROJECT_ID")}
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		fmt.Printf("error initializing app: %v\n", err)
		return nil, err
	}
	return app, nil
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeader) != 2 {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		idToken := authHeader[1]
		app, err := initializeAppDefault()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.Background()
		client, err := app.Auth(ctx)
		if err != nil {
			http.Error(w, "FirebaseError", http.StatusInternalServerError)
			return
		}

		token, err := client.VerifyIDToken(ctx, idToken)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		if token.Audience != os.Getenv("PROJECT_ID") {
			http.Error(w, "Unauthorized user", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
