package bootstrap

const (
	DeleteBookings     = "DELETE FROM bookings"
	DeleteCampsites    = "DELETE FROM campsites"
	SelectByCampsiteID = "SELECT campsite_code FROM campsites WHERE campsite_id = $1"
)
