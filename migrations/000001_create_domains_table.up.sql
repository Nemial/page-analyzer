CREATE TABLE IF NOT EXISTS domains
(
    id         int primary key auto_increment,
    name       varchar(255),
    updated_at timestamp,
    created_at timestamp
);