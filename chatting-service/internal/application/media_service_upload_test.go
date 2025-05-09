package application

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMediaStorage struct {
	mock.Mock
}

func (m *MockMediaStorage) Upload(ctx context.Context, file io.Reader, filename string, contentType string, size int64) (string, error) {
	args := m.Called(ctx, file, filename, contentType, size)
	return args.String(0), args.Error(1)
}

func (m *MockMediaStorage) GetURL(ctx context.Context, path string) (string, error) {
	args := m.Called(ctx, path)
	return args.String(0), args.Error(1)
}

func (m *MockMediaStorage) Delete(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(0)
}

func (m *MockMediaStorage) GetSignedURL(ctx context.Context, path string, duration time.Duration) (string, error) {
	args := m.Called(ctx, path, duration)
	return args.String(0), args.Error(1)
}

func TestMediaService_Upload(t *testing.T) {
	tests := []struct {
		name          string
		fileSize      int64
		contentType   string
		mockSetup     func(*MockMediaStorage)
		expectedError error
	}{
		{
			name:        "ValidUpload",
			fileSize:    5 * 1024 * 1024, // 5MB
			contentType: "image/jpeg",
			mockSetup: func(ms *MockMediaStorage) {
				ms.On("Upload", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return("user_1/test.jpg", nil)
				ms.On("GetURL", mock.Anything, "user_1/test.jpg").
					Return("http://example.com/user_1/test.jpg", nil)
			},
			expectedError: nil,
		},
		{
			name:          "FileTooLarge",
			fileSize:      11 * 1024 * 1024, // 11MB
			contentType:   "image/jpeg",
			mockSetup:     func(ms *MockMediaStorage) {},
			expectedError: shared.ErrValidation.WithDetails("file size exceeds 10MB limit"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MockMediaStorage{}
			if tt.mockSetup != nil {
				tt.mockSetup(storage)
			}

			service := NewMediaService(storage)
			file := bytes.NewBufferString("test file content")

			_, err := service.Upload(
				context.WithValue(context.Background(), "userID", uint(1)),
				uint(1),
				file,
				"test.jpg",
				tt.contentType,
				tt.fileSize,
			)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
