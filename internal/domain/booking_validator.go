package domain

type BookingValidator interface {
	Validate(b *Booking) error
}
