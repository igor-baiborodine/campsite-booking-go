# Campsite Bookings API (Go)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Technical Task](#technical-task)
  - [Booking Constraints](#booking-constraints)
  - [System Requirements](#system-requirements)
- [Implementation Details](#implementation-details)
  - [Technology Stack](#technology-stack)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Technical Task

### Booking Constraints

* The campsite can be reserved for a maximum of 3 days.
* The campsite can be reserved a minimum of 1 day(s) ahead of arrival and up to 1 month in advance.
* Reservations can be canceled anytime.
* For the sake of simplicity, assume the check-in & check-out time is 12:00 AM.

### System Requirements

* The users will need to find out when the campsite is available. So, the system should expose an API
  to provide information on the availability of the campsite for a given date range, with the default
  being 1 month.
* Provide an endpoint for reserving the campsite. The user will provide his/her email & full name
  at the time of reserving the campsite along with the intended arrival date and departure date. Return
  a unique booking identifier to the caller if the reservation succeeds.
* The unique booking identifier can be used to modify or cancel the reservation later on. Provide
  appropriate endpoint (s) to allow modification/cancellation of an existing reservation.
* Due to the popularity of the campsite, there is a high likelihood of multiple users attempting to
  reserve the campsite for the same/overlapping date(s). Demonstrate with appropriate test cases
  that the system can gracefully handle concurrent requests to reserve the campsite.
* Provide appropriate error messages to the caller to indicate the error cases.
* The system should be able to handle a large volume of requests to determine campsite
  availability.
* There are no restrictions on how reservations are stored as long as system constraints are not
  violated.

## Implementation Details

### Technology Stack

* [Go](https://github.com/golang/go), [gRPC](https://github.com/grpc/grpc-go), 
* [protovalidate-go](https://github.com/bufbuild/protovalidate-go) (requests validation)
* [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports), [golines](https://github.com/segmentio/golines), [gofumpt](https://github.com/mvdan/gofumpt) (code style & formatting)
* [PostgreSQL](https://www.postgresql.org/)  
* [Docker](https://www.docker.com/), [Docker Compose](https://docs.docker.com/compose/)
* [Goose](https://pressly.github.io/goose/) (DB migrations)
* TODO

### Project Setup

**Prerequisites**:
1. [Git](https://git-scm.com/), see this [guide](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) on how to install Git.
2. [Make](https://man7.org/linux/man-pages/man1/make.1.html) 
3 [Go](https://go.dev/) (version >= 1.21) see this [guide](https://go.dev/doc/install) on how to install Go.

* Clone the project: 

```shell
$ git clone https://github.com/igor-baiborodine/campsite-booking-go.git
```

* Install necessary tools executing the following command from the project's root:
```shell
$ make install-tools
```

If you use either [IntelliJ IDEA](https://www.jetbrains.com/idea/) or [GoLand](https://www.jetbrains.com/go/) IDEs, 
follow this [guide](/readme/IDE-SETUP.md) to configure it. 
