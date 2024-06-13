package validators

import (
	"testing"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestBookingAllowedStartDateValidator_Validate(t *testing.T) {
	now := bootstrap.AsStartOfDayUTC(time.Now())

	var tests = map[string]struct {
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
	v := BookingAllowedStartDateValidator{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given

			// when
			err := v.Validate(tc.booking)
			// then
			if tc.wantErr != nil {
				assert.ErrorIs(
					t, err, tc.wantErr,
					"BookingAllowedStartDateValidator.Validate() error = %v, wantErr %v",
					err, tc.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestBookingMaximumStayValidator_Validate(t *testing.T) {
	now := bootstrap.AsStartOfDayUTC(time.Now())

	var tests = map[string]struct {
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
	v := BookingMaximumStayValidator{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given

			// when
			err := v.Validate(tc.booking)
			// then
			if tc.wantErr != nil {
				assert.ErrorIs(
					t, err, tc.wantErr,
					"BookingMaximumStayValidator.Validate() error = %v, wantErr %v",
					err, tc.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestBookingStartDateBeforeEndDateValidator_Validate(t *testing.T) {
	now := bootstrap.AsStartOfDayUTC(time.Now())

	var tests = map[string]struct {
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
	v := BookingStartDateBeforeEndDateValidator{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given

			// when
			err := v.Validate(tc.booking)
			// then
			if tc.wantErr != nil {
				assert.ErrorIs(
					t, err, tc.wantErr,
					"BookingStartDateBeforeEndDateValidator.Validate() error = %v, wantErr %v",
					err, tc.wantErr)
				return
			}
			assert.Nil(t, err)
		})
	}
}

func TestApply(t *testing.T) {
	now := bootstrap.AsStartOfDayUTC(time.Now())

	var tests = map[string]struct {
		booking  *domain.Booking
		wantErrs []string
	}{
		"Success": {
			booking: &domain.Booking{
				StartDate: now.AddDate(0, 0, 1),
				EndDate:   now.AddDate(0, 0, 2),
			},
			wantErrs: nil,
		},
		"Error_BookingStartDateBeforeEndDateValidator": {
			booking: &domain.Booking{
				StartDate: now.AddDate(0, 0, 2),
				EndDate:   now.AddDate(0, 0, 1),
			},
			wantErrs: []string{
				codes.InvalidArgument.String(),
				ErrBookingStartDateBeforeEndDate{}.Error(),
			},
		},
		"Error_BookingAllowedStartDate_ErrBookingMaximumStay": {
			booking: &domain.Booking{
				StartDate: now.AddDate(0, 2, 2),
				EndDate:   now.AddDate(0, 4, 2),
			},
			wantErrs: []string{
				codes.InvalidArgument.String(),
				ErrBookingAllowedStartDate{}.Error(),
				ErrBookingMaximumStay{}.Error(),
			},
		},
	}

	validators := []BookingValidator{
		&BookingStartDateBeforeEndDateValidator{},
		&BookingAllowedStartDateValidator{},
		&BookingMaximumStayValidator{},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given

			// when
			err := Apply(validators, tc.booking)
			// then
			if tc.wantErrs != nil {
				for _, wantErr := range tc.wantErrs {
					assert.Containsf(t, err.Error(), wantErr, "Apply() error = %v, wantErr %v", err, wantErr)
				}
				return
			}
			assert.Nil(t, err)
		})
	}
}
