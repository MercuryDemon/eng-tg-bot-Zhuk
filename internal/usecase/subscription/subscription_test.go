package subscription

import (
	"context"
	"errors"
	"testing"

	"github.com/krezefal/eng-tg-bot/internal/domain"
	"./mocks"
	"github.com/rs/zerolog"
)

func TestSubscribe(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		userID        int64
		username      string
		dictionary    string
		expectedError error
	}{
		{
			name:          "test1",
			ctx:           context.Background(),
			userID:        1,
			username:      "",
			dictionary:    "",
			expectedError: domain.ErrDictionaryNotFound,
		},

		{
			name:          "test2",
			ctx:           context.Background(),
			userID:        1,
			username:      "",
			dictionary:    "",
			expectedError: domain.ErrAlreadySubscribed,
		},

		{
			name:          "test3",
			ctx:           context.Background(),
			userID:        1,
			username:      "",
			dictionary:    "",
			expectedError: domain.ErrEmptyReviewWordsList,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			useCase := NewUsecase(mocks.MockSubscriptionsRepo.MockUserRepo, mocks.MockDictionaryRepo, mocks.MockSubscriptionsRepo, &zerolog.Logger{})

			err := useCase.Subscribe(tt.ctx, tt.userID, tt.username, tt.dictionary)
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("expected error %v, have %v", tt.expectedError, err)
			}
		})
	}
}

// TODO
// write mocks (mockery maybe)
// write test cases
