# Campsite Bookings API (Go)

![ci](https://github.com/igor-baiborodine/campsite-booking-go/workflows/ci/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/igor-baiborodine/campsite-booking-go)](https://goreportcard.com/report/github.com/igor-baiborodine/campsite-booking-go)
[![codecov](https://codecov.io/gh/igor-baiborodine/campsite-booking-go/graph/badge.svg?token=XTDH6MGEDJ)](https://codecov.io/gh/igor-baiborodine/campsite-booking-go)


<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Technical Task](#technical-task)
  - [Booking Constraints](#booking-constraints)
  - [System Requirements](#system-requirements)
- [Implementation Details](#implementation-details)
- [Project Setup](#project-setup)
- [Up & Running Locally](#up--running-locally)
  - [Run with IntelliJ/GoLand IDE](#run-with-intellijgoland-ide)
  - [Run with Docker Compose](#run-with-docker-compose)
  - [Run with Kubernetes](#run-with-kubernetes)
- [Tests](#tests)
  - [Unit & Integration](#unit--integration)
  - [gRPCurl](#grpcurl)

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
* [Goose](https://pressly.github.io/goose/) (DB migrations)
* [PostgreSQL](https://www.postgresql.org/)  
* [Docker](https://www.docker.com/), [Docker Compose](https://docs.docker.com/compose/)
* TODO

TODO: elaborate on implementation details

## Project Setup

**Prerequisites**:
1. [Git](https://git-scm.com/), see this [guide](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) on how to install Git.
2. [Make](https://man7.org/linux/man-pages/man1/make.1.html)
3. [Go](https://go.dev/) (version >= 1.22), see this [guide](https://go.dev/doc/install) on how to install Go.

Clone the project and install the necessary tools(protoc, mockery, golines, goimports, gofumpt,
golangci-lint):

```shell
$ git clone https://github.com/igor-baiborodine/campsite-booking-go.git
$ cd campsite-booking-go
$ make install-tools
```

If you use either [IntelliJ IDEA](https://www.jetbrains.com/idea/) or [GoLand](https://www.jetbrains.com/go/) IDEs,
follow this [guide](/docs/ide-setup/README.md) to configure it.

## Up & Running Locally

⚠️ Please note that all commands listed below should be executed from the project's root.

### Run with IntelliJ/GoLand IDE

* Go to **Run | Edit Configurations...** and create a new `Run/Debug` configuration as follows:

![Run with IDE Config](/docs/run-with-ide-config.png)

* Start a PostgreSQL instance using **Docker Compose**:
```shell
$ docker compose -f docker/docker-compose.yml -p campsite-booking-go up -d postgres 
```

* Verify the health status of the running `postgres` container:
```shell
$ docker inspect --format="{{.State.Health.Status}}" postgres
```

* If the output is `healthy`, launch the `Run/Debug` configuration
  created in the previous step.

### Run with Docker Compose

* Start PostgreSQL and campgrounds app instances using **Docker Compose**:
```shell
$ docker compose -f docker/docker-compose.yml -p campsite-booking-go up -d 
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

1. Spin up a 3-node cluster: 
```bash 
$ make cluster-deploy
# which is equivalent of
$ kind create cluster --name local-k8s --config ./k8s/kind-config.yaml
$ kubectl cluster-info --context kind-local-k8s
```
2. Deploy PostgreSQL, Campgrounds application, and Envoy proxy:
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
   Campgrounds services:
```bash
$ PROXY_POD_NAME=$(kubectl get pods --selector=app=envoy -o jsonpath='{.items[0].metadata.name}')
$ kubectl port-forward "$PROXY_POD_NAME" 8080:8080
```

## Tests

### Unit & Integration

TODO

### gRPCurl

TODO