package query

import (
	"context"
	"testing"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	"github.com/stackus/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func parseDateStr(t *testing.T, d string) time.Time {
	date, err := time.Parse(time.DateOnly, d)
	if err != nil {
		t.Errorf("failed to parse date %s", d)
	}
	return date
}

func TestGetVacantDatesHandler(t *testing.T) {
	type mocks struct {
		bookings *domain.MockBookingRepository
	}
	campsiteID := "campsite-id"
	monthOutOfRangeDate := "2024-99-01"

	tests := map[string]struct {
		qry     GetVacantDates
		on      func(f mocks)
		want    []string
		wantErr error
	}{
		"Success_NoBookingsFoundForGivenDateRange": {
			qry: GetVacantDates{
				CampsiteID: campsiteID,
				StartDate:  "2006-01-02",
				EndDate:    "2006-01-03",
			},
			on: func(f mocks) {
				f.bookings.On(
					"FindForDateRange", context.TODO(), campsiteID,
					parseDateStr(t, "2006-01-02"), parseDateStr(t, "2006-01-03"),
				).Return([]*domain.Booking{}, nil)
			},
			want:    []string{"2006-01-02"},
			wantErr: nil,
		},
		"Success_BookingsFoundForGivenDateRange": {
			qry: GetVacantDates{
				CampsiteID: campsiteID,
				StartDate:  "2006-01-02",
				EndDate:    "2006-01-03",
			},
			on: func(f mocks) {
				f.bookings.On(
					"FindForDateRange", context.TODO(), campsiteID,
					parseDateStr(t, "2006-01-02"), parseDateStr(t, "2006-01-03"),
				).Return([]*domain.Booking{
					{
						StartDate: parseDateStr(t, "2006-01-02"),
						EndDate:   parseDateStr(t, "2006-01-03"),
					},
				}, nil)
			},
			want:    nil,
			wantErr: nil,
		},
		"Error_ParseStartDate": {
			qry: GetVacantDates{
				CampsiteID: campsiteID,
				StartDate:  monthOutOfRangeDate,
				EndDate:    "2006-01-02",
			},
			on:      nil,
			want:    nil,
			wantErr: &time.ParseError{Value: monthOutOfRangeDate},
		},
		"Error_ParseEndDate": {
			qry: GetVacantDates{
				CampsiteID: campsiteID,
				StartDate:  "2006-01-02",
				EndDate:    monthOutOfRangeDate,
			},
			on:      nil,
			want:    nil,
			wantErr: &time.ParseError{Value: monthOutOfRangeDate},
		},
		"Error_BeginTx": {
			qry: GetVacantDates{
				CampsiteID: campsiteID,
				StartDate:  "2006-01-02",
				EndDate:    "2006-01-03",
			},
			on: func(f mocks) {
				f.bookings.On(
					"FindForDateRange", context.TODO(), campsiteID,
					parseDateStr(t, "2006-01-02"), parseDateStr(t, "2006-01-03"),
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
			h := NewGetVacantDatesHandler(m.bookings)
			if tc.on != nil {
				tc.on(m)
			}
			// when
			got, err := h.Handle(context.TODO(), tc.qry)
			// then
			assert.Equal(t, tc.want, got)
			if err != nil {
				var parseErr *time.ParseError
				if errors.As(err, &parseErr) {
					assert.Equalf(t, monthOutOfRangeDate, parseErr.Value,
						"GetVacantDatesHandler.Handle() error = %v, wantErr %v", err, tc.wantErr)
				}
			}
			mock.AssertExpectationsForObjects(t, m.bookings)
		})
	}
}
