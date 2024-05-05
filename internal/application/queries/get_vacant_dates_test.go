package queries_test

import (
	"context"
	"errors"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/queries"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func parseDateStr(t *testing.T, d string) time.Time {
	date, err := time.Parse(time.DateOnly, d)
	if err != nil {
		t.Errorf("cannot parse date %s", d)
	}
	return date
}

func TestGetVacantDates(t *testing.T) {
	campsiteID := "campsite-id"

	type args struct {
		ctx context.Context
		qry queries.GetVacantDates
	}

	type mocks struct {
		bookings *domain.MockBookingRepository
	}

	tests := map[string]struct {
		args    args
		on      func(f mocks)
		want    []string
		wantErr string
	}{
		"ParseStartDateError": {
			args: args{
				ctx: context.Background(),
				qry: queries.GetVacantDates{
					CampsiteID: campsiteID,
					StartDate:  "9999-88-88",
					EndDate:    "2006-01-02",
				},
			},
			on:      nil,
			want:    nil,
			wantErr: "cannot parse start date 9999-88-88",
		},
		"ParseEndDateError": {
			args: args{
				ctx: context.Background(),
				qry: queries.GetVacantDates{
					CampsiteID: campsiteID,
					StartDate:  "2006-01-02",
					EndDate:    "9999-99-99",
				},
			},
			on:      nil,
			want:    nil,
			wantErr: "cannot parse end date 9999-99-99",
		},
		"FindForDateRangeError": {
			args: args{
				ctx: context.Background(),
				qry: queries.GetVacantDates{
					CampsiteID: campsiteID,
					StartDate:  "2006-01-02",
					EndDate:    "2006-01-03",
				},
			},
			on: func(f mocks) {
				f.bookings.On(
					"FindForDateRange", context.Background(), campsiteID,
					parseDateStr(t, "2006-01-02"), parseDateStr(t, "2006-01-03"),
				).Return(nil, errors.New("begin transaction"))
			},
			want:    nil,
			wantErr: "begin transaction",
		},
		"NoBookingsFoundForGivenDateRange": {
			args: args{
				ctx: context.Background(),
				qry: queries.GetVacantDates{
					CampsiteID: campsiteID,
					StartDate:  "2006-01-02",
					EndDate:    "2006-01-03",
				},
			},
			on: func(f mocks) {
				f.bookings.On(
					"FindForDateRange", context.Background(), campsiteID,
					parseDateStr(t, "2006-01-02"), parseDateStr(t, "2006-01-03"),
				).Return([]*domain.Booking{}, nil)
			},
			want:    []string{"2006-01-02"},
			wantErr: "",
		},
		"BookingsFoundForGivenDateRange": {
			args: args{
				ctx: context.Background(),
				qry: queries.GetVacantDates{
					CampsiteID: campsiteID,
					StartDate:  "2006-01-02",
					EndDate:    "2006-01-03",
				},
			},
			on: func(f mocks) {
				f.bookings.On(
					"FindForDateRange", context.Background(), campsiteID,
					parseDateStr(t, "2006-01-02"), parseDateStr(t, "2006-01-03"),
				).Return([]*domain.Booking{
					{StartDate: parseDateStr(t, "2006-01-02"), EndDate: parseDateStr(t, "2006-01-03")}}, nil)
			},
			want:    nil,
			wantErr: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			m := mocks{
				bookings: domain.NewMockBookingRepository(t),
			}
			h := queries.NewGetVacantDatesHandler(m.bookings)
			if tc.on != nil {
				tc.on(m)
			}
			// when
			vacantDates, err := h.GetVacantDates(tc.args.ctx, tc.args.qry)
			// then
			if tc.wantErr != "" {
				assert.Containsf(t, err.Error(), tc.wantErr, "GetVacantDates() error=%v, wantErr %v", err, tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, vacantDates)
		})
	}
}
