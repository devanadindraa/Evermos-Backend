package fileutils

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func IsValidImage(file *multipart.FileHeader) bool {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validImageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}

	for _, validExt := range validImageExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

func IsValidVideo(file *multipart.FileHeader) bool {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validVideoExtensions := []string{".mp4", ".avi", ".mov", ".mkv", ".flv", ".webm"}

	for _, validExt := range validVideoExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

// Example of image validation function
func IsValidImageExtension(ext string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

// Helper function to check if a file extension is a video
func IsVideo(extension string) bool {
	// List of video file extensions
	videoExtensions := []string{".mp4", ".avi", ".mov", ".mkv", ".flv", ".webm"}
	for _, ext := range videoExtensions {
		if strings.ToLower(extension) == ext {
			return true
		}
	}
	return false
}

// Helper function to generate a name for media files
func GenerateMediaName(productID string) (string, error) {
	// Generate a unique name for the media file
	return fmt.Sprintf("%s_%d", productID, time.Now().UnixNano()), nil
}

func SaveMedia(ctx context.Context, file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file on disk: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy file to disk: %w", err)
	}

	return nil
}
