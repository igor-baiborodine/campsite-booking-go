package postgres

const (
	InsertIntoCampsites = `
		INSERT INTO campgrounds.campsites (
		                                   campsite_id, 
		                                   campsite_code, 
		                                   capacity, 
		                                   restrooms, 
		                                   drinking_water, 
		                                   picnic_table, 
		                                   fire_pit, 
		                                   active, 
		                                   created_at, 
		                                   updated_at
		                                   )
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	FindAllInCampsites = `
		SELECT 
		    id,
		    campsite_id, 
		    campsite_code, 
		    capacity, 
		    restrooms, 
		    drinking_water, 
		    picnic_table, 
		    fire_pit, 
		    active
		FROM campgrounds.campsites
	`

	SelectByBookingIdFromBookings = `
		SELECT 
		    id,
		    booking_id, 
		    campsite_id, 
		    email, 
		    full_name, 
		    start_date, 
		    end_date, 
		    active
		FROM campgrounds.bookings
		WHERE booking_id = $1
	`

	InsertIntoBookings = `
		INSERT INTO campgrounds.bookings (
		                                  booking_id, 
		                                  campsite_id, 
		                                  email, 
		                                  full_name, 
		                                  start_date, 
		                                  end_date, 
		                                  active, 
		                                  created_at, 
		                                  updated_at
		                                  )
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	FindForDateRangeInBookings = `
		SELECT
		    id,
		    booking_id, 
		    campsite_id, 
		    email, 
		    full_name, 
		    start_date, 
		    end_date, 
		    active
		FROM campgrounds.bookings
		WHERE ((start_date < $1 and $2 < end_date) 
		           OR ($1 < end_date and end_date <= $2) 
		           OR ($1 <= start_date and start_date <= $2)) 
		  AND active = TRUE and campsite_id = $3
	`
)
