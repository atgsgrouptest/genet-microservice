package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"

	"rsc.io/pdf"
)

func ExtractTextFromPDF(file multipart.File) (string, error) {
	// First get the file size
	size := fileSize(file)

	// Create new reader
	reader, err := pdf.NewReader(file, size)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}

	var buf bytes.Buffer
	for i := 1; i <= reader.NumPage(); i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}

		content := page.Content()
		for _, txt := range content.Text {
			buf.WriteString(txt.S)
			buf.WriteString(" ")
		}
	}
	return buf.String(), nil
}

func fileSize(file multipart.File) int64 {
	if seeker, ok := file.(io.Seeker); ok {
		size, _ := seeker.Seek(0, io.SeekEnd)
		seeker.Seek(0, io.SeekStart)
		return size
	}
	return 0
}
