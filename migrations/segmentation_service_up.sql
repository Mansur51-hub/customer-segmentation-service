create table segments
(
    id   serial primary key,
    slug varchar(255) not null unique
);

create table users
(
    id int primary key
);

create table operations
(
    id             serial primary key,
    user_id        int          not null,
    segment_slug   varchar(255) not null,
    created_at     timestamp    not null default now(),
    operation_type varchar(255) not null,
    foreign key (user_id) references users (id)
);

create table memberships
(
    id           serial primary key,
    user_id      int          not null,
    segment_slug varchar(255) not null,
    created_at   timestamp    not null default now(),
    expired_at   timestamp             default null,
    foreign key (segment_slug) references segments (slug),
    UNIQUE (user_id, segment_slug)
);
