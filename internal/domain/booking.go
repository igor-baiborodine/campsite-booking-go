package domain

type Booking struct {
	// Persistence ID
	ID int64
	// Business ID
	BookingID  string
	CampsiteID string
	Email      string
	FullName   string
	StartDate  string
	EndDate    string
	Active     bool
}

type BookingBuilder struct {
	booking *Booking
}

func NewBookingBuilder() *BookingBuilder {
	booking := &Booking{}
	b := &BookingBuilder{booking: booking}
	return b
}

func (b *BookingBuilder) ID(ID int64) *BookingBuilder {
	b.booking.ID = ID
	return b
}

func (b *BookingBuilder) BookingID(bookingID string) *BookingBuilder {
	b.booking.BookingID = bookingID
	return b
}

func (b *BookingBuilder) CampsiteID(campsiteID string) *BookingBuilder {
	b.booking.CampsiteID = campsiteID
	return b
}

func (b *BookingBuilder) Email(email string) *BookingBuilder {
	b.booking.Email = email
	return b
}

func (b *BookingBuilder) FullName(fullName string) *BookingBuilder {
	b.booking.FullName = fullName
	return b
}

func (b *BookingBuilder) StartDate(startDate string) *BookingBuilder {
	b.booking.StartDate = startDate
	return b
}

func (b *BookingBuilder) EndDate(endDate string) *BookingBuilder {
	b.booking.EndDate = endDate
	return b
}

func (b *BookingBuilder) Active(active bool) *BookingBuilder {
	b.booking.Active = active
	return b
}

func (b *BookingBuilder) Build() (*Booking, error) {
	return b.booking, nil
}
