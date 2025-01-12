package models

import "time"

type Domain struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}

type DomainCheck struct {
	Id          int       `db:"id"`
	DomainId    int       `db:"domain_id"`
	StatusCode  int       `db:"status_code"`
	H1          string    `db:"h1"`
	Keywords    string    `db:"keywords"`
	Description string    `db:"description"`
	UpdatedAt   time.Time `db:"updated_at"`
	CreatedAt   time.Time `db:"created_at"`
}
