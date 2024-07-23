package query

import (
	"context"

	"github.com/igor-baiborodine/campsite-booking-go/internal/application/decorator"
	"github.com/igor-baiborodine/campsite-booking-go/internal/application/handler"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
)

type (
	GetCampsites struct{}

	// GetCampsitesHandler is a logging decorator for the getCampsitesHandler struct.
	GetCampsitesHandler handler.Query[GetCampsites, []*domain.Campsite]

	getCampsitesHandler struct {
		campsites domain.CampsiteRepository
	}
)

func NewGetCampsitesHandler(campsites domain.CampsiteRepository) GetCampsitesHandler {
	return decorator.ApplyQueryDecorator[GetCampsites, []*domain.Campsite](
		getCampsitesHandler{campsites: campsites},
	)
}

func (h getCampsitesHandler) Handle(
	ctx context.Context,
	_ GetCampsites,
) ([]*domain.Campsite, error) {
	return h.campsites.FindAll(ctx)
}
