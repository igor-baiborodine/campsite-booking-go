-- +goose Up
CREATE TABLE campgrounds.campsites
(
    id             bigint GENERATED BY DEFAULT AS IDENTITY NOT NULL,
    campsite_id    varchar(255)                            NOT NULL,
    campsite_code  varchar(255)                            NOT NULL,
    capacity       int                                     NOT NULL,
    restrooms      boolean                                 NOT NULL,
    drinking_water boolean                                 NOT NULL,
    picnic_table   boolean                                 NOT NULL,
    fire_pit       boolean                                 NOT NULL,
    active         boolean                                 NOT NULL,
    created_at     timestamptz                             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     timestamptz                             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_campsites PRIMARY KEY (id)
);

CREATE UNIQUE INDEX unique_campsites_campsite_id ON campgrounds.campsites (campsite_id);
CREATE UNIQUE INDEX unique_campsites_campsite_code ON campgrounds.campsites (campsite_code);

-- +goose Down
DROP TABLE IF EXISTS campgrounds.campsites;