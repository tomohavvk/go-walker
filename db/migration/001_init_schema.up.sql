CREATE EXTENSION IF NOT EXISTS timescaledb;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

create table if not exists devices
(
    id         varchar(64) not null primary key,
    status     varchar(7) not null,
    created_at timestamp   not null,
    updated_at timestamp   not null
);

create index if not exists device_id_status_idx on devices (id, status);


create table if not exists device_locations
(
    device_id         varchar(64) not null references devices (id) on delete cascade,
    latitude          int         not null,
    longitude         int         not null,
    accuracy          int         not null,
    altitude          int         not null,
    speed             int         not null,
    bearing           int         not null,
    altitude_accuracy int         not null,
    time              timestamp   not null
);

SELECT create_hypertable('device_locations', 'time');
SELECT set_chunk_time_interval('device_locations', interval '1 month');

create unique index if not exists device_locations_device_id_time_idx on device_locations (device_id, time);

create table if not exists groups
(
    id              varchar(21) PRIMARY KEY,
    owner_device_id varchar(64)  not null references devices (id) on delete cascade,
    name            varchar(128) not null,
    is_public       boolean      not null,
    public_id       varchar(32),
    description     varchar(256),
    created_at      timestamp    not null,
    updated_at      timestamp    not null
);

create unique index if not exists groups_public_id_uniq_idx on groups (public_id);
create index if not exists groups_public_id_idx on groups USING gin (public_id gin_trgm_ops);
create index if not exists groups_name_trgm_idx on groups USING gin (name gin_trgm_ops);
create index if not exists groups_updated_at_idx on groups (updated_at desc);

create table if not exists devices_groups
(
    device_id  varchar(64) not null references devices (id) on delete cascade,
    group_id   varchar(64) not null references groups (id) on delete cascade,
    created_at timestamp   not null
);

create unique index if not exists devices_groups_device_id_group_id_idx on devices_groups (device_id, group_id);


create table if not exists group_messages
(
    group_id         varchar(64)  not null references groups (id) on delete cascade,
    author_device_id varchar(64)  not null references devices (id) on delete cascade,
    message          varchar(512) not null,
    created_at        timestamp    not null
);

SELECT create_hypertable('group_messages', 'created_at');
SELECT set_chunk_time_interval('group_messages', interval '1 month');
