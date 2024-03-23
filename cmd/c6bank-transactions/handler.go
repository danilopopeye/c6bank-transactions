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
	"git.home/c6bank-transactions/internal/qif"
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

	number, ok := r.MultipartForm.Value["number"]
	if !ok || len(number) > 1 {
		http.Error(w, "field `number` should by unique or empty", http.StatusBadRequest)
		return
	}

	fileHeader := files[0]
	log.Printf("received upload file %s of size %d\n", fileHeader.Filename, fileHeader.Size)

	if fileHeader.Size > MAX_UPLOAD_SIZE {
		http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 10MB in size", fileHeader.Filename), http.StatusBadRequest)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if err = validateUploadFile(fileHeader.Filename, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	csv, err := qif.ParseCSV(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse %s: %s", fileHeader.Filename, err), http.StatusBadRequest)
		return
	}

	// csv, err := parser.Parse(fileHeader.Filename, file, number[0])
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("could not parse %s: %s", fileHeader.Filename, err), http.StatusBadRequest)
	// 	return
	// }

	out, err := io.ReadAll(csv)
	if err != nil {
		http.Error(w, fmt.Sprintf("cloud not read csv: %s", err), http.StatusBadRequest)
		return
	}

	csvFilename := strings.TrimSuffix(fileHeader.Filename, filepath.Ext(fileHeader.Filename))
	log.Printf("upload parsed as %s.csv: %d\n", csvFilename, len(out))

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment;filename="%s.csv"`, csvFilename))

	_, err = w.Write(out)
	if err != nil {
		http.Error(w, fmt.Sprintf("cloud not write csv response: %s", err), http.StatusBadRequest)
		return
	}
}

func validateUploadFile(name string, file multipart.File) error {
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return err
	}

	filetype := http.DetectContentType(buff)
	log.Printf("received upload %s of type %s\n", name, filetype)

	if err := parser.IsValid(name, filetype); err != nil {
		return err
	}

	_, err := file.Seek(0, io.SeekStart)
	return err
}
