package model

import (
	"github.com/google/uuid"
)

type Route struct {
	ID     uuid.UUID `json:"id"`
	Method string    `json:"method"`
	Path   string    `json:"path"`
	Active bool      `json:"active"`
}

type Role struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Rbac struct {
	RouteID uuid.UUID `json:"route_id"`
	RoleID  uuid.UUID `json:"role_id"`
}
