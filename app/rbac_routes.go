package app

import (
	"net/http"

	handler "github.com/demkowo/rbac/handlers"
	auth "github.com/demkowo/utils/auth"
	"github.com/gin-gonic/gin"
)

func addRbacRoutes(h handler.Rbac) {
	router.POST("/api/v1/routes/:service", h.AddExternalRoutes)

	// === ROUTES ===
	routes := router.Group("/api/v1/routes", auth.AuthMiddleware())
	{
		routes.GET("", h.FindRoutes)
		routes.GET("role/:role_id", h.FindRoutesByRole)
		routes.POST("", h.AddRoute)
		routes.POST("mark-active", func(c *gin.Context) {
			res, err := h.MarkActiveRoutes(router)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "activation routes failed"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"routes": res})
		})
		routes.PUT("/:route_id", h.UpdateRoute)
		routes.DELETE(":route_id", h.DeleteRoute)
	}

	// === ROLES ===
	roles := router.Group("/api/v1/roles", auth.AuthMiddleware())
	{
		roles.GET("", h.FindRoles)
		roles.GET("route/:route_id", h.FindRolesByRoute)
		roles.POST("routes", h.FindRolesByRoutes)
		roles.POST("", h.AddRole)
		roles.PUT("/:role_id", h.UpdateRole)
		roles.DELETE("/:role_id", h.DeleteRole)
	}

	// === RBAC (role-route relation) ===
	rbac := router.Group("/api/v1/rbac", auth.AuthMiddleware())
	{
		rbac.GET("", h.FindRbac)
		rbac.POST("", h.AddRbac)
		rbac.DELETE("", h.DeleteRbac)
	}
}
