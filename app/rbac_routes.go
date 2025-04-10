package app

import (
	"net/http"

	handler "github.com/demkowo/rbac/handlers"
	auth "github.com/demkowo/utils/auth"
	"github.com/gin-gonic/gin"
)

func addRbacRoutes(h handler.Rbac) {
	routes := router.Group("/api/v1/routes", auth.AuthMiddleware())
	{
		routes.POST("", h.AddRoute)
		routes.POST("/:service", h.AddExternalRoutes)
		routes.POST("mark-active", func(c *gin.Context) {
			res, err := h.MarkActiveRoutes(router)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "activation routes failed"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"routes": res})
		})
		routes.DELETE(":route_id", h.DeleteRoute)
		routes.GET("", h.FindRoutes)
		routes.GET("role/:role_id", h.FindRoutesByRole)
		routes.PUT("/:route_id", h.UpdateRoute)
	}

	roles := router.Group("/api/v1/roles", auth.AuthMiddleware())
	{
		roles.POST("", h.AddRole)
		roles.DELETE("/:role_id", h.DeleteRole)
		roles.GET("", h.FindRoles)
		roles.GET("route/:route_id", h.FindRolesByRoute)
		roles.PUT("/:role_id", h.UpdateRole)
	}

	rbac := router.Group("/api/v1/rbac", auth.AuthMiddleware())
	{
		rbac.POST("", h.AddRbac)
		rbac.DELETE("", h.DeleteRbac)
		rbac.GET("", h.FindRbac)
	}
}
