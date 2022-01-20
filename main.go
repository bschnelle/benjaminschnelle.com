package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func writePostHTML(sourcePath string, destPath string) {
	handle, err := os.Open(sourcePath)
	handleError(err)

	contents, err := io.ReadAll(handle)
	handleError(err)

	var sb strings.Builder
	sb.WriteString("<!DOCTYPE html>")
	sb.WriteString("<html>")
	sb.WriteString("<body>")
	sb.WriteString(strings.Replace(string(contents), "\n", "<br/>", -1))
	sb.WriteString("</body>")
	sb.WriteString("</html>")

	err = os.WriteFile(destPath, []byte(sb.String()), 0644)
	handleError(err)
}

func writeIndexHTML(path string, fileName string) {
	handle, err := os.Open(filepath.Join(path, fileName))
	handleError(err)

	contents, err := io.ReadAll(handle)
	handleError(err)

	err = os.WriteFile(filepath.Join(path, "index.html"), contents, 0644)
	handleError(err)
}

func createPublicDirectory() {
	err := os.RemoveAll("public")
	handleError(err)

	publicDir := "public"
	err = os.Mkdir(publicDir, 0777)
	handleError(err)

	postsDir := "posts"
	files, err := os.ReadDir(postsDir)
	handleError(err)

	wd, err := os.Getwd()
	handleError(err)

	for _, file := range files {
		fileName := file.Name()
		sourcePath := filepath.Join(wd, postsDir, fileName)
		publicPath := filepath.Join(wd, publicDir, strings.Replace(fileName, ".md", ".html", 1))
		writePostHTML(sourcePath, publicPath)
	}

	writeIndexHTML(filepath.Join(wd, publicDir), "blog-from-scratch.html")
}

func main() {
	fmt.Println("Starting up...")
	createPublicDirectory()
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Fatal(http.ListenAndServe(":80", nil))
}
