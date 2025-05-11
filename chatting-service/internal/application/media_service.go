package application

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"go.uber.org/zap"
)

type MediaService struct {
	storage domain.MediaStorage
}

func NewMediaService(storage domain.MediaStorage) *MediaService {
	return &MediaService{
		storage: storage,
	}
}

func (s *MediaService) Upload(ctx context.Context, userID uint, file io.Reader, filename string, contentType string, size int64, userId uint) (*domain.MediaResponse, error) {
	// Validate file size
	if size > 10*1024*1024 { // 10MB
		shared.Log.Debug("file too large",
			zap.Int64("size", size),
			zap.Uint("userID", userID))
		return nil, shared.ErrValidation.WithDetails("file size exceeds 10MB limit")
	}

	// Generate unique filename
	uniqueFilename := generateUniqueFilename(userID, filename)

	// Upload to storage
	path, err := s.storage.Upload(ctx, file, uniqueFilename, contentType, size, userID)
	if err != nil {
		shared.Log.Error("failed to upload media",
			zap.String("filename", filename),
			zap.Uint("userID", userID),
			zap.Error(err))
		return nil, err
	}

	// Get public URL
	url, err := s.storage.GetURL(ctx, path)
	if err != nil {
		shared.Log.Error("failed to get media URL",
			zap.String("path", path),
			zap.Uint("userID", userID),
			zap.Error(err))
		return nil, err
	}

	return &domain.MediaResponse{
		URL:         url,
		Size:        size,
		ContentType: contentType,
		UploadedAt:  time.Now().UTC(),
	}, nil
}

func (s *MediaService) Delete(ctx context.Context, userID uint, path string) error {
	err := s.storage.Delete(ctx, path)
	if err != nil {
		shared.Log.Error("failed to delete media",
			zap.String("path", path),
			zap.Uint("userID", userID),
			zap.Error(err))
		return err
	}
	return nil
}

func (s *MediaService) GetByUser(ctx context.Context, userID uint) ([]domain.MediaResponse, error) {
	return nil, nil
}

func generateUniqueFilename(userID uint, original string) string {
	ext := filepath.Ext(original)
	// base := strings.TrimSuffix(filepath.Base(original), ext)

	randStr, _ := generateRandomString(8) // 8 character random string

	return fmt.Sprintf("%d_%d_%s%s",
		userID,
		time.Now().UnixNano(),
		randStr,
		ext)
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
