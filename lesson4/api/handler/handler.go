package handler

import (
	"encoding/json"
	"fmt"
	"gb-backend/lesson4/app/repos/files"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Router struct {
	*http.ServeMux
	files *files.Files
}

func NewRouter(files *files.Files) *Router {
	r := &Router{
		ServeMux: http.NewServeMux(),
		files:    files,
	}

	r.Handle("/", http.FileServer(http.Dir(r.files.Dir)))
	r.Handle("/upload", http.HandlerFunc(r.Upload))
	r.Handle("/files", http.HandlerFunc(r.ListFiles))

	return r
}

func (rt *Router) Upload(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to read file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read file", http.StatusBadRequest)
			return
		}

		filePath := filepath.Join(rt.files.Dir, header.Filename)

		err = os.WriteFile(filePath, data, 0777)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}

		fileLink := r.Host + "/" + header.Filename
		fmt.Fprintln(w, fileLink)

	default:
	}
}

func (rt *Router) ListFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	ext := r.FormValue("ext")
	files, err := rt.files.Search(ext)
	if err != nil {
		http.Error(w, "Unable to get list of files", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(files)
}
