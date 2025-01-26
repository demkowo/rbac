package app

import (
	"net/http"

	handler "github.com/demkowo/rbac/handlers"
	"github.com/gin-gonic/gin"
)

func addRbacRoutes(h handler.Rbac) {
	routes := router.Group("/api/v1/routes")
	{
		routes.POST("add", h.AddRoute)
		routes.POST("add-active", func(c *gin.Context) {
			res, err := h.AddActiveRoutes(router)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "activation routes failed"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"routes": res})
		})
		routes.GET("delete/:route_id", h.DeleteRoute)
		routes.GET("find", h.FindRoutes)
		routes.GET("find-by-role/:role_id", h.FindRoutesByRole)
		routes.POST("update", h.UpdateRoute)
	}

	roles := router.Group("/api/v1/roles")
	{
		roles.POST("add", h.AddRole)
		roles.POST("delete", h.DeleteRole)
		roles.GET("find", h.FindRoles)
		roles.GET("find-by-route/:route_id", h.FindRolesByRoute)
		roles.POST("update", h.UpdateRole)
	}

	rbac := router.Group("/api/v1/rbac")
	{
		rbac.POST("add", h.AddRbac)
		rbac.POST("delete", h.DeleteRbac)
		rbac.GET("find", h.FindRbac)
	}
}
