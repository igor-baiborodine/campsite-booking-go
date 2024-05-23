package commands

import (
	"context"

	"github.com/google/uuid"
	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
)

type (
	CreateCampsite struct {
		CampsiteId    string
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
	campsite := domain.Campsite{}
	campsite.CampsiteID = uuid.New().String()
	campsite.CampsiteCode = cmd.CampsiteCode
	campsite.Capacity = cmd.Capacity
	campsite.DrinkingWater = cmd.DrinkingWater
	campsite.Restrooms = cmd.Restrooms
	campsite.PicnicTable = cmd.PicnicTable
	campsite.FirePit = cmd.FirePit
	campsite.Active = true

	return h.campsites.Insert(ctx, &campsite)
}
