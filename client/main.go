package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Missing env: %s", k)
	}
	return v
}

func main() {
	_ = godotenv.Load()

	wsURL := mustEnv("SERVER_WS")
	uploadURL := mustEnv("SERVER_UPLOAD")
	clientID := mustEnv("CLIENT_ID")
	token := mustEnv("API_TOKEN")
	filePath := mustEnv("FILE_PATH")

	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)
	headers.Set("X-Client-ID", clientID)

	for {
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
		if err != nil {
			log.Println("WS connect failed, retrying...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("Connected to server as", clientID)

		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.Println("WS disconnected, reconnecting...")
				ws.Close()
				break
			}

			if string(msg) == "DOWNLOAD" {
				log.Println("DOWNLOAD command received")
				uploadFile(uploadURL, token, clientID, filePath)
			}
		}
	}
}

func uploadFile(url, token, clientID, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Println("File error:", err)
		return
	}
	defer file.Close()

	req, err := http.NewRequest("POST", url, file)
	if err != nil {
		log.Println("Request error:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Client-ID", clientID)
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Upload failed:", err)
		return
	}
	resp.Body.Close()

	log.Println("File uploaded successfully")
}
