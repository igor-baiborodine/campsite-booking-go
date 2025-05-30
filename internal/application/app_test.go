package application

import (
	"testing"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestApp_New(t *testing.T) {
	// given
	campsiteRepository := domain.NewMockCampsiteRepository(t)
	bookingRepository := domain.NewMockBookingRepository(t)
	// when
	got := New(campsiteRepository, bookingRepository)
	// then
	assert.NotNil(t, got)
	assert.NotNil(t, got.CreateCampsiteHandler)
	assert.NotNil(t, got.CreateBookingHandler)
	assert.NotNil(t, got.UpdateBookingHandler)
	assert.NotNil(t, got.CancelBookingHandler)
	assert.NotNil(t, got.GetCampsitesHandler)
	assert.NotNil(t, got.GetBookingHandler)
	assert.NotNil(t, got.GetVacantDatesHandler)
}
