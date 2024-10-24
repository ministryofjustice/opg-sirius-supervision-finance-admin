package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func requestDownload(envVars EnvironmentVars) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := r.URL.Query().Get("uid")

		if !IsHxRequest(r) {
			http.Error(w, "Bad routing", http.StatusBadRequest)
			return
		}
		w.Header().Add("HX-Redirect", fmt.Sprintf("%s/download/callback?uid=%s", envVars.Prefix, uid))
	})
}

func downloadCallback(client ApiClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := getContext(r)
		uid := r.URL.Query().Get("uid")
		resp, err := client.Download(ctx, uid)
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
