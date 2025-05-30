################################################################################
# Stage 1: cache modules
################################################################################
FROM golang:1.24.2-bullseye AS modules

COPY go.mod go.sum /usr/src/modules/
WORKDIR /usr/src/modules
RUN go mod download

################################################################################
# Stage 2: build executable
################################################################################
FROM golang:1.24.2-bullseye AS builder

COPY --from=modules /go/pkg /go/pkg
COPY . /usr/src/app
WORKDIR /usr/src/app

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -v -o /bin/app ./cmd

################################################################################
# Stage 3: build app image
################################################################################
FROM debian:bullseye-slim

LABEL \
    maintainer="Igor Baiborodine <igor@kiroule.com>" \
    org.label-schema.schema-version="1.0" \
    org.label-schema.name="campgrounds-app" \
    org.label-schema.vcs-url="https://github.com/igor-baiborodine/campsite-booking-go" \
    org.label-schema.usage="https://github.com/igor-baiborodine/campsite-booking-go/blob/main/readme/README.md"

ENV APP_NAME=campgrounds
ENV APP_HOME=/opt/${APP_NAME}
ENV APP_USER=${APP_NAME}
ENV APP_GROUP=${APP_NAME}
ENV PATH="${APP_HOME}/bin:${PATH}"

RUN groupadd ${APP_GROUP}; \
    useradd -g ${APP_GROUP} ${APP_USER}

RUN set -ex; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
        # su tool for easy step-down from root
        gosu; \
    rm -rf /var/lib/apt/lists/*; \
    gosu nobody true

COPY --from=builder /bin/app ${APP_HOME}/bin/
COPY docker/docker-entrypoint.sh ${APP_HOME}/bin/

RUN chmod a+x ${APP_HOME}/bin/*

WORKDIR ${APP_HOME}
ENTRYPOINT ["docker-entrypoint.sh"]
EXPOSE 8085 6060
CMD ["app"]
