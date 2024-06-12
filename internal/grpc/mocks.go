package grpc

import (
	"github.com/igor-baiborodine/campsite-booking-go/internal/application"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
)

type mocks struct {
	app       *application.MockApp
	campsites *domain.MockCampsiteRepository
	bookings  *domain.MockBookingRepository
}
