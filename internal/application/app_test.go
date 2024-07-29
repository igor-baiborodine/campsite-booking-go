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
	app := New(campsiteRepository, bookingRepository)
	// then
	assert.NotNil(t, app)
	assert.NotNil(t, app.commands.CreateCampsiteHandler)
	assert.NotNil(t, app.commands.CreateBookingHandler)
	assert.NotNil(t, app.commands.UpdateBookingHandler)
	assert.NotNil(t, app.commands.CancelBookingHandler)
	assert.NotNil(t, app.queries.GetCampsitesHandler)
	assert.NotNil(t, app.queries.GetBookingHandler)
	assert.NotNil(t, app.queries.GetVacantDatesHandler)
}
