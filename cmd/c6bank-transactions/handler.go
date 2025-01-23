package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"git.home/c6bank-transactions/internal/parser"
)

const (
	qifMIME = "text/qif"
)

//go:embed index.html
var indexHTML []byte

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	_, err := w.Write(indexHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	var invoiceRef string

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 32 MB is the default used by FormFile
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get a reference to the fileHeaders

	files, ok := r.MultipartForm.File["file"]
	if !ok {
		http.Error(w, "missing `file` param", http.StatusBadRequest)
		return
	}
	file := files[0]

	// get invoice month if given
	if invoiceReference, ok := r.MultipartForm.Value["invoice_reference"]; ok {
		invoiceRef = invoiceReference[0]
	}

	number, ok := r.MultipartForm.Value["number"]
	if !ok || len(number) > 1 {
		http.Error(w, "field `number` should by unique or empty", http.StatusBadRequest)
		return
	}

	if file.Size > MAX_UPLOAD_SIZE {
		http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 10MB in size", file.Filename), http.StatusBadRequest)
		return
	}

	body, err := file.Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer body.Close()

	filetype, err := validateUploadFile(file.Filename, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := parser.Parse(file.Filename, body, file.Size, number[0], invoiceRef)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse %s: %s", file.Filename, err), http.StatusBadRequest)
		return
	}

	filename := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename)) + ".qif"
	log.Printf("INFO received upload %s of type %s and parsed as %s\n", file.Filename, filetype, filename)

	w.Header().Set("Content-Type", qifMIME)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment;filename="%s"`, filename))

	_, err = io.Copy(w, output)
	if err != nil {
		http.Error(w, fmt.Sprintf("cloud not write response: %s", err), http.StatusBadRequest)
		return
	}
}

func validateUploadFile(name string, file multipart.File) (string, error) {
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return "", err
	}

	filetype := http.DetectContentType(buff)

	if err := parser.IsValid(name, filetype); err != nil {
		return "", err
	}

	_, err := file.Seek(0, io.SeekStart)
	return filetype, err
}
