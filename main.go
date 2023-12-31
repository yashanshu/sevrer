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
    "strings"
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

const constantDirPath = "static/uploads"
var homeDir string
var fullDirPath string

func init() {
    var err error
    homeDir, err = os.UserHomeDir()
    if err != nil {
       log.Fatal(err)
    }
}

func main() {
<<<<<<< HEAD
    // inside paul
=======

    fullDirPath = filepath.Join(homeDir, constantDirPath)
>>>>>>> refs/remotes/origin/main

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
    defer file.Close()

    // Validate the file
    if err := validateFile(file, header); err != nil {
        log.Printf("Validation Failed: %s", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err = os.MkdirAll(fullDirPath, 0755)
    if err != nil {
        log.Printf("Error creating the uploads directory: %s", err)
        http.Error(w, "Error creating file", http.StatusInternalServerError)
        return
    }

    if _, err := file.Seek(0, io.SeekStart); err != nil {
        log.Fatal(err)
    }
    uploadedFile, err := os.Create(filepath.Join(fullDirPath, header.Filename))
    if err != nil {
        log.Printf("Error creating the file: \n", err)
        http.Error(w, "Error creating file", http.StatusInternalServerError)
        return
    }
    defer uploadedFile.Close()

    // copy the file content using buffer size
    bufferSize := 8192
    bytesRead, err := io.CopyBuffer(uploadedFile, file, make([]byte, bufferSize))
    if err != nil {
        log.Printf("Error copying the file:\n", err)
        http.Error(w, "Error copying the file\n", http.StatusInternalServerError)
        return
    }
    fmt.Printf("BytesRead during copy: % d, actual bytes: % d\n", bytesRead, header.Size)

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
    ext = strings.ToLower(ext)
    if !allowedExtensions[ext] {
        return errors.New("Invalid file extension")
    }

    // Validate MIME type
    buffer := make([]byte, 512) // read the first 512 bytes to detect the mime type
    _, err := file.Read(buffer)
    if err != nil && err != io.EOF {
        return errors.New("Error reading the file")
    }
    mimeType := http.DetectContentType(buffer)
    if !isValidMimeType(mimeType) {
            return errors.New("Invalid MIME type")
    }

    fmt.Printf("Initial Bytes: %x\n", buffer[:32])

    // Validate file size (limit 5mb)
    const maxFileSize = 5 << 20 // shift 5(`101`) to 20 places left.
    if header.Size > maxFileSize {
        return errors.New("File size exceeds the limit")
    }

    // some more

    return nil
}

func isValidMimeType(mimeType string) bool {
    validMimeTypes := map[string]bool{
    "image/jpeg":       true,
    "image/png":        true,
    "text/markdown":    true,
    "application/epub+zip": true,
    "application/pdf":  true,
    "text/plain":       true,
    }
    return validMimeTypes[mimeType]
}
