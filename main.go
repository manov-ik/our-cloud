package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Serve HTML and static files
	http.Handle("/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/files", listHandler)
	http.HandleFunc("/download", downloadHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File error", 500)
		return
	}
	defer file.Close()

	dst, err := os.Create(filepath.Join("shared", handler.Filename))
	if err != nil {
		http.Error(w, "Failed to save", 500)
		return
	}
	defer dst.Close()

	io.Copy(dst, file)
	w.Write([]byte("Uploaded!"))
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("shared")
	if err != nil {
		http.Error(w, "Failed to read directory", 500)
		return
	}

	var filenames []string
	for _, file := range files {
		if !file.IsDir() {
			filenames = append(filenames, file.Name())
		}
	}

	json.NewEncoder(w).Encode(filenames)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing file name", 400)
		return
	}

	filePath := filepath.Join("shared", name)
	http.ServeFile(w, r, filePath)
}
