package validators

import (
	"testing"
	"time"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestBookingAllowedStartDateValidator_Validate(t *testing.T) {
	now := time.Now()

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
