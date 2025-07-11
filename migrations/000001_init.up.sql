create table users (
    id serial primary key,
    login varchar(255) not null unique,
    pwd_hash varchar(255) not null,
    created_at timestamp with time zone default current_timestamp
);