package handler

import (
	"bytes"
	"gb-backend/lesson4/app/repos/files"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRouter_Upload(t *testing.T) {
	f := files.NewFiles("../../upload")
	rt := NewRouter(f)

	file, _ := os.Open("../../test.txt")
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	h := http.HandlerFunc(rt.Upload).ServeHTTP

	w := &httptest.ResponseRecorder{}
	r := httptest.NewRequest(http.MethodPost, "/upload", body)
	r.Header.Add("Content-Type", writer.FormDataContentType())

	h(w, r)

	if w.Code != http.StatusOK {
		t.Error("status wrong:", w.Code)
	}
}

func TestRouter_ListFiles(t *testing.T) {
	f := files.NewFiles("upload")
	rt := NewRouter(f)

	h := http.HandlerFunc(rt.ListFiles).ServeHTTP

	w := &httptest.ResponseRecorder{}
	r := httptest.NewRequest(http.MethodGet, "/files?ext=txt", strings.NewReader(""))

	h(w, r)

	if w.Code != http.StatusOK {
		t.Error("status wrong:", w.Code)
	}
}
