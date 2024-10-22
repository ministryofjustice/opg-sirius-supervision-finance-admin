package server

import (
	"io"
	"log"
	"net/http"
)

func downloadProxy(client ApiClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := getContext(r)
		filename := r.URL.Query().Get("filename")
		resp, err := client.Download(ctx, filename)
		if err != nil {
			log.Printf("Error calling download API: %v", err)
			http.Error(w, "Failed to stream file", http.StatusInternalServerError)
		}

		defer resp.Body.Close()

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(resp.StatusCode)

		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Printf("Error streaming response from API: %v", err)
			http.Error(w, "Failed to stream file", http.StatusInternalServerError)
		}
	})
}
