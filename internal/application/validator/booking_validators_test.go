package validator

import (
	"testing"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	"github.com/stretchr/testify/assert"
)

func TestBookingAllowedStartDate_Validate(t *testing.T) {
	now := bootstrap.AsStartOfDayUTC(time.Now())

	tests := map[string]struct {
		booking *domain.Booking
		wantErr error
	}{
		"Success_OneDayFromNow": {
			booking: &domain.Booking{StartDate: now.AddDate(0, 0, 1)},
			wantErr: nil,
		},
		"Success_OneMonthFromNow": {
			booking: &domain.Booking{StartDate: now.AddDate(0, 1, 0)},
			wantErr: nil,
		},
		"Error_Now": {
			booking: &domain.Booking{StartDate: now},
			wantErr: ErrBookingAllowedStartDate{},
		},
		"Error_OneMonthAndOneDayFromNow": {
			booking: &domain.Booking{StartDate: now.AddDate(0, 1, 1)},
			wantErr: ErrBookingAllowedStartDate{},
		},
	}
	v := BookingAllowedStartDate{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given

			// when
			err := v.Validate(tc.booking)
			// then
			assert.Equal(t, tc.wantErr, err,
				"BookingAllowedStartDate.Validate() error = %v, wantErr %v",
				err, tc.wantErr)
		})
	}
}

func TestBookingMaximumStay_Validate(t *testing.T) {
	now := bootstrap.AsStartOfDayUTC(time.Now())

	tests := map[string]struct {
		booking *domain.Booking
		wantErr error
	}{
		"Success_OneDay": {
			booking: &domain.Booking{
				StartDate: now.AddDate(0, 0, 1),
				EndDate:   now.AddDate(0, 0, 2),
			},
			wantErr: nil,
		},
		"Success_ThreeDays": {
			booking: &domain.Booking{
				StartDate: now.AddDate(0, 0, 1),
				EndDate:   now.AddDate(0, 0, 4),
			},
			wantErr: nil,
		},
		"Error_FourDays": {
			booking: &domain.Booking{
				StartDate: now.AddDate(0, 0, 1),
				EndDate:   now.AddDate(0, 0, 5),
			},
			wantErr: ErrBookingMaximumStay{},
		},
	}
	v := BookingMaximumStay{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given

			// when
			err := v.Validate(tc.booking)
			// then
			assert.Equal(t, tc.wantErr, err,
				"BookingMaximumStay.Validate() error = %v, wantErr %v",
				err, tc.wantErr)
		})
	}
}

func TestBookingStartDateBeforeEndDate_Validate(t *testing.T) {
	now := bootstrap.AsStartOfDayUTC(time.Now())

	tests := map[string]struct {
		booking *domain.Booking
		wantErr error
	}{
		"Success": {
			booking: &domain.Booking{
				StartDate: now.AddDate(0, 0, 1),
				EndDate:   now.AddDate(0, 0, 2),
			},
			wantErr: nil,
		},
		"Error_StartDateEqualsEndDate": {
			booking: &domain.Booking{
				StartDate: now.AddDate(0, 0, 2),
				EndDate:   now.AddDate(0, 0, 2),
			},
			wantErr: ErrBookingStartDateBeforeEndDate{},
		},
		"Error_StartDateAfterEndDate": {
			booking: &domain.Booking{
				StartDate: now.AddDate(0, 0, 2),
				EndDate:   now.AddDate(0, 0, 1),
			},
			wantErr: ErrBookingStartDateBeforeEndDate{},
		},
	}
	v := BookingStartDateBeforeEndDate{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given

			// when
			err := v.Validate(tc.booking)
			// then
			assert.Equal(t, tc.wantErr, err,
				"BookingStartDateBeforeEndDate.Validate() error = %v, wantErr %v",
				err, tc.wantErr)
		})
	}
}

func TestApply(t *testing.T) {
	now := bootstrap.AsStartOfDayUTC(time.Now())

	tests := map[string]struct {
		booking *domain.Booking
		wantErr error
	}{
		"Success": {
			booking: &domain.Booking{
				StartDate: now.AddDate(0, 0, 1),
				EndDate:   now.AddDate(0, 0, 2),
			},
			wantErr: nil,
		},
		"Error_BookingStartDateBeforeEndDateValidator": {
			booking: &domain.Booking{
				StartDate: now.AddDate(0, 0, 2),
				EndDate:   now.AddDate(0, 0, 1),
			},
			wantErr: domain.ErrBookingValidation{
				MultiErr: multierror.Append(ErrBookingStartDateBeforeEndDate{}),
			},
		},
		"Error_BookingAllowedStartDate_ErrBookingMaximumStay": {
			booking: &domain.Booking{
				StartDate: now.AddDate(0, 2, 2),
				EndDate:   now.AddDate(0, 4, 2),
			},
			wantErr: domain.ErrBookingValidation{
				MultiErr: multierror.Append(
					multierror.Append(ErrBookingAllowedStartDate{}),
					ErrBookingMaximumStay{},
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			validators := []domain.BookingValidator{
				BookingStartDateBeforeEndDate{},
				BookingAllowedStartDate{},
				BookingMaximumStay{},
			}
			// when
			err := Apply(validators, tc.booking)
			// then
			assert.Equal(t, tc.wantErr, err,
				"Apply() error = %v, wantErr %v", err, tc.wantErr)
		})
	}
}
