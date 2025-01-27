package postgres

import (
	"database/sql"
	"errors"
	"log"

	model "github.com/demkowo/rbac/models"
	"github.com/google/uuid"
)

const (
	ADD_ACTIVE_ROUTES      = "INSERT INTO routes (id, method, path, active) VALUES ($1,$2,$3,$4) ON CONFLICT (method, path) DO UPDATE SET active=$4"
	ADD_ROUTE              = `INSERT INTO routes (id, method, path, active) VALUES ($1, $2, $3, $4)ON CONFLICT (method, path) DO UPDATE SET active = EXCLUDED.active`
	DELETE_ROUTE           = "DELETE FROM routes WHERE id = $1"
	ROUTE_EXISTS_BY_ID     = "SELECT EXISTS(SELECT 1 FROM routes WHERE id=$1)"
	FIND_ROUTES            = "SELECT id, method, path, active FROM routes ORDER BY path, method"
	FIND_ROUTES_BY_ROLE_ID = `
        SELECT routes.id, routes.method, routes.path, routes.active
        FROM routes
        INNER JOIN rbac ON routes.id = rbac.route_id
        WHERE rbac.role_id = $1
    `
	SET_ROUTE_INACTIVE = "UPDATE routes SET active = false"
	UPDATE_ROUTE       = "UPDATE routes SET method = $2, path = $3, active = $4 WHERE id = $1"
)

type Routes interface {
	AddActive([]*model.Route) error
	Add(*model.Route) error
	Delete(uuid.UUID) error
	ExistsByID(uuid.UUID) (bool, error)
	Find() ([]*model.Route, error)
	FindByRole(uuid.UUID) ([]*model.Route, error)
	SetInactive() error
	Update(*model.Route) error
}

type routes struct {
	db *sql.DB
}

func NewRoutes(db *sql.DB) Routes {
	return &routes{db: db}
}

func (r *routes) AddActive(routes []*model.Route) error {
	for _, route := range routes {
		_, err := r.db.Exec(ADD_ACTIVE_ROUTES, route.ID, route.Method, route.Path, route.Active)
		if err != nil {
			log.Printf("failed to execute db.Exec ADD_ACTIVE_ROUTES: %v", err)
			return errors.New("failed to update list of routes")
		}
	}
	return nil
}

func (r *routes) Add(route *model.Route) error {
	_, err := r.db.Exec(ADD_ROUTE, route.ID, route.Method, route.Path, route.Active)
	if err != nil {
		log.Printf("failed to execute db.Exec ADD_ROUTE: %v", err)
		return errors.New("failed to add route")
	}
	return nil
}

func (r *routes) Delete(id uuid.UUID) error {
	_, err := r.db.Exec(DELETE_ROUTE, id)
	if err != nil {
		log.Printf("failed to execute db.Exec DELETE_ROUTE: %v", err)
		return errors.New("failed to delete route")
	}

	return nil
}

func (r *routes) ExistsByID(id uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ROUTE_EXISTS_BY_ID, id).Scan(&exists)
	if err != nil {
		log.Printf("failed to execute db.QueryRow EXISTS_BY_ID: %v", err)
		return false, errors.New("failed to check if route exists")
	}
	return exists, nil
}

func (r *routes) Find() ([]*model.Route, error) {
	rows, err := r.db.Query(FIND_ROUTES)
	if err != nil {
		log.Printf("failed to execute db.Query FIND_ROUTES failed: %v", err)
		return nil, errors.New("failed to fetch routes")
	}
	defer rows.Close()

	var routes []*model.Route
	for rows.Next() {
		var route model.Route
		if err := rows.Scan(&route.ID, &route.Method, &route.Path, &route.Active); err != nil {
			log.Printf("failed to scan FIND_ROUTES record: %v", err)
			return nil, errors.New("failed to fetch routes")
		}
		routes = append(routes, &route)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error while iterating over routes: %v", err)
		return nil, errors.New("failed to fetch routes")
	}

	log.Printf("Retrieved %d routes", len(routes))
	return routes, nil
}

func (r *routes) FindByRole(roleID uuid.UUID) ([]*model.Route, error) {
	rows, err := r.db.Query(FIND_ROUTES_BY_ROLE_ID, roleID)
	if err != nil {
		log.Printf("failed to execute db.Query FIND_ROUTES_BY_ROLE_ID: %v", err)
		return nil, errors.New("failed to find routes for role")
	}
	defer rows.Close()

	var routes []*model.Route
	for rows.Next() {
		var route model.Route
		if err := rows.Scan(&route.ID, &route.Method, &route.Path, &route.Active); err != nil {
			log.Printf("failed to scan FIND_ROUTES_BY_ROLE_ID record: %v", err)
			return nil, errors.New("failed to find routes for role")
		}
		routes = append(routes, &route)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error while iterating over routes: %v", err)
		return nil, errors.New("failed to find routes for role")
	}

	return routes, nil
}

func (r *routes) SetInactive() error {
	_, err := r.db.Exec(SET_ROUTE_INACTIVE)
	if err != nil {
		log.Printf("failed to execute db.Exec SET_ROUTE_INACTIVE: %v", err)
		return errors.New("failed to set routes inactive")
	}

	return nil
}

func (r *routes) Update(route *model.Route) error {
	_, err := r.db.Exec(UPDATE_ROUTE, route.ID, route.Method, route.Path, route.Active)
	if err != nil {
		log.Printf("failed to execute db.Exec UPDATE_ROUTE: %v", err)
		return errors.New("failed to update route")
	}

	return nil
}
