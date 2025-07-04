package controllers

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
)

// ProcessImage converts an image file header into a raw base64 string.
// This raw base64 string is what Ollama expects for multimodal input.
func ProcessImage(image *multipart.FileHeader) (string, error) {
	// Open the file
	file, err := image.Open()
	fmt.Println("Image file opened:-------------------------------------------------------", file)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the file into a byte slice
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	fmt.Println("Image bytes length:-------------------------------------------------------", imageBytes)

	// Convert the byte slice to a base64 encoded string
	base64Str := base64.StdEncoding.EncodeToString(imageBytes)

	// Return just the raw base64 string
	return base64Str, nil
}