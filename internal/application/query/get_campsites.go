package query

import (
	"context"
	"log/slog"

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

func NewGetCampsitesHandler(
	campsites domain.CampsiteRepository,
	logger *slog.Logger,
) GetCampsitesHandler {
	return decorator.ApplyQueryDecorator[GetCampsites, []*domain.Campsite](
		getCampsitesHandler{campsites: campsites},
		logger,
	)
}

func (h getCampsitesHandler) Handle(
	ctx context.Context,
	_ GetCampsites,
) ([]*domain.Campsite, error) {
	return h.campsites.FindAll(ctx)
}
