package checkimage

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

func CheckImage(file multipart.File) error {
	// Create a buffer to store the header of the file
	// The http.DetectContentType function only needs the first 512 bytes
	buff := make([]byte, 512)

	_, err := file.Read(buff)
	if err != nil && err != io.EOF {
		return err
	}

	// Detect the Content-Type
	filetype := http.DetectContentType(buff)

	// Validate the type
	switch filetype {
	case "image/jpeg", "image/png":
		// It is a valid image type

	default:
		// fmt.Errorf("invalid file type. Only JPEG, PNG, and GIF allowed.")
		return errors.New("invalid file type. Only JPEG, PNG, and GIF allowed.")
	}

	// --- CRITICAL STEP ---
	// Because the first 512 bytes are read, the file cursor is now at offset 512.
	// "rewind" the cursor back to the start (0), or the saved file will be corrupted.
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	return nil
}
