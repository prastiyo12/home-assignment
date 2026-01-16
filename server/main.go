package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients = make(map[string]*websocket.Conn)
	mu      sync.Mutex
	apiKey  string
)

func main() {
	_ = godotenv.Load()

	apiKey = mustEnv("API_TOKEN")
	port := mustEnv("SERVER_PORT")

	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/trigger", triggerHandler)

	log.Println("Server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Missing env: %s", k)
	}
	return v
}

func authorized(r *http.Request) bool {
	return r.Header.Get("Authorization") == "Bearer "+apiKey
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	if !authorized(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	clientID := r.Header.Get("X-Client-ID")
	if clientID == "" {
		http.Error(w, "missing client id", 400)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	mu.Lock()
	clients[clientID] = conn
	mu.Unlock()

	log.Println("Client connected:", clientID)
}

func triggerHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	for id, c := range clients {
		c.WriteMessage(websocket.TextMessage, []byte("DOWNLOAD"))
		log.Println("Download command sent to:", id)
	}

	w.Write([]byte("Trigger sent"))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if !authorized(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	clientID := r.Header.Get("X-Client-ID")
	if clientID == "" {
		http.Error(w, "missing client id", 400)
		return
	}

	file, err := os.Create("received_" + clientID + ".txt")
	if err != nil {
		http.Error(w, "file error", 500)
		return
	}
	defer file.Close()

	io.Copy(file, r.Body)
	log.Println("File received from:", clientID)

	w.WriteHeader(http.StatusOK)
}
