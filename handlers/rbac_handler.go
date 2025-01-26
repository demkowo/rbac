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

type Rbac interface {
	AddRbac(*gin.Context)
	DeleteRbac(*gin.Context)
	FindRbac(*gin.Context)

	AddRole(*gin.Context)
	DeleteRole(*gin.Context)
	FindRoles(*gin.Context)
	FindRolesByRoute(*gin.Context)
	UpdateRole(*gin.Context)

	AddActiveRoutes(*gin.Engine) ([]*model.Route, error)
	AddRoute(*gin.Context)
	DeleteRoute(*gin.Context)
	FindRoutes(*gin.Context)
	FindRoutesByRole(*gin.Context)
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
	var auth []*model.Rbac
	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	for _, a := range auth {
		if a.RouteID == uuid.Nil || a.RoleID == uuid.Nil {
			log.Println("RouteID and RoleID are required")
			c.JSON(http.StatusBadRequest, gin.H{"error": "RouteID and RoleID are required"})
			return
		}
	}

	for _, a := range auth {
		if err := h.service.AddRbac(a); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Authorizations added successfully"})
}

func (h *rbac) DeleteRbac(c *gin.Context) {
	var auth []*model.Rbac
	if err := c.ShouldBindJSON(&auth); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	for _, a := range auth {
		if err := h.service.DeleteRbac(a); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Authorizations deleted successfully"})
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
	var role *model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if err := h.service.AddRole(role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"role": role})
}

func (h *rbac) DeleteRole(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if err := h.service.DeleteRole(req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Role %s deleted successfully", req.Name)})
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
	routeID, err := uuid.Parse(c.Param("route_id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route ID"})
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
	if err := c.ShouldBindJSON(&role); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if err := h.service.UpdateRole(role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}

func (h *rbac) AddActiveRoutes(router *gin.Engine) ([]*model.Route, error) {
	if err := h.service.SetRoutesInactive(); err != nil {
		log.Println("setting routes inactive failed", err)
	}

	var routes []*model.Route

	for _, route := range router.Routes() {
		r := &model.Route{
			ID:     uuid.New(),
			Method: route.Method,
			Path:   route.Path,
			Active: true,
		}

		routes = append(routes, r)
	}

	if err := h.service.AddActiveRoutes(routes); err != nil {
		log.Println("adding routes failed", err)
		return nil, err
	}
	return routes, nil
}

func (h *rbac) AddRoute(c *gin.Context) {
	var route *model.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if err := h.service.AddRoute(route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"AddRoute failed": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"route": route})
}

func (h *rbac) DeleteRoute(c *gin.Context) {
	routeID, err := uuid.Parse(c.Param("route_id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid route ID"})
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
	roleID, err := uuid.Parse(c.Param("role_id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	routes, err := h.service.FindRoutesByRole(roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"routes": routes})
}

func (h *rbac) UpdateRoute(c *gin.Context) {
	var route model.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if err := h.service.UpdateRoute(&route); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"route": route})
}
