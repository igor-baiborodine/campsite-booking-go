package sql

const (
	InsertCampsiteQuery = `
		INSERT INTO campsites (
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

	FindAllCampsitesQuery = `
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
		FROM campsites
	`

	FindBookingByBookingIDQuery = `
		SELECT 
		    id,
		    booking_id, 
		    campsite_id, 
		    email, 
		    full_name, 
		    start_date, 
		    end_date, 
		    active
		FROM bookings
		WHERE booking_id = $1
	`

	InsertBookingQuery = `
		INSERT INTO bookings (
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

	FindAllBookingsForDateRangeQuery = `
		SELECT
		    id,
		    booking_id, 
		    campsite_id, 
		    email, 
		    full_name, 
		    start_date, 
		    end_date, 
		    active
		FROM bookings
		WHERE active = TRUE 
		  	AND campsite_id = $1
		  	AND ((start_date < $2 AND $3 < end_date) 
		            OR ($2 < end_date AND end_date <= $3) 
		            OR ($2 <= start_date AND start_date <= $3)) 
	`

	UpdateBookingQuery = `
		UPDATE bookings
		SET 
		    campsite_id = $2, 
		    email = $3, 
		    full_name = $4, 
		    start_date = $5,
		    end_date = $6,
		    active = $7,
		    updated_at = $8
		WHERE booking_id = $1
	`
)
