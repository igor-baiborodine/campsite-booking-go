package domain

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
