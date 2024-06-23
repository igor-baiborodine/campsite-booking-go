package domain

import (
	"encoding/json"
	"fmt"
)

type Campsite struct {
	// Persistence ID
	ID int64
	// Business ID
	CampsiteID    string
	CampsiteCode  string
	Capacity      int32
	DrinkingWater bool
	Restrooms     bool
	PicnicTable   bool
	FirePit       bool
	Active        bool
}

func (c *Campsite) String() string {
	result, _ := json.Marshal(c)
	return fmt.Sprintf("%s", result)
}
