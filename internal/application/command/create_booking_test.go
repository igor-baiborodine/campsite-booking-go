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

func TestCreateBookingHandler(t *testing.T) {
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
	errBookingDatesNotAvailable := domain.ErrBookingDatesNotAvailable{
		StartDate: booking.StartDate,
		EndDate:   booking.EndDate,
	}
	cmd := CreateBooking{
		BookingID:  booking.BookingID,
		CampsiteID: booking.CampsiteID,
		Email:      booking.Email,
		FullName:   booking.FullName,
		StartDate:  booking.StartDate.Format(time.DateOnly),
		EndDate:    booking.EndDate.Format(time.DateOnly),
	}

	tests := map[string]struct {
		cmd     CreateBooking
		on      func(f mocks)
		wantErr error
	}{
		"Success": {
			cmd: cmd,
			on: func(f mocks) {
				f.bookings.On(
					"Insert", context.TODO(), booking,
				).Return(nil)
			},
			wantErr: nil,
		},
		"Error_ErrBookingDatesNotAvailable": {
			cmd: cmd,
			on: func(f mocks) {
				f.bookings.On(
					"Insert", context.TODO(), booking,
				).Return(errBookingDatesNotAvailable)
			},
			wantErr: errBookingDatesNotAvailable,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			m := mocks{
				bookings: domain.NewMockBookingRepository(t),
			}
			h := NewCreateBookingHandler(m.bookings)
			if tc.on != nil {
				tc.on(m)
			}
			// when
			err := h.Handle(context.TODO(), tc.cmd)
			// then
			assert.ErrorIs(t, err, tc.wantErr,
				"CreateBookingHandler.Handle() error = %v, wantErr %v", err, tc.wantErr)
			mock.AssertExpectationsForObjects(t, m.bookings)
		})
	}
}
