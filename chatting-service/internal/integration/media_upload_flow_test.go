package integration

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AmeerHeiba/chatting-service/internal/application"
	"github.com/AmeerHeiba/chatting-service/internal/infrastructure/database"
	"github.com/AmeerHeiba/chatting-service/internal/infrastructure/storage"
	"github.com/stretchr/testify/assert"
)

func TestMediaUploadFlow(t *testing.T) {
	db := setupTestDB(t)
	userRepo := database.NewUserRepository(db)

	// Setup storage
	tempDir := t.TempDir()
	storage := storage.NewLocalStorage(tempDir, "http://localhost:8080/media")
	mediaService := application.NewMediaService(storage)

	// Create test user
	user, err := userRepo.Create(context.Background(), "mediauser", "media@test.com", "password")
	assert.NoError(t, err)

	// Create test file
	fileContent := "test file content"
	fileReader := strings.NewReader(fileContent)

	// Upload media
	ctx := context.WithValue(context.Background(), "userID", user.ID)
	mediaResp, err := mediaService.Upload(
		ctx,
		user.ID,
		fileReader,
		"test.txt",
		"text/plain",
		int64(len(fileContent)),
		uint(1),
	)
	assert.NoError(t, err)
	assert.Contains(t, mediaResp.URL, "http://localhost:8080/media/user_")
	assert.Equal(t, int64(len(fileContent)), mediaResp.Size)

	// Verify file exists
	filePath := filepath.Join(tempDir, strings.TrimPrefix(mediaResp.URL, "http://localhost:8080/media/"))
	fileData, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, fileContent, string(fileData))
}
