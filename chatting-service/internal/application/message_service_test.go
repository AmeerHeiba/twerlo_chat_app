package application

import (
	"context"
	"io"
	"testing"

	"github.com/AmeerHeiba/chatting-service/internal/domain"
	"github.com/AmeerHeiba/chatting-service/internal/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) CreateWithRecipients(ctx context.Context, message *domain.Message, recipientIDs []uint) (*domain.Message, error) {
	args := m.Called(ctx, message, recipientIDs)
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageRepository) Create(ctx context.Context, senderID uint, content, mediaURL string, messageType domain.MessageType) (*domain.Message, error) {
	args := m.Called(ctx, senderID, content, mediaURL, messageType)
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageRepository) Update(ctx context.Context, messageID uint, recipientID *uint, broadcasterID *uint) error {
	args := m.Called(ctx, messageID, recipientID, broadcasterID)
	return args.Error(0)
}

func (m *MockMessageRepository) FindByID(ctx context.Context, messageID uint) (*domain.Message, error) {
	args := m.Called(ctx, messageID)
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageRepository) Delete(ctx context.Context, messageID uint, recipientID uint) error {
	args := m.Called(ctx, messageID, recipientID)
	return args.Error(0)
}

func (m *MockMessageRepository) FindBroadcasts(ctx context.Context, broadcasterID uint, query domain.MessageQuery) ([]domain.Message, error) {
	args := m.Called(ctx, broadcasterID, query)
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (m *MockMessageRepository) FindConversation(ctx context.Context, userID1, userID2 uint, query domain.MessageQuery) ([]domain.Message, error) {
	args := m.Called(ctx, userID1, userID2, query)
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (m *MockMessageRepository) FindUserMessages(ctx context.Context, userID uint, query domain.MessageQuery) ([]domain.Message, error) {
	args := m.Called(ctx, userID, query)
	return args.Get(0).([]domain.Message), args.Error(1)
}

func (m *MockMessageRepository) MarkAsDelivered(ctx context.Context, messageID uint) error {
	args := m.Called(ctx, messageID)
	return args.Error(0)
}

func (m *MockMessageRepository) MarkAsRead(ctx context.Context, messageID uint, recipientID uint) error {
	args := m.Called(ctx, messageID, recipientID)
	return args.Error(0)
}

type MockMessageNotifier struct {
	mock.Mock
}

func (m *MockMessageNotifier) Notify(ctx context.Context, message *domain.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageNotifier) Broadcast(ctx context.Context, message *domain.Message, recipientIDs []uint) error {
	args := m.Called(ctx, message, recipientIDs)
	return args.Error(0)
}

type MockMediaUploader struct {
	mock.Mock
}

func (m *MockMediaUploader) Upload(ctx context.Context, userID uint, reader io.Reader, fileName, contentType string, fileSize int64) (*domain.MediaResponse, error) {
	args := m.Called(ctx, userID, reader, fileName, contentType, fileSize)
	return args.Get(0).(*domain.MediaResponse), args.Error(1)
}

func TestMessageService_SendDirectMessage(t *testing.T) {
	shared.InitLogger("test")

	tests := []struct {
		name          string
		senderID      uint
		recipientID   uint
		content       string
		mediaURL      string
		mockSetup     func(*MockMessageRepository, *MockUserRepository, *MockMessageNotifier)
		expectedError error
	}{
		{
			name:        "Success",
			senderID:    1,
			recipientID: 2,
			content:     "Hello",
			mockSetup: func(mr *MockMessageRepository, ur *MockUserRepository, mn *MockMessageNotifier) {
				ur.On("Exists", mock.Anything, uint(2)).Return(true, nil)
				mr.On("Create", mock.Anything, uint(1), "Hello", "", domain.MessageDirect).
					Return(&domain.Message{Content: "Hello"}, nil)
				mr.On("Update", mock.Anything, uint(1), mock.Anything, nil).Return(nil)
				mr.On("FindByID", mock.Anything, uint(1)).Return(&domain.Message{Content: "Hello"}, nil)
				mn.On("Notify", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "Recipient does not exist",
			senderID:    1,
			recipientID: 2,
			content:     "Hello",
			mockSetup: func(mr *MockMessageRepository, ur *MockUserRepository, mn *MockMessageNotifier) {
				ur.On("Exists", mock.Anything, uint(2)).Return(false, nil)
			},
			expectedError: shared.ErrNotFound,
		},
		{
			name:        "Empty content and media URL",
			senderID:    1,
			recipientID: 2,
			content:     "",
			mediaURL:    "",
			mockSetup: func(mr *MockMessageRepository, ur *MockUserRepository, mn *MockMessageNotifier) {
				ur.On("Exists", mock.Anything, uint(2)).Return(true, nil)
			},
			expectedError: shared.ErrValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr := &MockMessageRepository{}
			ur := &MockUserRepository{}
			mn := &MockMessageNotifier{}
			mu := &MockMediaUploader{}

			if tt.mockSetup != nil {
				tt.mockSetup(mr, ur, mn)
			}

			service := NewMessageService(mr, nil, ur, mn, mu)
			_, err := service.SendDirectMessage(context.Background(), tt.senderID, tt.recipientID, tt.content, tt.mediaURL)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mr.AssertExpectations(t)
			ur.AssertExpectations(t)
			mn.AssertExpectations(t)
		})
	}
}
