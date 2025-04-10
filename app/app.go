package app

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	handler "github.com/demkowo/rbac/handlers"

	postgres "github.com/demkowo/rbac/repositories/postgres"
	service "github.com/demkowo/rbac/services"

	_ "github.com/lib/pq"
)

const (
	portNumber = ":5001"
)

var (
	dbConnection = os.Getenv("DB_RBAC")
	router       = gin.Default()
)

func Start() {
	db, err := sql.Open("postgres", dbConnection)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
