package query

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/jba/slog/handlers/loghandler"
	"github.com/stretchr/testify/assert"
)

func parseDateStr(t *testing.T, d string) time.Time {
	date, err := time.Parse(time.DateOnly, d)
	if err != nil {
		t.Errorf("failed to parse date %s", d)
	}
	return date
}

func TestGetVacantDates(t *testing.T) {
	campsiteID := "campsite-id"

	type args struct {
		ctx context.Context
		qry GetVacantDates
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
		"Success": {
			args: args{
				ctx: context.Background(),
				qry: GetVacantDates{
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
				qry: GetVacantDates{
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
		"ParseStartDateError": {
			args: args{
				ctx: context.Background(),
				qry: GetVacantDates{
					CampsiteID: campsiteID,
					StartDate:  "9999-88-88",
					EndDate:    "2006-01-02",
				},
			},
			on:      nil,
			want:    nil,
			wantErr: "failed to parse start date 9999-88-88",
		},
		"ParseEndDateError": {
			args: args{
				ctx: context.Background(),
				qry: GetVacantDates{
					CampsiteID: campsiteID,
					StartDate:  "2006-01-02",
					EndDate:    "9999-99-99",
				},
			},
			on:      nil,
			want:    nil,
			wantErr: "failed to parse end date 9999-99-99",
		},
		"FindForDateRangeError": {
			args: args{
				ctx: context.Background(),
				qry: GetVacantDates{
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			m := mocks{
				bookings: domain.NewMockBookingRepository(t),
			}
			h := NewGetVacantDatesHandler(
				m.bookings,
				slog.New(loghandler.New(os.Stdout, nil)),
			)
			if tc.on != nil {
				tc.on(m)
			}
			// when
			vacantDates, err := h.Handle(tc.args.ctx, tc.args.qry)
			// then
			if tc.wantErr != "" {
				assert.Containsf(t, err.Error(), tc.wantErr, "GetVacantDates() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, vacantDates)
		})
	}
}
