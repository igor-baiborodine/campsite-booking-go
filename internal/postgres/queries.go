package postgres

const (
	InsertIntoCampsites = "INSERT INTO campsites.campsites " +
		"(campsite_id, campsite_code, capacity, restrooms, drinking_water, picnic_table, fire_pit, active, created_at, updated_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"

	FindAllInCampsites = "SELECT campsite_id, campsite_code, capacity, restrooms, drinking_water, picnic_table, fire_pit, active " +
		"FROM campsites.campsites"

	SelectByBookingIdFromBookings = "SELECT id, booking_id, campsite_id, email, full_name, start_date, end_date, active " +
		"FROM bookings.bookings " +
		"WHERE booking_id = $1"
)
