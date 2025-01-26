package service

import (
	"errors"

	model "github.com/demkowo/rbac/models"
	"github.com/google/uuid"
)

type RbacRepo interface {
	Add(*model.Rbac) error
	Delete(*model.Rbac) error
	Find() ([]*model.Rbac, error)
}

type RolesRepo interface {
	Add(*model.Role) error
	Delete(string) error
	ExistsByID(uuid.UUID) (bool, error)
	Find() ([]*model.Role, error)
	FindByRoute(uuid.UUID) ([]*model.Role, error)
	Update(*model.Role) error
}

type RoutesRepo interface {
	AddActive([]*model.Route) error
	Add(*model.Route) error
	Delete(uuid.UUID) error
	ExistsByID(uuid.UUID) (bool, error)
	Find() ([]*model.Route, error)
	FindByRole(uuid.UUID) ([]*model.Route, error)
	SetInactive() error
	Update(*model.Route) error
}

type Rbac interface {
	AddRbac(*model.Rbac) error
	DeleteRbac(*model.Rbac) error
	FindRbac() ([]*model.Rbac, error)

	AddRole(*model.Role) error
	DeleteRole(string) error
	FindRoles() ([]*model.Role, error)
	FindRolesByRoute(uuid.UUID) ([]*model.Role, error)
	UpdateRole(*model.Role) error

	AddActiveRoutes([]*model.Route) error
	AddRoute(*model.Route) error
	DeleteRoute(uuid.UUID) error
	FindRoutes() ([]*model.Route, error)
	FindRoutesByRole(uuid.UUID) ([]*model.Route, error)
	UpdateRoute(*model.Route) error
	SetRoutesInactive() error
}

type rbac struct {
	rbac   RbacRepo
	roles  RolesRepo
	routes RoutesRepo
}

func NewRbac(rbacRepo RbacRepo, rolesRepo RolesRepo, routesRepo RoutesRepo) Rbac {
	return &rbac{
		rbac:   rbacRepo,
		roles:  rolesRepo,
		routes: routesRepo,
	}
}

func (s *rbac) AddRbac(rbac *model.Rbac) error {
	routeExists, err := s.routes.ExistsByID(rbac.RouteID)
	if err != nil {
		return err
	}
	if !routeExists {
		return errors.New("route does not exist")
	}

	roleExists, err := s.roles.ExistsByID(rbac.RoleID)
	if err != nil {
		return err
	}
	if !roleExists {
		return errors.New("role does not exist")
	}

	return s.rbac.Add(rbac)
}

func (s *rbac) DeleteRbac(auth *model.Rbac) error {
	err := s.rbac.Delete(auth)
	if err != nil {
		return err
	}

	return nil
}

func (s *rbac) FindRbac() ([]*model.Rbac, error) {
	return s.rbac.Find()
}

func (s *rbac) AddRole(role *model.Role) error {
	if role.ID == uuid.Nil {
		role.ID = uuid.New()
	}

	return s.roles.Add(role)
}

func (s *rbac) DeleteRole(role string) error {
	return s.roles.Delete(role)
}

func (s *rbac) FindRoles() ([]*model.Role, error) {
	return s.roles.Find()
}

func (s *rbac) FindRolesByRoute(routeID uuid.UUID) ([]*model.Role, error) {
	return s.roles.FindByRoute(routeID)
}

func (s *rbac) UpdateRole(role *model.Role) error {
	return s.roles.Update(role)
}

func (s *rbac) AddActiveRoutes(routes []*model.Route) error {
	return s.routes.AddActive(routes)
}

func (s *rbac) AddRoute(route *model.Route) error {
	if route.ID == uuid.Nil {
		route.ID = uuid.New()
	}

	return s.routes.Add(route)
}

func (s *rbac) DeleteRoute(routeID uuid.UUID) error {
	return s.routes.Delete(routeID)
}

func (s *rbac) FindRoutes() ([]*model.Route, error) {
	return s.routes.Find()
}

func (s *rbac) FindRoutesByRole(roleID uuid.UUID) ([]*model.Route, error) {
	return s.routes.FindByRole(roleID)
}

func (s *rbac) SetRoutesInactive() error {
	return s.routes.SetInactive()
}

func (s *rbac) UpdateRoute(route *model.Route) error {
	return s.routes.Update(route)
}
