CREATE TABLE IF NOT EXISTS service_user
(
    id       INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name     varchar,
    secret   varchar
);