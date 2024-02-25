package domain

type Campsite struct {
	// Persistence ID
	ID int64
	// Business ID
	CampsiteID    string
	Capacity      int32
	DrinkingWater bool
	Restrooms     bool
	PicnicTable   bool
	FirePit       bool
}

type CampsiteBuilder struct {
	campsite *Campsite
}

func NewCampsiteBuilder() *CampsiteBuilder {
	campsite := &Campsite{}
	b := &CampsiteBuilder{campsite: campsite}
	return b
}

func (b *CampsiteBuilder) ID(ID int64) *CampsiteBuilder {
	b.campsite.ID = ID
	return b
}

func (b *CampsiteBuilder) CampsiteId(campsiteID string) *CampsiteBuilder {
	b.campsite.CampsiteID = campsiteID
	return b
}

func (b *CampsiteBuilder) Capacity(capacity int32) *CampsiteBuilder {
	b.campsite.Capacity = capacity
	return b
}

func (b *CampsiteBuilder) DrinkingWater(drinkingWater bool) *CampsiteBuilder {
	b.campsite.DrinkingWater = drinkingWater
	return b
}

func (b *CampsiteBuilder) Restrooms(restrooms bool) *CampsiteBuilder {
	b.campsite.Restrooms = restrooms
	return b
}

func (b *CampsiteBuilder) PicnicTable(picnicTable bool) *CampsiteBuilder {
	b.campsite.PicnicTable = picnicTable
	return b
}

func (b *CampsiteBuilder) FirePit(firePit bool) *CampsiteBuilder {
	b.campsite.FirePit = firePit
	return b
}

func (b *CampsiteBuilder) Build() (*Campsite, error) {
	return b.campsite, nil
}
