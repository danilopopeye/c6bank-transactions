package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"git.home/c6bank-transactions/internal/parser"
)

const (
	qifMIME       = "text/qif"
	csvMIME       = "text/csv"
	maxUploadSize = 32 << 20
)

//go:embed index.html
var indexHTML []byte

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	if _, err := w.Write(indexHTML); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := validate(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	number := r.PostFormValue("number")

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	filename := fileHeader.Filename

	filetype, err := validateUploadFile(fileHeader.Filename, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, outputname, err := parser.Parse(filename, file, fileHeader.Size, number)
	if err != nil {
		fmt.Printf("ERROR file=%q: %s\n", filename, err)
		http.Error(w, fmt.Sprintf("could not parse %s: %s", filename, err), http.StatusBadRequest)

		return
	}

	log.Printf("%s INFO received upload %s of type %s and parsed as %s\n", time.Now().Format(time.RFC3339), filename, filetype, outputname)

	contentType := qifMIME
	if strings.HasSuffix(outputname, ".csv") {
		contentType = csvMIME
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment;filename="%s"`, outputname))

	_, err = io.Copy(w, output)
	if err != nil {
		http.Error(w, fmt.Sprintf("cloud not write response: %s", err), http.StatusBadRequest)

		return
	}
}

func validate(r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method %q not allowed", r.Method)
	}

	files, ok := r.MultipartForm.File["file"]
	if !ok || len(files) != 1 {
		return fmt.Errorf("%w: expected one `file` param", http.ErrMissingFile)
	}

	if files[0].Size > MAX_UPLOAD_SIZE {
		return fmt.Errorf("the uploaded image %q is too big. Please use an image less than 10MB in size", files[0].Filename)
	}

	return nil
}

func validateUploadFile(name string, file io.ReadSeeker) (string, error) {
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
