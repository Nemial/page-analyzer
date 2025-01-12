package mysql

import (
	"github.com/jmoiron/sqlx"
	"page-analyzer/internal/models"
)

type DomainCheckModel struct {
	DB *sqlx.DB
}

func (m *DomainCheckModel) Insert(domainId int, statusCode int, h1 string, keywords string, description string) (int64, error) {
	result, err := m.DB.Exec("INSERT INTO domain_checks(domain_id, status_code, h1, keywords, description, updated_at, created_at) VALUES(?, ?, ?, ?, ?, UTC_TIMESTAMP, UTC_TIMESTAMP)", domainId, statusCode, h1, keywords, description)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *DomainCheckModel) GetByDomain(domainId int) (*[]models.DomainCheck, error) {
	var domainChecks []models.DomainCheck
	err := m.DB.Select(&domainChecks, "SELECT * FROM domain_checks WHERE domain_id = ?", domainId)
	if err != nil {
		return &[]models.DomainCheck{}, err
	}

	return &domainChecks, nil
}

func (m *DomainCheckModel) GetAll() (*map[int]models.DomainCheck, error) {
	var domainCheck []models.DomainCheck

	err := m.DB.Select(&domainCheck, "SELECT DISTINCT domain_id, created_at, status_code FROM domain_checks  ORDER BY created_at DESC")
	if err != nil {
		return &map[int]models.DomainCheck{}, err
	}

	result := make(map[int]models.DomainCheck)

	for _, v := range domainCheck {
		if _, ok := result[v.DomainId]; !ok {
			result[v.DomainId] = v
		}
	}
	return &result, nil
}
