package grpc

import (
	"context"
	"testing"
	"time"

	api "github.com/igor-baiborodine/campsite-booking-go/campgroundspb/v1"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/queries"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
)

type mocks struct {
	app *application.MockApp
}

func TestGetBooking(t *testing.T) {

	type args struct {
		ctx context.Context
		req api.GetBookingRequest
	}

	nonExistingID := "non-existing-id"
	booking, err := bootstrap.NewBooking("campsite-id")
	assert.NoError(t, err)

	tests := map[string]struct {
		args    args
		on      func(f mocks)
		want    *api.GetBookingResponse
		wantErr string
	}{
		"Success": {
			args: args{
				ctx: context.Background(),
				req: api.GetBookingRequest{
					BookingId: booking.BookingID,
				},
			},
			on: func(f mocks) {
				f.app.On(
					"GetBooking", context.Background(),
					queries.GetBooking{BookingID: booking.BookingID},
				).Return(booking, nil)
			},
			want:    &api.GetBookingResponse{Booking: bookingFromDomain(booking)},
			wantErr: "",
		},
		"NotFound": {
			args: args{
				ctx: context.Background(),
				req: api.GetBookingRequest{
					BookingId: nonExistingID,
				},
			},
			on: func(f mocks) {
				f.app.On(
					"GetBooking", context.Background(),
					queries.GetBooking{BookingID: nonExistingID},
				).Return(nil, domain.ErrBookingNotFound{BookingID: nonExistingID})
			},
			want:    &api.GetBookingResponse{Booking: bookingFromDomain(booking)},
			wantErr: codes.NotFound.String(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			m := mocks{
				app: application.NewMockApp(t),
			}
			s := server{app: m.app}
			if tc.on != nil {
				tc.on(m)
			}
			// when
			resp, err := s.GetBooking(tc.args.ctx, &tc.args.req)
			// then
			if tc.wantErr != "" {
				assert.Containsf(t, err.Error(), tc.wantErr, "GetBooking() error=%v, wantErr %v", err, tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, resp)
		})
	}
}

func TestCreateBooking(t *testing.T) {

	type args struct {
		ctx context.Context
		req api.CreateBookingRequest
	}

	booking, err := bootstrap.NewBooking("campsite-id")
	assert.NoError(t, err)

	tests := map[string]struct {
		args    args
		on      func(f mocks)
		want    *api.CreateBookingResponse
		wantErr string
	}{
		"Success": {
			args: args{
				ctx: context.Background(),
				req: api.CreateBookingRequest{
					CampsiteId: booking.CampsiteID,
					Email:      booking.Email,
					FullName:   booking.FullName,
					StartDate:  booking.StartDate.Format(time.DateOnly),
					EndDate:    booking.EndDate.Format(time.DateOnly),
				},
			},
			on: func(f mocks) {
				f.app.On(
					"CreateBooking", context.Background(), mock.Anything,
				).Return(booking, nil)
			},
			want:    &api.CreateBookingResponse{BookingId: booking.BookingID},
			wantErr: "",
		},
		"BookingDatesNotAvailable": {
			args: args{
				ctx: context.Background(),
				req: api.CreateBookingRequest{
					CampsiteId: booking.CampsiteID,
					Email:      booking.Email,
					FullName:   booking.FullName,
					StartDate:  booking.StartDate.Format(time.DateOnly),
					EndDate:    booking.EndDate.Format(time.DateOnly),
				},
			},
			on: func(f mocks) {
				f.app.On(
					"CreateBooking", context.Background(), mock.Anything,
				).Return(
					nil,
					domain.ErrBookingDatesNotAvailable{
						StartDate: booking.StartDate,
						EndDate:   booking.EndDate,
					},
				)
			},
			want:    nil,
			wantErr: codes.FailedPrecondition.String(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			m := mocks{
				app: application.NewMockApp(t),
			}
			s := server{app: m.app}
			if tc.on != nil {
				tc.on(m)
			}
			// when
			resp, err := s.CreateBooking(tc.args.ctx, &tc.args.req)
			// then
			if tc.wantErr != "" {
				assert.Containsf(t, err.Error(), tc.wantErr, "CreateBooking() error=%v, wantErr %v", err, tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, resp)
		})
	}
}
