package command

import (
	"context"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
)

type (
	CreateCampsite struct {
		CampsiteID    string
		CampsiteCode  string
		Capacity      int32
		DrinkingWater bool
		Restrooms     bool
		PicnicTable   bool
		FirePit       bool
	}

	CreateCampsiteHandler struct {
		campsites domain.CampsiteRepository
	}
)

func NewCreateCampsiteHandler(campsites domain.CampsiteRepository) CreateCampsiteHandler {
	return CreateCampsiteHandler{campsites: campsites}
}

func (h CreateCampsiteHandler) Handle(ctx context.Context, cmd CreateCampsite) error {
	campsite := domain.Campsite{
		CampsiteID:    cmd.CampsiteID,
		CampsiteCode:  cmd.CampsiteCode,
		Capacity:      cmd.Capacity,
		DrinkingWater: cmd.DrinkingWater,
		Restrooms:     cmd.Restrooms,
		PicnicTable:   cmd.PicnicTable,
		FirePit:       cmd.FirePit,
		Active:        true,
	}
	return h.campsites.Insert(ctx, &campsite)
}
