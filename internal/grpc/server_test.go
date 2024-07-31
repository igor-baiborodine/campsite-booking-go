//go:build !integration

package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/go-multierror"
	api "github.com/igor-baiborodine/campsite-booking-go/campgroundspb/v1"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/query"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/validator"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mocks struct {
	app *application.MockApp
}

func TestServer_GetCampsites(t *testing.T) {
	campsite, err := bootstrap.NewCampsite()
	assert.NoError(t, err)

	tests := map[string]struct {
		req     *api.GetCampsitesRequest
		on      func(f mocks)
		want    *api.GetCampsitesResponse
		wantErr error
	}{
		"Success": {
			req: &api.GetCampsitesRequest{},
			on: func(f mocks) {
				f.app.
					On("GetCampsites", mock.Anything, mock.Anything).
					Return([]*domain.Campsite{campsite}, nil)
			},
			want: &api.GetCampsitesResponse{
				Campsites: []*api.Campsite{CampsiteFromDomain(campsite)},
			},
			wantErr: nil,
		},
		"Error_ErrQuery": {
			req: &api.GetCampsitesRequest{},
			on: func(f mocks) {
				f.app.
					On("GetCampsites", mock.Anything, mock.Anything).
					Return(nil, bootstrap.ErrQuery)
			},
			want:    nil,
			wantErr: bootstrap.ErrQuery,
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
			got, err := s.GetCampsites(context.TODO(), tc.req)
			// then
			assert.Equal(t, tc.want, got,
				"GetCampsites() got = %v, want %v", got, tc.want)
			assert.Equal(t, tc.wantErr, err,
				"GetCampsites() error = %v, wantErr %v", err, tc.wantErr)
			mock.AssertExpectationsForObjects(t, m.app)
		})
	}
}

func TestServer_CreateCampsite(t *testing.T) {
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
		req     *api.CreateCampsiteRequest
		on      func(f mocks)
		wantErr error
	}{
		"Success": {
			req: req,
			on: func(f mocks) {
				f.app.
					On("CreateCampsite", context.TODO(), mock.Anything).
					Return(nil)
			},
			wantErr: nil,
		},
		"Error_CommitTx": {
			req: req,
			on: func(f mocks) {
				f.app.
					On("CreateCampsite", context.TODO(), mock.Anything).
					Return(bootstrap.ErrCommitTx)
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
			got, err := s.CreateCampsite(context.TODO(), req)
			// then
			defer mock.AssertExpectationsForObjects(t, m.app)

			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err,
					"CreateCampsite() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			assert.NotEmpty(t, got.CampsiteId)
		})
	}
}

func TestServer_GetBooking(t *testing.T) {
	booking, err := bootstrap.NewBooking("campsite-id")
	assert.NoError(t, err)
	nonExistingID := "non-existing-id"
	errBookingNotFound := domain.ErrBookingNotFound{BookingID: nonExistingID}

	tests := map[string]struct {
		req     *api.GetBookingRequest
		on      func(f mocks)
		want    *api.GetBookingResponse
		wantErr error
	}{
		"Success": {
			req: &api.GetBookingRequest{BookingId: booking.BookingID},
			on: func(f mocks) {
				f.app.
					On(
						"GetBooking",
						context.TODO(),
						query.GetBooking{BookingID: booking.BookingID},
					).
					Return(booking, nil)
			},
			want:    &api.GetBookingResponse{Booking: BookingFromDomain(booking)},
			wantErr: nil,
		},
		"Error_NotFound_ErrBookingNotFound": {
			req: &api.GetBookingRequest{BookingId: nonExistingID},
			on: func(f mocks) {
				f.app.
					On("GetBooking", context.TODO(), query.GetBooking{BookingID: nonExistingID}).
					Return(nil, errBookingNotFound)
			},
			want:    nil,
			wantErr: status.Error(codes.NotFound, errBookingNotFound.Error()),
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
			got, err := s.GetBooking(context.TODO(), tc.req)
			// then
			assert.Equal(t, tc.want, got,
				"GetBooking() got = %v, want %v", got, tc.want)
			assert.Equal(t, tc.wantErr, err,
				"GetBooking() error = %v, wantErr %v", err, tc.wantErr)
			mock.AssertExpectationsForObjects(t, m.app)
		})
	}
}

func TestServer_CreateBooking(t *testing.T) {
	booking, err := bootstrap.NewBooking("campsite-id")
	assert.NoError(t, err)
	errBookingDatesNotAvailable := domain.ErrBookingDatesNotAvailable{
		StartDate: booking.StartDate,
		EndDate:   booking.EndDate,
	}
	req := &api.CreateBookingRequest{
		CampsiteId: booking.CampsiteID,
		Email:      booking.Email,
		FullName:   booking.FullName,
		StartDate:  booking.StartDate.Format(time.DateOnly),
		EndDate:    booking.EndDate.Format(time.DateOnly),
	}

	tests := map[string]struct {
		req     *api.CreateBookingRequest
		on      func(f mocks)
		wantErr error
	}{
		"Success": {
			req: req,
			on: func(f mocks) {
				f.app.
					On("CreateBooking", context.TODO(), mock.Anything).
					Return(nil)
			},
			wantErr: nil,
		},
		"Error_FailedPrecondition_BookingDatesNotAvailable": {
			req: req,
			on: func(f mocks) {
				f.app.
					On("CreateBooking", context.TODO(), mock.Anything).
					Return(errBookingDatesNotAvailable)
			},
			wantErr: status.Error(codes.FailedPrecondition, errBookingDatesNotAvailable.Error()),
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
			got, err := s.CreateBooking(context.TODO(), tc.req)
			// then
			defer mock.AssertExpectationsForObjects(t, m.app)

			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err,
					"CreateBooking() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			assert.NotEmpty(t, got.BookingId)
		})
	}
}

