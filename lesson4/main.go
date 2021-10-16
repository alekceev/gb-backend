package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 1. Возможность получить список файлов на сервере (имя, расширение, размер в байтах)
// 2. Фильтрация списка по расширению
// 3. Тесты с использованием библиотеки httptest

type UploadHandler struct {
	HostAddr  string
	UploadDir string
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

		filePath := filepath.Join(h.UploadDir, header.Filename)

		err = os.WriteFile(filePath, data, 0777)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}

		fileLink := h.HostAddr + "/" + header.Filename
		fmt.Fprintln(w, fileLink)

	default:
	}
}

type File struct {
	Name      string `json:"name"`
	Extension string `json:"ext"`
	Size      int64  `json:"size"`
}

type SearcherInFolder struct {
	Dir string
}

func (fs SearcherInFolder) Search(extension string) ([]File, error) {
	dirFiles, err := os.ReadDir(fs.Dir)
	if err != nil {
		return nil, err
	}
	files := make([]File, 0)
	for _, dirFile := range dirFiles {
		info, err := dirFile.Info()
		if err != nil {
			continue
		}

		arr := strings.Split(info.Name(), ".")
		file_ext := arr[len(arr)-1]

		if extension != file_ext {
			continue
		}

		files = append(files, File{
			Name:      info.Name(),
			Extension: file_ext,
			Size:      info.Size(),
		})
	}
	return files, nil
}

type GetFileListHandler struct {
	FileSearcher SearcherInFolder
}

func (h *GetFileListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ext := r.FormValue("ext")
	files, err := h.FileSearcher.Search(ext)
	if err != nil {
		http.Error(w, "Unable to get list of files", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(files)
	if err != nil {
		http.Error(w, "Error when encoding", http.StatusInternalServerError)
	}
}

func main() {
	uploadHandler := &UploadHandler{
		UploadDir: "upload",
		HostAddr:  "localhost:8080",
	}
	http.Handle("/upload", uploadHandler)

	http.Handle("/files", &GetFileListHandler{
		FileSearcher: SearcherInFolder{Dir: uploadHandler.UploadDir},
	})

	http.Handle("/", http.FileServer(http.Dir(uploadHandler.UploadDir)))

	fs := &http.Server{
		Addr: uploadHandler.HostAddr,
		// Handler:      http.FileServer(http.Dir(uploadHandler.UploadDir)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fs.ListenAndServe()
}
