CREATE TABLE IF NOT EXISTS domain_checks
(
    id          int primary key auto_increment,
    domain_id   int,
    FOREIGN KEY (domain_id) REFERENCES domains (id),
    status_code int,
    h1          text,
    keywords    text,
    description text,
    updated_at  timestamp,
    created_at  timestamp

)