func TestServer_UpdateBooking(t *testing.T) {
	booking, err := bootstrap.NewBooking("campsite-id")
	assert.NoError(t, err)
	errBookingDatesNotAvailable := domain.ErrBookingDatesNotAvailable{
		StartDate: booking.StartDate,
		EndDate:   booking.EndDate,
	}
	errBookingValidation := domain.ErrBookingValidation{
		MultiErr: multierror.Append(validator.ErrBookingStartDateBeforeEndDate{}),
	}
	req := &api.UpdateBookingRequest{Booking: BookingFromDomain(booking)}

	tests := map[string]struct {
		req     *api.UpdateBookingRequest
		on      func(f mocks)
		want    *api.UpdateBookingResponse
		wantErr error
	}{
		"Success": {
			req: req,
			on: func(f mocks) {
				f.app.
					On("UpdateBooking", context.TODO(), mock.Anything).
					Return(nil)
			},
			want:    &api.UpdateBookingResponse{},
			wantErr: nil,
		},
		"Error_FailedPrecondition_BookingDatesNotAvailable": {
			req: req,
			on: func(f mocks) {
				f.app.
					On("UpdateBooking", context.TODO(), mock.Anything).
					Return(errBookingDatesNotAvailable)
			},
			want:    nil,
			wantErr: status.Error(codes.FailedPrecondition, errBookingDatesNotAvailable.Error()),
		},
		"Error_InvalidArgument_BookingValidation": {
			req: req,
			on: func(f mocks) {
				f.app.
					On("UpdateBooking", context.TODO(), mock.Anything).
					Return(errBookingValidation)
			},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, errBookingValidation.Error()),
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
			got, err := s.UpdateBooking(context.TODO(), tc.req)
			// then
			assert.Equal(t, tc.want, got,
				"UpdateBooking() got = %v, want %v", got, tc.want)
			assert.Equal(t, tc.wantErr, err,
				"UpdateBooking() error = %v, wantErr %v", err, tc.wantErr)
			mock.AssertExpectationsForObjects(t, m.app)
		})
	}
}

func TestServer_CancelBooking(t *testing.T) {
	booking, err := bootstrap.NewBooking("campsite-id")
	assert.NoError(t, err)
	errBookingAlreadyCancelled := domain.ErrBookingAlreadyCancelled{BookingID: booking.BookingID}
	req := &api.CancelBookingRequest{BookingId: booking.BookingID}

	tests := map[string]struct {
		req     *api.CancelBookingRequest
		on      func(f mocks)
		want    *api.CancelBookingResponse
		wantErr error
	}{
		"Success": {
			req: req,
			on: func(f mocks) {
				f.app.
					On("CancelBooking", context.TODO(), mock.Anything).
					Return(nil)
			},
			want:    &api.CancelBookingResponse{},
			wantErr: nil,
		},
		"Error_FailedPrecondition_BookingAlreadyCancelled": {
			req: req,
			on: func(f mocks) {
				f.app.
					On("CancelBooking", context.TODO(), mock.Anything).
					Return(errBookingAlreadyCancelled)
			},
			want:    nil,
			wantErr: status.Error(codes.FailedPrecondition, errBookingAlreadyCancelled.Error()),
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
			got, err := s.CancelBooking(context.TODO(), tc.req)
			// then
			assert.Equal(t, tc.want, got,
				"CancelBooking() got = %v, want %v", got, tc.want)
			assert.Equal(t, tc.wantErr, err,
				"CancelBooking() error = %v, wantErr %v", err, tc.wantErr)
			mock.AssertExpectationsForObjects(t, m.app)
		})
	}
}

func TestServer_GetVacantDates(t *testing.T) {
	req := &api.GetVacantDatesRequest{
		CampsiteId: "campsite-id",
		StartDate:  "2006-01-02",
		EndDate:    "2006-01-03",
	}

	tests := map[string]struct {
		req     *api.GetVacantDatesRequest
		on      func(f mocks)
		want    *api.GetVacantDatesResponse
		wantErr error
	}{
		"Success": {
			req: req,
			on: func(f mocks) {
				f.app.
					On("GetVacantDates", context.TODO(), mock.Anything).
					Return([]string{"2006-01-02"}, nil)
			},
			want:    &api.GetVacantDatesResponse{VacantDates: []string{"2006-01-02"}},
			wantErr: nil,
		},
		"Error_CommitTx": {
			req: req,
			on: func(f mocks) {
				f.app.
					On("GetVacantDates", context.TODO(), mock.Anything).
					Return(nil, bootstrap.ErrCommitTx)
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
			got, err := s.GetVacantDates(context.TODO(), tc.req)
			// then
			assert.Equal(t, tc.want, got,
				"GetVacantDates() got = %v, want %v", got, tc.want)
			assert.Equal(t, tc.wantErr, err,
				"GetVacantDates() error = %v, wantErr %v", err, tc.wantErr)
			mock.AssertExpectationsForObjects(t, m.app)
		})
	}
}
