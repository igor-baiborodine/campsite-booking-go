package domain

import (
	"context"
)

type CampsiteRepository interface {
	FindAll(ctx context.Context) ([]*Campsite, error)
	Insert(ctx context.Context, campsite *Campsite) error
}
