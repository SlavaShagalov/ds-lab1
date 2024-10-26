create table if not exists persons
(
    id      bigserial primary key,
    name    text not null,
    age     int  not null,
    address text not null,
    work    text not null
);
