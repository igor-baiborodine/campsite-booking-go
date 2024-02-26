package queries

import (
	"context"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
)

type (
	GetCampsites struct{}

	GetCampsitesHandler struct {
		campsites domain.CampsiteRepository
	}
)

func NewGetCampsitesHandler(campsites domain.CampsiteRepository) GetCampsitesHandler {
	return GetCampsitesHandler{campsites: campsites}
}

func (h GetCampsitesHandler) GetCampsites(ctx context.Context, _ GetCampsites) ([]*domain.Campsite, error) {
	return h.campsites.FindAll(ctx)
}
