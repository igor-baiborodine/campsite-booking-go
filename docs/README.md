# Campsite Bookings API (Go)

![ci](https://github.com/igor-baiborodine/campsite-booking-go/workflows/ci/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/igor-baiborodine/campsite-booking-go)](https://goreportcard.com/report/github.com/igor-baiborodine/campsite-booking-go)
[![codecov](https://codecov.io/gh/igor-baiborodine/campsite-booking-go/graph/badge.svg?token=XTDH6MGEDJ)](https://codecov.io/gh/igor-baiborodine/campsite-booking-go)

🔥 A Java-based implementation is available in the [campsite-booking](https://github.com/igor-baiborodine/campsite-booking) repository. 

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Technical Task](#technical-task)
  - [Booking Constraints](#booking-constraints)
  - [System Requirements](#system-requirements)
- [Implementation Details](#implementation-details)
- [Project Setup](#project-setup)
- [Up and Running Locally](#up-and-running-locally)
  - [Run with IntelliJ/GoLand IDE](#run-with-intellijgoland-ide)
  - [Run with Docker Compose](#run-with-docker-compose)
  - [Run with Kubernetes](#run-with-kubernetes)
- [Tests](#tests)
  - [Unit and Integration](#unit-and-integration)
  - [Service and Method Discovery](#service-and-method-discovery)
  - [Functional and Error Handling](#functional-and-error-handling)
  - [Performance](#performance)
    - [GetCampsites](#getcampsites)
    - [GetVacantDates](#getvacantdates)
  - [Load](#load)
    - [GetCampsites](#getcampsites-1)
    - [GetVacantDates](#getvacantdates-1)

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

**Technologies used**:

* [Go](https://github.com/golang/go), [gRPC](https://github.com/grpc/grpc-go) 
* [protovalidate-go](https://github.com/bufbuild/protovalidate-go) (requests validation)
* [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports), [golines](https://github.com/segmentio/golines), [gofumpt](https://github.com/mvdan/gofumpt) (code style & formatting)
* [PostgreSQL](https://www.postgresql.org/)
* [Goose](https://pressly.github.io/goose/) (DB migrations)
* [Docker](https://www.docker.com/), [Docker Compose](https://docs.docker.com/compose/)
* [Kubernetes](https://kubernetes.io/)

The implementation of the Campsite Bookings API(or Campgrounds API) in Go is based on a
domain-centric architecture using the command and query responsibility segregation(CQRS) pattern.
It's greatly inspired by the Mallbots example application in Michael Stack's
book ["Event-Driven Architecture in Golang"](https://www.amazon.com/Event-Driven-Architecture-Golang-asynchronicity-consistency/dp/1803238011)([GitHub](https://github.com/igor-baiborodine/event-driven-architecture-in-golang-workshop)).

These resources were also used during the work on this project:
* ["Distributed Services with Go"](https://www.amazon.com/Distributed-Services-Go-Reliable-Maintainable/dp/1680507605) by Jeffrey Travis, [GitHub](https://github.com/igor-baiborodine/distributed-services-with-go-workshop)
* ["gRPC Go for Professionals"](https://www.amazon.com/gRPC-Professionals-Implement-production-grade-microservices/dp/1837638845) by Clément Jean, [GitHub](https://github.com/PacktPublishing/gRPC-Go-for-Professionals)
* ["Test-Driven Development in Go"](https://www.amazon.com/Test-Driven-Development-practical-idiomatic-real-world/dp/1803247878) by Adelina Simion, [GitHub](https://github.com/PacktPublishing/Test-Driven-Development-in-Go)

## Project Setup

**Prerequisites**:
- [Git](https://git-scm.com/), see this [guide](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) on how to install Git.
- [Make](https://man7.org/linux/man-pages/man1/make.1.html)
- [Go](https://go.dev/) (version >= 1.22), see this [guide](https://go.dev/doc/install) on how to install Go.

Clone the project and install the necessary tools(protoc, mockery, golines, goimports, gofumpt,
golangci-lint):

```shell
$ git clone https://github.com/igor-baiborodine/campsite-booking-go.git
$ cd campsite-booking-go
$ make install-tools
```

If you use either [IntelliJ IDEA](https://www.jetbrains.com/idea/) or [GoLand](https://www.jetbrains.com/go/) IDEs,
follow this [guide](/docs/ide-setup/README.md) to configure it.

>
> ⚠️ **Please note that all commands listed below should be executed from the project's root.**
>

## Up and Running Locally

### Run with IntelliJ/GoLand IDE

* Go to **Run | Edit Configurations...** and create a new `Run/Debug` configuration for the
  Campgrounds API as follows:

![Run with IDE Config](/docs/run-with-ide-config.png)

* Start a PostgreSQL DB instance using **Docker Compose**:
```shell
$ make compose-up-postgres
# which is equivalent of
$ docker compose -f docker/docker-compose.yml -p campsite-booking-go up -d postgres 
```

* Verify the health status of the running `postgres` container:
```shell
$ docker inspect --format="{{.State.Health.Status}}" postgres
```

* If the output is `healthy`, launch the `Run/Debug` configuration
  created in the previous step.

### Run with Docker Compose

> Docker images for Campgrounds API are available on [Docker Hub](https://hub.docker.com/repository/docker/ibaiborodine/campsite-booking-go/tags). 

* Start PostgreSQL DB and Campgrounds API instances using **Docker Compose**:
```shell
$ make compose-up-all
# which is equivalent of
$ docker compose -f docker/docker-compose.yml -p campsite-booking-go up -d --build 
```

### Run with Kubernetes

**Prerequisites**:

- Install [kind](https://kind.sigs.k8s.io/):
```shell
$ go install sigs.k8s.io/kind@$latest
$ kind version
```
- Install [kubectl](https://kubernetes.io/docs/reference/kubectl/):
```shell
$ curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
$ curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl.sha256"
$ echo "$(cat kubectl.sha256)  kubectl" | sha256sum --check
$ sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
$ kubectl version --client
```
---

1. Spin up a 3-node cluster: 
```bash 
$ make cluster-deploy
# which is equivalent of
$ kind create cluster --name local-k8s --config ./k8s/kind-config.yaml
$ kubectl cluster-info --context kind-local-k8s
```
2. Deploy PostgreSQL DB, Campgrounds API, and Envoy proxy:
```bash
$ make all-deploy
# which is equivalent of
# db-deploy
$ kubectl create secret generic postgres-secret --from-literal=POSTGRES_PASSWORD=postgres
$ kubectl create secret generic campgrounds-secret --from-literal=CAMPGROUNDS_PASSWORD=campgrounds_pass
$ kubectl create configmap initdb-config --from-file=./db/init/
$ kubectl apply -f ./k8s/postgres.yaml
# api-deploy
$ kubectl apply -f ./k8s/campgrounds.yaml
# proxy-deploy:
$ kubectl create configmap envoy-config --from-file=./k8s/envoy-config.yaml
$ kubectl apply -f ./k8s/envoy.yaml
```
3. Verify the status of created pods:
```bash
$ kubectl get pods 
# which may look like this
NAME                           READY   STATUS    RESTARTS      AGE
campgrounds-796fff564f-dgfsm   1/1     Running   2 (81s ago)   89s
campgrounds-796fff564f-qj44x   1/1     Running   2 (81s ago)   89s
campgrounds-796fff564f-vqfjz   1/1     Running   3 (61s ago)   89s
envoy-9dbcd5c66-h4p9v          1/1     Running   0             89s
postgres-0                     1/1     Running   0             2m13s
```
4. Use the `port-forward` command to forward Envoy’s port `8080` to `localhost:8080` to test the
   Campgrounds services using a gRPC client:
```bash
$ PROXY_POD_NAME=$(kubectl get pods --selector=app=envoy -o jsonpath='{.items[0].metadata.name}')
$ kubectl port-forward "$PROXY_POD_NAME" 8080:8080
```

## Tests

### Unit and Integration

1. Execute only unit tests:
```bash
$ make test
# which is equivalent of
$ go test -race ./internal/...
```
2. Execute only integration tests:
```bash
$ make test-integration
# which is equivalent of
$ go test -tags=integration ./internal/...
```

### Service and Method Discovery

**Prerequisites**:

- Install [gRPCurl](https://github.com/fullstorydev/grpcurl):
```shell
$ go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
$ grpcurl -version
```
- The Campgrounds API should be up & running using either the [Run with IntelliJ/GoLand IDE](#run-with-intellijgoland-ide) or [Run with Docker Compose](#run-with-docker-compose).
---

1. List the services present on the gRPC server:
```bash
$ grpcurl -plaintext localhost:8085 list
# output
campgroundspb.v1.CampgroundsService
grpc.reflection.v1.ServerReflection
grpc.reflection.v1alpha.ServerReflection
```
2. List all the RPC endpoints the `campgroundspb.v1.CampgroundsService` contains:
```bash
$ grpcurl -plaintext localhost:8085 describe campgroundspb.v1.CampgroundsService
# output
campgroundspb.v1.CampgroundsService is a service:
service CampgroundsService {
  rpc CancelBooking ( .campgroundspb.v1.CancelBookingRequest ) returns ( .campgroundspb.v1.CancelBookingResponse );
  rpc CreateBooking ( .campgroundspb.v1.CreateBookingRequest ) returns ( .campgroundspb.v1.CreateBookingResponse );
  rpc CreateCampsite ( .campgroundspb.v1.CreateCampsiteRequest ) returns ( .campgroundspb.v1.CreateCampsiteResponse );
  rpc GetBooking ( .campgroundspb.v1.GetBookingRequest ) returns ( .campgroundspb.v1.GetBookingResponse );
  rpc GetCampsites ( .campgroundspb.v1.GetCampsitesRequest ) returns ( .campgroundspb.v1.GetCampsitesResponse );
  rpc GetVacantDates ( .campgroundspb.v1.GetVacantDatesRequest ) returns ( .campgroundspb.v1.GetVacantDatesResponse );
  rpc UpdateBooking ( .campgroundspb.v1.UpdateBookingRequest ) returns ( .campgroundspb.v1.UpdateBookingResponse );
}
```
3. Get a gRPC message definition, for example for `campgroundspb.v1.GetBookingRequest`:
```bash
grpcurl -plaintext localhost:8085 describe campgroundspb.v1.GetBookingRequest
# output
campgroundspb.v1.GetBookingRequest is a message:
message GetBookingRequest {
  string booking_id = 1 [(.buf.validate.field) = { string: { uuid: true } }];
}
```

### Functional and Error Handling

**Prerequisites**:
- The same as for the Service and Method Discovery tests.
---

1. Create a campsite:
```bash
$ grpcurl -plaintext -d \
    '{"campsite_code": "CAMP01", "capacity": 4, "drinking_water": true, "fire_pit": true, "picnic_table": true, "restrooms": false}' \
    localhost:8085 campgroundspb.v1.CampgroundsService/CreateCampsite
# output
{
  "campsiteId": "07df7f35-9c7a-4b10-a702-66844a7ec08c"
}
```
2. Get campsites:
```bash
$ grpcurl -plaintext -d '{}' localhost:8085 campgroundspb.v1.CampgroundsService/GetCampsites
# output
{
  "campsites": [
    {
      "campsiteId": "07df7f35-9c7a-4b10-a702-66844a7ec08c",
      "campsiteCode": "CAMP01",
      "capacity": 4,
      "drinkingWater": true,
      "picnicTable": true,
      "firePit": true,
      "active": true
    }
  ]
}
```
3. Create a booking:
```bash
$ grpcurl -plaintext -d \
    '{"campsite_id": "07df7f35-9c7a-4b10-a702-66844a7ec08c", "email": "john.smith@example.com", "full_name": "John Smith", "start_date": "2024-09-09", "end_date": "2024-09-12"}' \
    localhost:8085 campgroundspb.v1.CampgroundsService/CreateBooking
# output
{
  "bookingId": "692abbc0-5457-4f2b-8a6e-061ba2e5dd90"
}
```
4. Get a booking:
```bash
$ grpcurl -plaintext -d \
    '{"booking_id": "692abbc0-5457-4f2b-8a6e-061ba2e5dd90"}' \
    localhost:8085 campgroundspb.v1.CampgroundsService/GetBooking
# output
{
  "booking": {
    "bookingId": "692abbc0-5457-4f2b-8a6e-061ba2e5dd90",
    "campsiteId": "07df7f35-9c7a-4b10-a702-66844a7ec08c",
    "email": "john.smith@example.com",
    "fullName": "John Smith",
    "startDate": "2024-09-09",
    "endDate": "2024-09-12",
    "active": true
  }
}
```
5. Create a booking that does not meet the [booking constraints](#booking-constraints), for example a
   maximum stay of three days:
```bash
$ grpcurl -plaintext -d \
    '{"campsite_id": "07df7f35-9c7a-4b10-a702-66844a7ec08c", "email": "john.smith@example.com", "full_name": "John Smith", "start_date": "2024-09-15", "end_date": "2024-09-20"}' \
    localhost:8085 campgroundspb.v1.CampgroundsService/CreateBooking
# output
ERROR:
  Code: InvalidArgument
  Message: booking validation: 1 error occurred:
        * maximum stay: must be less or equal to three days
```
6. Create booking for non-existing campsite ID:
```bash
$ grpcurl -plaintext -d \
  '{"campsite_id": "a2432518-0fc0-496f-8f78-ac9902a44e3d", "start_date": "2024-11-21", "end_date": "2024-11-23", "email": "john.smith.1@email.com", "full_name": "John Smith 1"}' \
  localhost:8085 campgroundspb.v1.CampgroundsService/CreateBooking 
# output
ERROR:
  Code: Internal
  Message: insert booking: ERROR: insert or update on table "bookings" violates foreign key constraint "fk_bookings_campsite_id_campsites" (SQLSTATE 23503)
  Details:
  1)    {
          "@type": "type.googleapis.com/errors.ErrorType",
          "GRPCCode": "13",
          "HTTPCode": "500",
          "TypeCode": "INTERNAL_SERVER_ERROR"
        }
```

### Performance

**Prerequisites**:
- The Campgrounds API should be up & running using the [Run with Docker Compose](#run-with-docker-compose).
- The `pprof` tool should reachable at http://localhost:6060/debug/pprof/ in a browser of your choice.
- Run the data generator to create, for example, 100 campsites and non-consecutive bookings for each
  campsite:
```bash
$ go run ./datagenerator/main.go localhost:8085 100
# output
igor@lptacr:~/GitRepos/igor-baiborodine/campsite-booking-go$ go run ./datagenerator/main.go localhost:8085 100
2024/09/22 19:03:06 server address: localhost:8085, campsites count: 100
2024/09/22 19:03:06 created 100 campsites
2024/09/22 19:03:06 ...created 10 bookings for campsite ID bdf7e4fb-4d35-49aa-aca7-2876c4e25135
2024/09/22 19:03:06 ...created 10 bookings for campsite ID 408de8b6-b552-4905-b1c3-c54e038bbfaf
... more created bookings output
2024/09/22 19:03:09 ...created 9 bookings for campsite ID f954f06b-e5c8-4b04-8f68-1b1e722ce0fb
2024/09/22 19:03:10 ...created 9 bookings for campsite ID aada6ebf-9c5c-46fe-ad0e-f4a27741085f
2024/09/22 19:03:10 created total 946 bookings
```
---

#### GetCampsites

1. Start downloading the profiling data for the `GetCampsites` endpoint from the past 10 seconds and
   save it to a local file named `get-campsites-profile.pprof`. Then immediately execute the
   corresponding benchmark test:
```bash
$ make pprof-get-campsites
# which is equivalent of
$ curl --output ./tests/perf/get-campsites-profile.pprof "http://localhost:6060/debug/pprof/profile?seconds=10"
$ SERVER_ADDR=localhost:8085 go test -bench BenchmarkGetCampsites ./tests/perf
```
2. Validate the profiling data for the `GetCampsites` endpoint by launching the `pprof` tool. When
   prompted, enter the `web` option to generate a report in `SVG` format on a temp file, and start a
   web browser to view it. Alternatively, you can use the `png` option, to generate a report in `PNG`
   format:
```bash
$ make pprof-get-campsites-data
# which is equivalent of
go tool pprof ./tests/perf/get-campsites-profile.pprof
# output
File: app
Type: cpu
Time: Sep 21, 2024 at 5:52pm (EDT)
Duration: 10.01s, Total samples = 1.20s (11.99%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) web
(pprof) png
Generating report in profile001.png
(pprof) 
```
3. See https://git.io/JfYMW on how to read the graph:
   ![pprof GetCampsites data](/docs/pprof-get-campsites-data.png)

#### GetVacantDates

1. Start downloading the profiling data for the `GetVacantDates` endpoint from the past 10 seconds and
   save it to a local file named `get-vacant-dates-profile.pprof`. Then immediately execute the
   corresponding benchmark test:
```bash
$ make pprof-get-vacant-dates
# which is equivalent of
$ curl --output ./tests/perf/get-vacant-dates-profile.pprof "http://localhost:6060/debug/pprof/profile?seconds=10"
$ SERVER_ADDR=localhost:8085 go test -bench BenchmarkGetVacantDates ./tests/perf
```
2. Validate the profiling data for the `GetVacantDates` endpoint by launching the `pprof` tool. When
   prompted, use either the `web` or `png` option to generate a corresponding report.
```bash
$ make pprof-get-vacant-dates-data
# which is equivalent of
go tool pprof ./tests/perf/get-vacant-dates-profile.pprof
```

### Load

**Prerequisites**:
- Install [ghz](https://ghz.sh/):
```bash
$ go install github.com/bojand/ghz/cmd/ghz@latest
```
- The Campgrounds API should be up & running using the [Run with Docker Compose](#run-with-docker-compose).
- Run the data generator to create, for example, 100 campsites and non-consecutive bookings for each
  campsite:
```bash
$ go run ./datagenerator/main.go localhost:8085 100
```
- When using the `buf generate` command, `buf` fetches the dependencies and uses them to generate
  the necessary files. These dependencies are not stored on your local file system in a directly
  accessible way. Instead, `buf` manages these dependencies in a non-visible, internal cache.
  Therefore, execute the following command to load and save the `protovalidate` dependency:
```bash
$ buf export buf.build/bufbuild/protovalidate --output ./campgroundspb/v1/
```

#### GetCampsites

Execute the following command to perform a basic load testing of the `GetCampsites` endpoint:
```bash
$ ghz --insecure --proto ./campgroundspb/v1/api.proto \
  --import-paths ./campgroundspb/buf/validate/validate.proto \
  --call campgroundspb.v1.CampgroundsService/GetCampsites \
  -n 10000 -c 10 -d '{}' localhost:8085
```
Where:
* `-n 10000` - number of requests to run
* `-c 10` - number of request workers to run concurrently

The output may look like the one below:
```text
Summary:
  Count:        10000
  Total:        21.53 s
  Slowest:      100.66 ms
  Fastest:      2.70 ms
  Average:      21.25 ms
  Requests/sec: 464.49

Response time histogram:
  2.704   [1]    |
  12.499  [5177] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  22.295  [1778] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  32.090  [268]  |∎∎
  41.885  [859]  |∎∎∎∎∎∎∎
  51.680  [1081] |∎∎∎∎∎∎∎∎
  61.476  [563]  |∎∎∎∎
  71.271  [184]  |∎
  81.066  [70]   |∎
  90.861  [16]   |
  100.657 [3]    |

Latency distribution:
  10 % in 6.47 ms 
  25 % in 8.44 ms 
  50 % in 12.17 ms 
  75 % in 36.18 ms 
  90 % in 50.05 ms 
  95 % in 56.51 ms 
  99 % in 69.96 ms 

Status code distribution:
  [OK]   10000 responses   
```

#### GetVacantDates

Execute the following command to perform a basic load testing of the `GetVacantDates` endpoint:
```bash
$ ghz --insecure --proto ./campgroundspb/v1/api.proto \
  --import-paths ./campgroundspb/buf/validate/validate.proto \
  --call campgroundspb.v1.CampgroundsService/GetVacantDates \
  -n 10000 -c 10 \
  -d '{"campsite_id":"167ce4b6-8616-4757-9de0-3bbed703d51a","start_date":"2024-09-23","end_date":"2024-10-23"}' localhost:8085
```
Where:
* `-n 10000` - number of requests to run
* `-c 10` - number of request workers to run concurrently

The output may look like the one below:
```text
Summary:
  Count:        10000
  Total:        14.12 s
  Slowest:      69.36 ms
  Fastest:      1.22 ms
  Average:      13.90 ms
  Requests/sec: 708.25

Response time histogram:
  1.224  [1]    |
  8.038  [6845] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  14.852 [615]  |∎∎∎∎
  21.666 [8]    |
  28.479 [77]   |
  35.293 [616]  |∎∎∎∎
  42.107 [969]  |∎∎∎∎∎∎
  48.921 [659]  |∎∎∎∎
  55.735 [178]  |∎
  62.549 [27]   |
  69.363 [5]    |

Latency distribution:
  10 % in 3.14 ms 
  25 % in 4.11 ms 
  50 % in 5.74 ms 
  75 % in 26.75 ms 
  90 % in 41.12 ms 
  95 % in 45.23 ms 
  99 % in 51.71 ms 

Status code distribution:
  [OK]   10000 responses 
```
