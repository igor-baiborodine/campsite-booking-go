package command

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	"github.com/stackus/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateBookingHandler(t *testing.T) {
	type mocks struct {
		bookings  *domain.MockBookingRepository
		validator *domain.MockBookingValidator
	}
	campsiteID := uuid.New().String()
	booking, err := bootstrap.NewBooking(campsiteID)
	if err != nil {
		t.Fatalf("create booking error: %v", err)
	}
	booking.ID = 0
	booking.Active = true
	errBookingAlreadyCancelled := domain.ErrBookingAlreadyCancelled{BookingID: booking.BookingID}
	monthOutOfRangeDate := "2024-99-01"

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
				booking.Active = true
				f.bookings.
					On("Find", context.TODO(), booking.BookingID).
					Return(booking, nil).
					On("Update", context.TODO(), booking).
					Return(nil)
				f.validator.
					On("Validate", booking).
					Return(nil)
			},
			wantErr: nil,
		},
		"Error_Find_BeginTx": {
			cmd: cmd,
			on: func(f mocks) {
				booking.Active = true
				f.bookings.
					On("Find", context.TODO(), booking.BookingID).
					Return(nil, bootstrap.ErrBeginTx)
			},
			wantErr: bootstrap.ErrBeginTx,
		},
		"Error_BookingAlreadyCancelled": {
			cmd: cmd,
			on: func(f mocks) {
				booking.Active = false
				f.bookings.
					On("Find", context.TODO(), booking.BookingID).
					Return(booking, nil)
			},
			wantErr: errBookingAlreadyCancelled,
		},
		"Error_ParseStartDate": {
			cmd: UpdateBooking{
				BookingID:  cmd.BookingID,
				CampsiteID: cmd.CampsiteID,
				Email:      cmd.Email,
				FullName:   cmd.FullName,
				StartDate:  monthOutOfRangeDate,
				EndDate:    cmd.EndDate,
			},
			on: func(f mocks) {
				booking.Active = true
				f.bookings.
					On("Find", context.TODO(), booking.BookingID).
					Return(booking, nil)
			},
			wantErr: &time.ParseError{Value: monthOutOfRangeDate},
		},
		"Error_ParseEndDate": {
			cmd: UpdateBooking{
				BookingID:  cmd.BookingID,
				CampsiteID: cmd.CampsiteID,
				Email:      cmd.Email,
				FullName:   cmd.FullName,
				StartDate:  cmd.StartDate,
				EndDate:    monthOutOfRangeDate,
			},
			on: func(f mocks) {
				booking.Active = true
				f.bookings.
					On("Find", context.TODO(), booking.BookingID).
					Return(booking, nil)
			},
			wantErr: &time.ParseError{Value: monthOutOfRangeDate},
		},
		// TODO: add test case for validate error
		"Error_Update_CommitTx": {
			cmd: cmd,
			on: func(f mocks) {
				booking.Active = true
				f.bookings.
					On("Find", context.TODO(), booking.BookingID).
					Return(booking, nil).
					On("Update", context.TODO(), booking).
					Return(bootstrap.ErrCommitTx)
				f.validator.
					On("Validate", booking).
					Return(nil)
			},
			wantErr: bootstrap.ErrCommitTx,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			m := mocks{
				bookings:  domain.NewMockBookingRepository(t),
				validator: domain.NewMockBookingValidator(t),
			}
			var validators []domain.BookingValidator
			validators = append(validators, m.validator)
			h := NewUpdateBookingHandler(m.bookings, validators)

			if tc.on != nil {
				tc.on(m)
			}
			// when
			err := h.Handle(context.TODO(), tc.cmd)
			// then
			defer mock.AssertExpectationsForObjects(t, m.bookings)

			var parseErr *time.ParseError
			if errors.As(err, &parseErr) {
				assert.Equalf(t, monthOutOfRangeDate, parseErr.Value,
					"UpdateBookingHandler.Handle() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			assert.ErrorIs(t, err, tc.wantErr,
				"UpdateBookingHandler.Handle() error = %v, wantErr %v", err, tc.wantErr)
		})
	}
}
