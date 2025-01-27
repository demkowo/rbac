package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	handler "github.com/demkowo/rbac/handlers"

	postgres "github.com/demkowo/rbac/repositories/postgres"
	service "github.com/demkowo/rbac/services"

	_ "github.com/lib/pq"
)

const (
	portNumber = ":5000"
)

var (
	dbConnection = os.Getenv("DB_CLIENT")
	router       = gin.Default()
)

func Start() {
	db, err := sql.Open("postgres", dbConnection)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	rbacRepo := postgres.NewRbac(db)
	rolesRepo := postgres.NewRoles(db)
	routesRepo := postgres.NewRoutes(db)
	rbacService := service.NewRbac(rbacRepo, rolesRepo, routesRepo)
	rbacHandler := handler.NewRbac(rbacService)
	addRbacRoutes(rbacHandler)

	CreateTables(db)

	rbacHandler.MarkActiveRoutes(router)

	router.Run(portNumber)
}
