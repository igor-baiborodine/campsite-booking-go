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
	assert.NotNil(t, got.commands.CreateCampsiteHandler)
	assert.NotNil(t, got.commands.CreateBookingHandler)
	assert.NotNil(t, got.commands.UpdateBookingHandler)
	assert.NotNil(t, got.commands.CancelBookingHandler)
	assert.NotNil(t, got.queries.GetCampsitesHandler)
	assert.NotNil(t, got.queries.GetBookingHandler)
	assert.NotNil(t, got.queries.GetVacantDatesHandler)
}
