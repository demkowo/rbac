package postgres

import (
	"database/sql"
	"errors"
	"log"

	model "github.com/demkowo/rbac/models"
)

const (
	ADD_RBAC    = "INSERT INTO rbac (route_id, role_id) VALUES ($1, $2) ON CONFLICT (route_id, role_id) DO NOTHING;"
	DELETE_RBAC = "DELETE FROM rbac WHERE route_id = $1 AND role_id = $2;"
	FIND_RBAC   = "SELECT route_id, role_id FROM rbac;"
)

type Rbac interface {
	Add(*model.Rbac) error
	Delete(*model.Rbac) error
	Find() ([]*model.Rbac, error)
}

type rbac struct {
	db *sql.DB
}

func NewRbac(db *sql.DB) Rbac {
	return &rbac{db: db}
}

func (r *rbac) Add(rbac *model.Rbac) error {
	_, err := r.db.Exec(ADD_RBAC, rbac.RouteID, rbac.RoleID)
	if err != nil {
		log.Printf("failed to execute db.Exec AUTH_ADD: %v", err)
		return errors.New("failed to add rbac record")
	}
	return nil
}

func (r *rbac) Delete(rbac *model.Rbac) error {
	_, err := r.db.Exec(DELETE_RBAC, rbac.RouteID, rbac.RoleID)
	if err != nil {
		log.Printf("failed to execute db.Exec DELETE_RBAC: %v", err)
		return errors.New("failed to delete rbac record")
	}

	return nil
}

func (r *rbac) Find() ([]*model.Rbac, error) {
	rows, err := r.db.Query(FIND_RBAC)
	if err != nil {
		log.Printf("failed to execute db.Query FIND_RBAC failed: %v", err)
		return nil, errors.New("failed to find rbac records")
	}
	defer rows.Close()

	var rbacs []*model.Rbac
	for rows.Next() {
		var rbac model.Rbac
		if err := rows.Scan(&rbac.RouteID, &rbac.RoleID); err != nil {
			log.Printf("failed to scan FIND_RBAC rows: %v", err)
			return nil, errors.New("failed to find rbac records")
		}
		rbacs = append(rbacs, &rbac)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error while iterating over rbacs: %v", err)
		return nil, errors.New("failed to find rbac records")
	}

	log.Printf("Retrieved %d records", len(rbacs))
	return rbacs, nil
}
