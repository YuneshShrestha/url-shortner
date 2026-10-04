package main

import (
	"YuneshShrestha/url-shortner/utils"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var ctx = context.Background()

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found")
	}
	// Changing this secret changes every short code, so set it once and keep it.
	if os.Getenv("URL_SHORTENER_SECRET") == "" {
		fmt.Println("URL_SHORTENER_SECRET must be set")
		return
	}

	dbClient := utils.NewRedisClient()

	if dbClient == nil {
		fmt.Println("Failed to connect to Redis")
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Server index page
		fmt.Fprintln(w, "Welcome to URL Shortener!")
	})

	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		// Shorten URL
		url := r.FormValue("url")

		fmt.Println("Payload: ", url)

		id, err := utils.NextID(ctx, dbClient)
		if err != nil {
			http.Error(w, "Failed to generate short URL", http.StatusInternalServerError)
			return
		}
		shortUrl, err := utils.EncodeID(id)
		if err != nil {
			http.Error(w, "Failed to generate short URL", http.StatusInternalServerError)
			return
		}

		fullShortUrl := fmt.Sprintf("http://localhost:8080/r/%s", shortUrl)

		// Generated short URL
		fmt.Printf("Generated short URL: %s\n", fullShortUrl)

		if err := utils.SetLongURL(ctx, dbClient, id, url); err != nil {
			http.Error(w, "Failed to save short URL", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/r/{code}", func(w http.ResponseWriter, r *http.Request) {
		// Get short code from URL
		code := r.PathValue("code")

		// Find original URL in Redis
		longURL, err := utils.GetLongURL(
			ctx,
			dbClient,
			code,
		)

		if err != nil {
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		}

		// Redirect to original URL
		http.Redirect(
			w,
			r,
			longURL,
			http.StatusPermanentRedirect,
		)
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
