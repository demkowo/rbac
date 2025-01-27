package postgres

import (
	"database/sql"
	"errors"
	"log"

	model "github.com/demkowo/rbac/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	ADD_ROLE               = "INSERT INTO roles (id, name) VALUES ($1, $2) ON CONFLICT (name) DO NOTHING;"
	DELETE_ROLE            = "DELETE FROM roles WHERE id = $1;"
	FIND_ROLES             = "SELECT id, name FROM roles ORDER BY name;"
	FIND_ROLES_BY_ROUTE_ID = `
        SELECT roles.id, roles.name
        FROM roles
        INNER JOIN rbac ON roles.id = rbac.role_id
        WHERE rbac.route_id = $1
    `
	ROLE_EXISTS_BY_ID = "SELECT EXISTS(SELECT 1 FROM roles WHERE id=$1)"
	UPDATE_ROLE       = "UPDATE roles SET name = $2 WHERE id = $1;"
)

type Roles interface {
	Add(*model.Role) error
	Delete(string) error
	ExistsByID(uuid.UUID) (bool, error)
	Find() ([]*model.Role, error)
	FindByRoute(uuid.UUID) ([]*model.Role, error)
	Update(*model.Role) error
}

type roles struct {
	db *sql.DB
}

func NewRoles(db *sql.DB) Roles {
	return &roles{db: db}
}

func (r *roles) Add(role *model.Role) error {
	_, err := r.db.Exec(ADD_ROLE, role.ID, role.Name)
	if err != nil {
		log.Printf("failed to execute db.Exec ADD_ROLE: %v", err)
		return errors.New("failed to add role")
	}
	return nil
}

func (r *roles) Delete(roleId string) error {
	_, err := r.db.Exec(DELETE_ROLE, roleId)
	if err != nil {
		log.Printf("failed to execute db.Exec DELETE_ROLE: %v", err)
		return errors.New("failed to delete role")
	}
	return nil
}

func (r *roles) ExistsByID(id uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ROLE_EXISTS_BY_ID, id).Scan(&exists)
	if err != nil {
		log.Printf("failed to execute db.QueryRow ROLE_EXISTS_BY_ID: %v", err)
		return false, errors.New("failed to check if role exists")
	}
	return exists, nil
}

func (r *roles) Find() ([]*model.Role, error) {
	rows, err := r.db.Query(FIND_ROLES)
	if err != nil {
		log.Printf("failed to execute db.Query FIND_ROLES: %v", err)
		return nil, errors.New("failed to find roles")
	}
	defer rows.Close()

	var roles []*model.Role
	for rows.Next() {
		var role model.Role
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			log.Printf("failed to scan FIND_ROLES record: %v", err)
			return nil, errors.New("failed to find roles")
		}
		roles = append(roles, &role)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error while iterating over roles: %v", err)
		return nil, errors.New("failed to find roles")
	}

	log.Printf("Retrieved %d roles", len(roles))
	return roles, nil
}

func (r *roles) FindByRoute(routeID uuid.UUID) ([]*model.Role, error) {
	rows, err := r.db.Query(FIND_ROLES_BY_ROUTE_ID, routeID)
	if err != nil {
		log.Printf("failed to execute db.Query FIND_ROLES_BY_ROUTE_ID: %v", err)
		return nil, errors.New("failed to find roles for route")
	}
	defer rows.Close()

	var roles []*model.Role
	for rows.Next() {
		var role model.Role
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			log.Printf("failed to scan FIND_ROLES_BY_ROUTE_ID record: %v", err)
			return nil, errors.New("failed to find roles for route")
		}
		roles = append(roles, &role)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error while iterating over roles: %v", err)
		return nil, errors.New("failed to find roles for route")
	}

	return roles, nil
}

func (r *roles) Update(role *model.Role) error {
	_, err := r.db.Exec(UPDATE_ROLE, role.ID, role.Name)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				log.Printf("duplicate key error on UPDATE_ROLE: %v", pqErr.Detail)
				return errors.New("role with the given name already exists")
			}
		}
		log.Printf("failed to execute db.Exec UPDATE_ROLE: %v", err)
		return errors.New("failed to update role")
	}

	return nil
}
