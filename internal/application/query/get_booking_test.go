package query

import (
	"context"
	"testing"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetBookingHandler(t *testing.T) {
	type mocks struct {
		bookings *domain.MockBookingRepository
	}
	booking, err := bootstrap.NewBooking("campsite-id")
	if err != nil {
		t.Fatalf("create booking error: %v", err)
	}

	tests := map[string]struct {
		qry     GetBooking
		on      func(f mocks)
		want    *domain.Booking
		wantErr error
	}{
		"Success": {
			qry: GetBooking{BookingID: booking.BookingID},
			on: func(f mocks) {
				f.bookings.On(
					"Find", context.TODO(), booking.BookingID,
				).Return(booking, nil)
			},
			want:    booking,
			wantErr: nil,
		},
		"Error_BeginTx": {
			qry: GetBooking{BookingID: booking.BookingID},
			on: func(f mocks) {
				f.bookings.On(
					"Find", context.TODO(), booking.BookingID,
				).Return(nil, bootstrap.ErrBeginTx)
			},
			want:    nil,
			wantErr: bootstrap.ErrBeginTx,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			m := mocks{
				bookings: domain.NewMockBookingRepository(t),
			}
			h := NewGetBookingHandler(m.bookings)
			if tc.on != nil {
				tc.on(m)
			}
			// when
			got, err := h.Handle(context.TODO(), tc.qry)
			// then
			assert.Equal(t, tc.want, got,
				"GetBookingHandler.Handle() got = %v, want %v", got, tc.want)
			assert.ErrorIs(t, err, tc.wantErr,
				"GetBookingHandler.Handle() error = %v, wantErr %v", err, tc.wantErr)
			mock.AssertExpectationsForObjects(t, m.bookings)
		})
	}
}
