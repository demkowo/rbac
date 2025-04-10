package handler

import (
	"fmt"
	"log"
	"net/http"

	model "github.com/demkowo/rbac/models"
	service "github.com/demkowo/rbac/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	e error
)

type Rbac interface {
	AddRbac(*gin.Context)
	DeleteRbac(*gin.Context)
	FindRbac(*gin.Context)

	AddRole(*gin.Context)
	DeleteRole(*gin.Context)
	FindRoles(*gin.Context)
	FindRolesByRoute(*gin.Context)
	UpdateRole(*gin.Context)

	AddRoute(*gin.Context)
	AddExternalRoutes(*gin.Context)
	DeleteRoute(*gin.Context)
	FindRoutes(*gin.Context)
	FindRoutesByRole(*gin.Context)
	MarkActiveRoutes(*gin.Engine) ([]model.Route, error)
	UpdateRoute(*gin.Context)
}

type rbac struct {
	service service.Rbac
}

func NewRbac(service service.Rbac) Rbac {
	return &rbac{
		service: service,
	}
}

func (h *rbac) AddRbac(c *gin.Context) {
	var req struct {
		RouteID string `json:"route_id"`
		RoleID  string `json:"role_id"`
	}

	rbac := &model.Rbac{}

	if !bindJSON(c, &req) {
		return
	}

	if rbac.RoleID, e = parseUUID(c, "route_id", req.RoleID); e != nil {
		return
	}

	if rbac.RouteID, e = parseUUID(c, "role_id", req.RouteID); e != nil {
		return
	}

	if err := h.service.AddRbac(rbac); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "RBAC record added successfully"})
}

func (h *rbac) DeleteRbac(c *gin.Context) {
	var req struct {
		RouteID string `json:"route_id"`
		RoleID  string `json:"role_id"`
	}

	rbac := &model.Rbac{}

	if !bindJSON(c, &req) {
		return
	}

	if rbac.RoleID, e = parseUUID(c, "route_id", req.RoleID); e != nil {
		return
	}

	if rbac.RouteID, e = parseUUID(c, "role_id", req.RouteID); e != nil {
		return
	}

	if err := h.service.DeleteRbac(rbac); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "RBAC record deleted successfully"})
}

func (h *rbac) FindRbac(c *gin.Context) {
	res, err := h.service.FindRbac()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rbac records": res})
}

func (h *rbac) AddRole(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}

	if !bindJSON(c, &req) {
		return
	}

	role := &model.Role{
		ID:   uuid.New(),
		Name: req.Name,
	}

	if err := h.service.AddRole(role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"role": role.Name})
}

func (h *rbac) DeleteRole(c *gin.Context) {
	if err := h.service.DeleteRole(c.Param("role_id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted role successfully"})
}

func (h *rbac) FindRoles(c *gin.Context) {
	roles, err := h.service.FindRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

func (h *rbac) FindRolesByRoute(c *gin.Context) {
	routeID, err := parseUUID(c, "route_id", c.Param("route_id"))
	if err != nil {
		return
	}

	roles, err := h.service.FindRolesByRoute(routeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

func (h *rbac) UpdateRole(c *gin.Context) {
	var role *model.Role

	if !bindJSON(c, &role) {
		return
	}

	if role.ID, e = parseUUID(c, "role_id", c.Param("role_id")); e != nil {
		return
	}

	if err := h.service.UpdateRole(role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}

func (h *rbac) AddRoute(c *gin.Context) {
	var req struct {
		Method string `json:"method"`
		Path   string `json:"path"`
		Active bool   `json:"active"`
	}

	if !bindJSON(c, &req) {
		return
	}

	route := &model.Route{
		ID:     uuid.New(),
		Method: req.Method,
		Path:   req.Path,
		Active: req.Active,
	}

	if err := h.service.AddRoute(route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"AddRoute failed": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"route": route})
}

func (h *rbac) AddExternalRoutes(c *gin.Context) {
	var routes []model.Route
	svc := c.Param("service")

	if err := c.ShouldBindJSON(&routes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.SetRoutesInactive(svc); err != nil {
		log.Println("setting routes inactive failed", err)
	}
	if err := h.service.AddActiveRoutes(routes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"routes registered": routes})
}

func (h *rbac) DeleteRoute(c *gin.Context) {
	routeID, err := parseUUID(c, "route_id", c.Param("route_id"))
	if err != nil {
		return
	}

	if err := h.service.DeleteRoute(routeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Route deleted successfully"})
}

func (h *rbac) FindRoutes(c *gin.Context) {
	routes, err := h.service.FindRoutes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"routes": routes})
}

func (h *rbac) FindRoutesByRole(c *gin.Context) {
	roleID, err := parseUUID(c, "role_id", c.Param("role_id"))
	if err != nil {
		return
	}

	routes, err := h.service.FindRoutesByRole(roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"routes": routes})
}

func (h *rbac) MarkActiveRoutes(router *gin.Engine) ([]model.Route, error) {
	var routes []model.Route

	for _, route := range router.Routes() {
		r := model.Route{
			ID:      uuid.New(),
			Method:  route.Method,
			Path:    route.Path,
			Service: "rbac",
			Active:  true,
		}

		routes = append(routes, r)
	}

	if err := h.service.SetRoutesInactive("rbac"); err != nil {
		log.Println("setting routes inactive failed", err)
	}

	if err := h.service.AddActiveRoutes(routes); err != nil {
		log.Println("adding routes failed", err)
		return nil, err
	}
	return routes, nil
}

func (h *rbac) UpdateRoute(c *gin.Context) {
	var route model.Route

	if !bindJSON(c, &route) {
		return
	}

	if route.ID, e = parseUUID(c, "route_id", c.Param("route_id")); e != nil {
		return
	}

	if err := h.service.UpdateRoute(&route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"route": route})
}

func bindJSON(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("invalid JSON data: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid JSON data: %s", err.Error())})
		return false
	}
	return true
}

func parseUUID(c *gin.Context, field, txt string) (uuid.UUID, error) {
	id, err := uuid.Parse(txt)
	if err != nil {
		log.Printf("failed to parse %s: %v", field, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid value: %s", field)})
		return uuid.Nil, err
	}
	return id, nil
}
