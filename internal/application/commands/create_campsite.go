package commands

import (
	"context"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
)

type (
	CreateCampsite struct {
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

func (h CreateCampsiteHandler) CreateCampsite(ctx context.Context, cmd CreateCampsite) error {
	campsiteBuilder := domain.NewCampsiteBuilder().
		CampsiteCode(cmd.CampsiteCode).
		Capacity(cmd.Capacity).
		DrinkingWater(cmd.DrinkingWater).
		Restrooms(cmd.Restrooms).
		PicnicTable(cmd.PicnicTable).
		FirePit(cmd.FirePit)

	campsite, err := campsiteBuilder.Build()
	if err != nil {
		return err
	}
	return h.campsites.Insert(ctx, campsite)
}
