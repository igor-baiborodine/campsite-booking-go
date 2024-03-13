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
		WHERE active = TRUE 
		  AND campsite_id = $1
		  AND ((start_date < $2 AND $3 < end_date) 
		           OR ($2 < end_date AND end_date <= $3) 
		           OR ($2 <= start_date AND start_date <= $3)) 
	`
)
