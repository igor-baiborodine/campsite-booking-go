package command

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateBookingHandler(t *testing.T) {
	type mocks struct {
		bookings *domain.MockBookingRepository
	}
	campsiteID := uuid.New().String()
	booking, err := bootstrap.NewBooking(campsiteID)
	if err != nil {
		t.Fatalf("create booking error: %v", err)
	}
	booking.ID = 0
	booking.Active = true

	cmd := UpdateBooking{
		BookingID:  booking.BookingID,
		CampsiteID: booking.CampsiteID,
		Email:      booking.Email,
		FullName:   booking.FullName,
		StartDate:  booking.StartDate.Format(time.DateOnly),
		EndDate:    booking.EndDate.Format(time.DateOnly),
	}

	tests := map[string]struct {
		cmd     UpdateBooking
		on      func(f mocks)
		wantErr error
	}{
		"Success": {
			cmd: cmd,
			on: func(f mocks) {
				f.bookings.
					On(
						"Find", context.TODO(), booking.BookingID,
					).Return(booking, nil).
					On(
						"Update", context.TODO(), booking,
					).Return(nil)
			},
			wantErr: nil,
		},
		"Error_CommitTx": {
			cmd: cmd,
			on: func(f mocks) {
				f.bookings.
					On(
						"Find", context.TODO(), booking.BookingID,
					).Return(booking, nil).
					On(
						"Update", context.TODO(), booking,
					).Return(bootstrap.ErrCommitTx)
			},
			wantErr: bootstrap.ErrCommitTx,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			m := mocks{
				bookings: domain.NewMockBookingRepository(t),
			}
			h := NewUpdateBookingHandler(m.bookings)
			if tc.on != nil {
				tc.on(m)
			}
			// when
			err := h.Handle(context.TODO(), tc.cmd)
			// then
			assert.ErrorIs(t, err, tc.wantErr,
				"UpdateBookingHandler.Handle() error = %v, wantErr %v", err, tc.wantErr)
			mock.AssertExpectationsForObjects(t, m.bookings)
		})
	}
}
