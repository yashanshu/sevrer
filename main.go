package main

import (
    "errors"
    "fmt"
    "log"
    "mime/multipart"
    "net/http"
    "path/filepath"
    "io"
    "os"
)

var allowedExtensions = map[string]bool{
    ".md": true,
    ".epub": true,
    ".pdf": true,
    ".txt": true,
    ".jpg": true,
    ".jpeg": true,
    ".png": true,
}

const directoryPath = "~/static/uploads"

func main() {

    http.HandleFunc("/", handleRoot)
    http.HandleFunc("/upload", handleUpload)
    http.HandleFunc("/success", handleSuccess)
    port := ":6969"
    log.Printf("Server is running on port %s\n", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        //fmt.Println("Server Failed to Start:", err)
        log.Fatalf("Server Failed to Start:\n", err)
    }
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    log.Printf("Accessed root path: %s %s", r.Method, r.URL.Path)
    http.ServeFile(w, r, "upload.html")
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
    log.Printf("Recieved file upload page: %s %s\n", r.Method, r.URL.Path)

    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed\n", http.StatusMethodNotAllowed)
        return
    }
    file, header, err := r.FormFile("file")
    log.Printf("Recieved file upload request: %s %s\n", r.Method, r.URL.Path)
    if err != nil {
        log.Printf("Error retrieving the file: \n", err)
        http.Error(w, "Error retrieving the file\n", http.StatusBadRequest)
        return
    }

    // Validate the file
    if err := validateFile(file, header); err != nil {
        log.Printf("Validation Failed: %s", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    err = os.MkdirAll(directoryPath, 755)
    if err != nil {
        log.Printf("Error creating the uploads direcotry: %s", err)
        http.Error(w, "Error creating file", http.StatusInternalServerError)
        return
    }

    uploadedFile, err := os.Create(directoryPath + "/" + header.Filename)
    if err != nil {
        log.Printf("Error creating the file: \n", err)
        http.Error(w, "Error creating file", http.StatusInternalServerError)
        return
    }
    defer uploadedFile.Close()

    _, err = io.Copy(uploadedFile, file)
    if err != nil {
        log.Printf("Error copying the file:\n", err)
        http.Error(w, "Error copying the file\n", http.StatusInternalServerError)
        return
    }

    log.Printf("File %s uploaded successfully!\n", header.Filename)
    //fmt.Println(w, "File %s uploaded successfully!\n", header.Filename)
    http.Redirect(w, r, "/success", http.StatusTemporaryRedirect)
}


func handleSuccess(w http.ResponseWriter, r *http.Request) {
    log.Printf("Accessed success path: %s %s", r.Method, r.URL.Path)
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintf(w, "<h1>File uploaded successfully!</h1>")
}

func validateFile(file multipart.File, header *multipart.FileHeader) error {
    // validate file extensions
    log.Printf("Validating file: %s", header.Filename)
    ext := filepath.Ext(header.Filename)
    if !allowedExtensions[ext] {
        return errors.New("Invalid file extension")
    }

    // Validate MIME type
    buffer := make([]byte, 512) // read the first 512 bytes to detect the mime type
    _, err := file.Read(buffer)
    if err != nil {
        return errors.New("Error reading the file")
    }
    mimeType := http.DetectContentType(buffer)
    fmt.Printf(mimeType)
    if mimeType != "image/jpeg" && mimeType != "image/png" && mimeType != "image/png" && 
        mimeType != "text/Markdown" && mimeType != "application/epub+zip" && 
        mimeType != "application/pdf" && mimeType != "text/plain" {
            return errors.New("Invalid MIME type")
        }

    // Validate file size (limit 5mb)
    const maxFileSize = 5 << 20 // shift 5(`101`) to 20 places left.
    if header.Size > maxFileSize {
        return errors.New("File size exceeds the limit")
    }

    // some more

    return nil
}
