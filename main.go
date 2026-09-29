package main

import (
	"YuneshShrestha/url-shortner/utils"
	"context"
	"fmt"
	"net/http"
)

var ctx = context.Background()

func main() {
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

		shortUrl := utils.GetShortCode()

		fullShortUrl := fmt.Sprintf("http://localhost:8080/r/%s", shortUrl)

		// Generated short URL
		fmt.Printf("Generated short URL: %s\n", fullShortUrl)

		// Set the key in Redis
		utils.SetKey(&ctx, dbClient, shortUrl, url, 0)
	})

	http.HandleFunc("/r/{code}", func(w http.ResponseWriter, r *http.Request) {
		// Get short code from URL
		code := r.PathValue("code")

		// Find original URL in Redis
		longURL, err := utils.GetLongURL(
			&ctx,
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
