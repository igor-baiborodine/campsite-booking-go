package query

import (
	"context"
	"testing"

	"github.com/igor-baiborodine/campsite-booking-go/internal/domain"
	"github.com/igor-baiborodine/campsite-booking-go/internal/testing/bootstrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCampsitesHandler(t *testing.T) {
	type mocks struct {
		campsites *domain.MockCampsiteRepository
	}
	campsite, err := bootstrap.NewCampsite()
	if err != nil {
		t.Fatalf("create campsite error: %v", err)
	}

	tests := map[string]struct {
		qry     GetCampsites
		on      func(f mocks)
		want    []*domain.Campsite
		wantErr error
	}{
		"Success": {
			qry: GetCampsites{},
			on: func(f mocks) {
				f.campsites.
					On("FindAll", context.TODO()).
					Return([]*domain.Campsite{campsite}, nil)
			},
			want:    []*domain.Campsite{campsite},
			wantErr: nil,
		},
		"Success_NoCampsitesFound": {
			qry: GetCampsites{},
			on: func(f mocks) {
				f.campsites.
					On("FindAll", context.TODO()).
					Return(nil, nil)
			},
			want:    nil,
			wantErr: nil,
		},
		"Error_BeginTx": {
			qry: GetCampsites{},
			on: func(f mocks) {
				f.campsites.
					On("FindAll", context.TODO()).
					Return(nil, bootstrap.ErrBeginTx)
			},
			want:    nil,
			wantErr: bootstrap.ErrBeginTx,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given
			m := mocks{
				campsites: domain.NewMockCampsiteRepository(t),
			}
			h := NewGetCampsitesHandler(m.campsites)
			if tc.on != nil {
				tc.on(m)
			}
			// when
			got, err := h.Handle(context.TODO(), tc.qry)
			// then
			assert.Equalf(t, tc.want, got,
				"GetCampsitesHandler.Handle() got = %v, want %v", got, tc.want)
			assert.Equalf(t, tc.wantErr, err,
				"Find() error = %v, wantErr %v", err, tc.wantErr)
			mock.AssertExpectationsForObjects(t, m.campsites)
		})
	}
}
