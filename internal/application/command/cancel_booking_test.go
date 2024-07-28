package command

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCancelBookingHandler(t *testing.T) {
	type mocks struct {
		bookings *domain.MockBookingRepository
	}
	campsiteID := uuid.New().String()
	booking, err := bootstrap.NewBooking(campsiteID)
	if err != nil {
		t.Fatalf("create active booking error: %v", err)
	}
	booking.ID = 0
	booking.Active = true
	errBookingAlreadyCancelled := domain.ErrBookingAlreadyCancelled{BookingID: booking.BookingID}

	tests := map[string]struct {
		cmd     CancelBooking
		on      func(f mocks)
		wantErr error
	}{
		"Success": {
			cmd: CancelBooking{BookingID: booking.BookingID},
			on: func(f mocks) {
				booking.Active = true
				f.bookings.
					On("Find", context.TODO(), booking.BookingID).
					Return(booking, nil).
					On("Update", context.TODO(), booking).
					Return(nil)
			},
			wantErr: nil,
		},
		"Error_Find_BeginTx": {
			cmd: CancelBooking{BookingID: booking.BookingID},
			on: func(f mocks) {
				booking.Active = true
				f.bookings.
					On("Find", context.TODO(), booking.BookingID).
					Return(nil, bootstrap.ErrBeginTx)
			},
			wantErr: bootstrap.ErrBeginTx,
		},
		"Error_BookingAlreadyCancelled": {
			cmd: CancelBooking{BookingID: booking.BookingID},
			on: func(f mocks) {
				booking.Active = false
				f.bookings.
					On("Find", context.TODO(), booking.BookingID).
					Return(booking, nil)
			},
			wantErr: errBookingAlreadyCancelled,
		},
		"Error_Update_CommitTx": {
			cmd: CancelBooking{BookingID: booking.BookingID},
			on: func(f mocks) {
				booking.Active = true
				f.bookings.
					On("Find", context.TODO(), booking.BookingID).
					Return(booking, nil).
					On("Update", context.TODO(), booking).
					Return(bootstrap.ErrCommitTx)
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
			h := NewCancelBookingHandler(m.bookings)
			if tc.on != nil {
				tc.on(m)
			}
			// when
			err := h.Handle(context.TODO(), tc.cmd)
			// then
			assert.Equalf(t, tc.wantErr, err,
				"CancelBookingHandler.Handle() error = %v, wantErr %v", err, tc.wantErr)
			mock.AssertExpectationsForObjects(t, m.bookings)
		})
	}
}
