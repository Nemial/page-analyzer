package mysql

import (
	"github.com/jmoiron/sqlx"
	"page-analyzer/internal/models"
)

type DomainModel struct {
	DB *sqlx.DB
}

func (m *DomainModel) Insert(name string) (int64, error) {
	result, err := m.DB.Exec("INSERT INTO domains(name, updated_at, created_at) VALUES (?, UTC_TIMESTAMP, UTC_TIMESTAMP)", name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *DomainModel) Get(id int) (*models.Domain, error) {
	var domain models.Domain

	err := m.DB.Get(&domain, "SELECT * FROM domains WHERE id = ?", id)
	if err != nil {
		return &models.Domain{}, err
	}

	return &domain, nil
}

func (m *DomainModel) GetAll() (*[]models.Domain, error) {
	var domains []models.Domain

	err := m.DB.Select(&domains, "SELECT * FROM domains")
	if err != nil {
		return &[]models.Domain{}, err
	}

	return &domains, nil
}
