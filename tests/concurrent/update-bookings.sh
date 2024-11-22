#!/usr/bin/env bash

set -e

count="$1"
campsite_id="$2"
booking_id="$3"
start_date="$4"
end_date="$5"

data="{\"booking\": {\"campsite_id\": \"$campsite_id\", \"booking_id\": \"$booking_id\", \"start_date\": \"$start_date\", \"end_date\": \"$end_date\", \"email\": \"EMAIL\", \"full_name\": \"FULL_NAME\"}}"
i=1
requests=""

while [ "$i" -le "$count" ]; do
  request_data=$(echo "${data}" | sed -e "s/EMAIL/john.smith.$i@email.com/g" \
    -e "s/FULL_NAME/John Smith $i/g")
  requests+="grpcurl -plaintext -d '$request_data' localhost:8085 campgroundspb.v1.CampgroundsService/UpdateBooking & "
  i=$((i+1))
done

printf "✅ about to execute $count update request(s):\n  %s\n" "$requests"
eval "$requests"

sleep 1
printf "\n✅ concurrent bookings update completed\n"
