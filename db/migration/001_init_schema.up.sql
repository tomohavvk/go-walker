CREATE EXTENSION IF NOT EXISTS timescaledb;

create table if not exists devices
(
    id         varchar(64) not null PRIMARY KEY,
    name       varchar(64) not null,
    created_at timestamp   not null
);

create table if not exists devices_locations
(
    device_id         varchar(64) NOT NULL REFERENCES devices (id) ON DELETE CASCADE,
    latitude          numeric     not null,
    longitude         numeric     not null,
    accuracy          numeric     not null,
    altitude          numeric     not null,
    speed             numeric     not null,
    bearing           numeric     not null,
    altitude_accuracy numeric     not null,
    time              timestamp   not null
);

SELECT create_hypertable('devices_locations', 'time');
SELECT set_chunk_time_interval('devices_locations', INTERVAL '1 month');

create unique index if not exists devices_locations_device_id_time_idx on devices_locations (device_id, time);
