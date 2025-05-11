package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"go.uber.org/zap"
)

type LocalStorage struct {
	basePath string
	baseURL  string
}

func NewLocalStorage(basePath, baseURL string) *LocalStorage {
	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}
}

func (s *LocalStorage) Upload(ctx context.Context, file io.Reader, filename string, contentType string, size int64, userId uint) (string, error) {
	// Get user ID from context
	// userID, ok := ctx.Value("userID").(uint)
	// if !ok {
	// 	return "", errors.New("user ID not found in context")
	// }

	// Create user-specific directory
	userPath := filepath.Join(s.basePath, fmt.Sprintf("user_%d", userId))
	if err := os.MkdirAll(userPath, 0755); err != nil {
		shared.Log.Error("failed to create user directory",
			zap.String("path", userPath),
			zap.Error(err))
		return "", err
	}

	// Create file
	filePath := filepath.Join(userPath, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		shared.Log.Error("failed to create file",
			zap.String("path", filePath),
			zap.Error(err))
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		shared.Log.Error("failed to write file",
			zap.String("path", filePath),
			zap.Error(err))
		return "", err
	}

	// Return relative path in format "user_<ID>/filename"
	return filepath.Join(fmt.Sprintf("user_%d", userId), filename), nil
}

func (s *LocalStorage) GetURL(ctx context.Context, path string) (string, error) {
	return s.baseURL + "/" + path, nil
}

func (s *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)
	return os.Remove(fullPath)
}

// Dummy for now
func (s *LocalStorage) GetSignedURL(ctx context.Context, path string, expires time.Duration) (string, error) {
	return s.baseURL + "/" + path, nil
}
