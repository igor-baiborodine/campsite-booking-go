//go:build !integration

package grpc

import (
	"context"
	"testing"
	"time"

	api "github.com/igor-baiborodine/campsite-booking-go/campgroundspb/v1"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/query"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
)

type mocks struct {
	app *application.MockApp
}

func TestServer_CreateCampsite(t *testing.T) {
	type args struct {
		ctx context.Context
		req *api.CreateCampsiteRequest
	}
	campsite, err := bootstrap.NewCampsite()
	assert.NoError(t, err)

	req := &api.CreateCampsiteRequest{
		CampsiteCode:  campsite.CampsiteCode,
		Capacity:      campsite.Capacity,
		DrinkingWater: campsite.DrinkingWater,
		Restrooms:     campsite.Restrooms,
		PicnicTable:   campsite.PicnicTable,
		FirePit:       campsite.FirePit,
	}

	tests := map[string]struct {
		args    args
		on      func(f mocks)
		wantErr error
	}{
		"Success": {
			args: args{ctx: context.Background(), req: req},
			on: func(f mocks) {
				f.app.On(
					"CreateCampsite", context.Background(), mock.Anything,
				).Return(nil)
			},
			wantErr: nil,
		},
		"Error_CommitTx": {
			args: args{ctx: context.Background(), req: req},
			on: func(f mocks) {
				f.app.On(
					"CreateCampsite", context.Background(), mock.Anything,
				).Return(bootstrap.ErrCommitTx)
			},
			wantErr: bootstrap.ErrCommitTx,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			m := mocks{app: application.NewMockApp(t)}
			s := server{app: m.app}
			if tc.on != nil {
				tc.on(m)
			}
			// when
			got, err := s.CreateCampsite(tc.args.ctx, tc.args.req)
			// then
			mock.AssertExpectationsForObjects(t, m.app)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr, "CreateCampsite() error = %v, wantErr %v",
					err, tc.wantErr)
				return
			}
			assert.NotEmpty(t, got.CampsiteId)
		})
	}
}

