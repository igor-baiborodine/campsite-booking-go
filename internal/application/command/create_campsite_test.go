package command

import (
	"context"
	"testing"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCampsiteHandler(t *testing.T) {
	type mocks struct {
		campsites *domain.MockCampsiteRepository
	}
	campsite, err := bootstrap.NewCampsite()
	if err != nil {
		t.Fatalf("create campsite error: %v", err)
	}
	campsite.ID = 0
	campsite.Active = true

	cmd := CreateCampsite{
		CampsiteID:    campsite.CampsiteID,
		CampsiteCode:  campsite.CampsiteCode,
		Capacity:      campsite.Capacity,
		DrinkingWater: campsite.DrinkingWater,
		Restrooms:     campsite.Restrooms,
		PicnicTable:   campsite.PicnicTable,
		FirePit:       campsite.FirePit,
	}

	tests := map[string]struct {
		cmd     CreateCampsite
		on      func(f mocks)
		wantErr error
	}{
		"Success": {
			cmd: cmd,
			on: func(f mocks) {
				f.campsites.On(
					"Insert", context.TODO(), campsite,
				).Return(nil)
			},
			wantErr: nil,
		},
		"Error_CommitTx": {
			cmd: cmd,
			on: func(f mocks) {
				f.campsites.On(
					"Insert", context.TODO(), campsite,
				).Return(bootstrap.ErrCommitTx)
			},
			wantErr: bootstrap.ErrCommitTx,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			m := mocks{
				campsites: domain.NewMockCampsiteRepository(t),
			}
			h := NewCreateCampsiteHandler(m.campsites)
			if tc.on != nil {
				tc.on(m)
			}
			// when
			err := h.Handle(context.TODO(), tc.cmd)
			// then
			assert.ErrorIs(t, err, tc.wantErr,
				"CreateCampsiteHandler.Handle() error = %v, wantErr %v", err, tc.wantErr)
			mock.AssertExpectationsForObjects(t, m.campsites)
		})
	}
}