func TestServer_GetBooking(t *testing.T) {
	type args struct {
		ctx context.Context
		req *api.GetBookingRequest
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
				req: &api.GetBookingRequest{
					BookingId: booking.BookingID,
				},
			},
			on: func(f mocks) {
				f.app.On(
					"GetBooking", context.Background(),
					query.GetBooking{BookingID: booking.BookingID},
				).Return(booking, nil)
			},
			want:    &api.GetBookingResponse{Booking: BookingFromDomain(booking)},
			wantErr: "",
		},
		"Error_NotFound": {
			args: args{
				ctx: context.Background(),
				req: &api.GetBookingRequest{
					BookingId: nonExistingID,
				},
			},
			on: func(f mocks) {
				f.app.On(
					"GetBooking", context.Background(),
					query.GetBooking{BookingID: nonExistingID},
				).Return(nil, domain.ErrBookingNotFound{BookingID: nonExistingID})
			},
			want:    &api.GetBookingResponse{Booking: BookingFromDomain(booking)},
			wantErr: codes.NotFound.String(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			m := mocks{app: application.NewMockApp(t)}
			s := server{app: m.app}
			if tc.on != nil {
				tc.on(m)
			}
			// when
			got, err := s.GetBooking(tc.args.ctx, tc.args.req)
			// then
			mock.AssertExpectationsForObjects(t, m.app)

			if tc.wantErr != "" {
				assert.Containsf(t, err.Error(),
					tc.wantErr, "GetBooking() error = %v, wantErr %v", err, tc.wantErr,
				)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestServer_CreateBooking(t *testing.T) {
	type args struct {
		ctx context.Context
		req *api.CreateBookingRequest
	}
	booking, err := bootstrap.NewBooking("campsite-id")
	assert.NoError(t, err)

	tests := map[string]struct {
		args    args
		on      func(f mocks)
		wantErr string
	}{
		"Success": {
			args: args{
				ctx: context.Background(),
				req: &api.CreateBookingRequest{
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
				).Return(nil)
			},
			wantErr: "",
		},
		"Error_FailedPrecondition_BookingDatesNotAvailable": {
			args: args{
				ctx: context.Background(),
				req: &api.CreateBookingRequest{
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
					domain.ErrBookingDatesNotAvailable{
						StartDate: booking.StartDate,
						EndDate:   booking.EndDate,
					},
				)
			},
			wantErr: codes.FailedPrecondition.String(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			m := mocks{app: application.NewMockApp(t)}
			s := server{app: m.app}
			if tc.on != nil {
				tc.on(m)
			}
			// when
			resp, err := s.CreateBooking(tc.args.ctx, tc.args.req)
			// then
			mock.AssertExpectationsForObjects(t, m.app)

			if tc.wantErr != "" {
				assert.Containsf(t, err.Error(),
					tc.wantErr, "CreateBooking() error = %v, wantErr %v", err, tc.wantErr,
				)
				return
			}
			assert.NotEmpty(t, resp.BookingId)
		})
	}
}

func TestServer_UpdateBooking(t *testing.T) {
	type args struct {
		ctx context.Context
		req *api.UpdateBookingRequest
	}
	booking, err := bootstrap.NewBooking("campsite-id")
	assert.NoError(t, err)

	tests := map[string]struct {
		args    args
		on      func(f mocks)
		want    *api.UpdateBookingResponse
		wantErr string
	}{
		"Success": {
			args: args{
				ctx: context.Background(),
				req: &api.UpdateBookingRequest{Booking: BookingFromDomain(booking)},
			},
			on: func(f mocks) {
				f.app.On(
					"UpdateBooking", context.Background(), mock.Anything,
				).Return(nil)
			},
			want:    &api.UpdateBookingResponse{},
			wantErr: "",
		},
		"Error_FailedPrecondition_BookingDatesNotAvailable": {
			args: args{
				ctx: context.Background(),
				req: &api.UpdateBookingRequest{Booking: BookingFromDomain(booking)},
			},
			on: func(f mocks) {
				f.app.On(
					"UpdateBooking", context.Background(), mock.Anything,
				).Return(
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
			m := mocks{app: application.NewMockApp(t)}
			s := server{app: m.app}
			if tc.on != nil {
				tc.on(m)
			}
			// when
			got, err := s.UpdateBooking(tc.args.ctx, tc.args.req)
			// then
			mock.AssertExpectationsForObjects(t, m.app)

			if tc.wantErr != "" {
				assert.Containsf(t, err.Error(),
					tc.wantErr, "UpdateBooking() error = %v, wantErr %v", err, tc.wantErr,
				)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestServer_CancelBooking(t *testing.T) {
	type args struct {
		ctx context.Context
		req *api.CancelBookingRequest
	}
	booking, err := bootstrap.NewBooking("campsite-id")
	assert.NoError(t, err)

	tests := map[string]struct {
		args    args
		on      func(f mocks)
		want    *api.CancelBookingResponse
		wantErr string
	}{
		"Success": {
			args: args{
				ctx: context.Background(),
				req: &api.CancelBookingRequest{BookingId: booking.BookingID},
			},
			on: func(f mocks) {
				f.app.On(
					"CancelBooking", context.Background(), mock.Anything,
				).Return(nil)
			},
			want:    &api.CancelBookingResponse{},
			wantErr: "",
		},
		"Error_FailedPrecondition_BookingAlreadyCancelled": {
			args: args{
				ctx: context.Background(),
				req: &api.CancelBookingRequest{BookingId: booking.BookingID},
			},
			on: func(f mocks) {
				f.app.On(
					"CancelBooking", context.Background(), mock.Anything,
				).Return(
					domain.ErrBookingAlreadyCancelled{BookingID: booking.BookingID},
				)
			},
			want:    nil,
			wantErr: codes.FailedPrecondition.String(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			m := mocks{app: application.NewMockApp(t)}
			s := server{app: m.app}
			if tc.on != nil {
				tc.on(m)
			}
			// when
			got, err := s.CancelBooking(tc.args.ctx, tc.args.req)
			// then
			mock.AssertExpectationsForObjects(t, m.app)

			if tc.wantErr != "" {
				assert.Containsf(t, err.Error(),
					tc.wantErr, "CancelBooking() error = %v, wantErr %v", err, tc.wantErr,
				)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestServer_GetVacantDates(t *testing.T) {
	type args struct {
		ctx context.Context
		req *api.GetVacantDatesRequest
	}
	tests := map[string]struct {
		args    args
		on      func(f mocks)
		want    *api.GetVacantDatesResponse
		wantErr error
	}{
		"Success": {
			args: args{
				ctx: context.Background(),
				req: &api.GetVacantDatesRequest{
					CampsiteId: "campsite-id",
					StartDate:  "2006-01-02",
					EndDate:    "2006-01-03",
				},
			},
			on: func(f mocks) {
				f.app.On(
					"GetVacantDates", context.Background(), mock.Anything,
				).Return([]string{"2006-01-02"}, nil)
			},
			want:    &api.GetVacantDatesResponse{VacantDates: []string{"2006-01-02"}},
			wantErr: nil,
		},
		"Error_CommitTx": {
			args: args{
				ctx: context.Background(),
				req: &api.GetVacantDatesRequest{
					CampsiteId: "campsite-id",
					StartDate:  "2006-01-02",
					EndDate:    "2006-01-03",
				},
			},
			on: func(f mocks) {
				f.app.On(
					"GetVacantDates", context.Background(), mock.Anything,
				).Return(nil, bootstrap.ErrCommitTx)
			},
			want:    nil,
			wantErr: bootstrap.ErrCommitTx,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			m := mocks{app: application.NewMockApp(t)}
			s := server{app: m.app}
			if tc.on != nil {
				tc.on(m)
			}
			// when
			got, err := s.GetVacantDates(tc.args.ctx, tc.args.req)
			// then
			mock.AssertExpectationsForObjects(t, m.app)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr, "GetVacantDates() error = %v, wantErr %v",
					err, tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
