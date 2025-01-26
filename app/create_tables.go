package app

import (
	"database/sql"
	"log"
)

const (
	RBAC_TABLE_EXIST   = "SELECT to_regclass('public.rbac')"
	ROLES_TABLE_EXIST  = "SELECT to_regclass('public.roles')"
	ROUTES_TABLE_EXIST = "SELECT to_regclass('public.routes')"

	RBAC_CREATE_TABLE = `
        CREATE TABLE IF NOT EXISTS rbac (
            route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
            role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
            PRIMARY KEY (route_id, role_id)
        );
    `

	ROLES_CREATE_TABLE = `
		CREATE TABLE IF NOT EXISTS roles (
			id UUID PRIMARY KEY, 
			name VARCHAR(255) NOT NULL UNIQUE
		);
	`

	ROUTES_CREATE_TABLE = `
        CREATE TABLE IF NOT EXISTS routes (
            id UUID PRIMARY KEY,
            method VARCHAR(10) NOT NULL,
            path VARCHAR(255) NOT NULL,
            active BOOLEAN NOT NULL DEFAULT TRUE,
            UNIQUE (method, path)
        );
	`
)

func CreateTables(db *sql.DB) {
	if !checkRbacExists(db) {
		createRbac(db)
	}

	if !checkRolesExists(db) {
		createRoles(db)
	}

	if !checkRoutesExists(db) {
		createRoutes(db)
	}

	log.Println("tables rbac, roles and routes are ready to go")
}

func checkRbacExists(db *sql.DB) bool {
	var tableName sql.NullString
	err := db.QueryRow(RBAC_TABLE_EXIST).Scan(&tableName)
	if err != nil {
		log.Panic("failed to check rbac table existence")
	}

	return tableName.Valid
}

func checkRolesExists(db *sql.DB) bool {
	var tableName sql.NullString
	err := db.QueryRow(ROLES_TABLE_EXIST).Scan(&tableName)
	if err != nil {
		log.Panic("failed to check roles table existence")
	}

	return tableName.Valid
}

func checkRoutesExists(db *sql.DB) bool {
	var tableName sql.NullString
	err := db.QueryRow(ROUTES_TABLE_EXIST).Scan(&tableName)
	if err != nil {
		log.Panic("failed to check routes table existence")
	}

	return tableName.Valid
}

func createRbac(db *sql.DB) {
	_, err := db.Exec(RBAC_CREATE_TABLE)
	if err != nil {
		log.Panic("failed to create rbac table")
	}
}

func createRoles(db *sql.DB) {
	_, err := db.Exec(ROLES_CREATE_TABLE)
	if err != nil {
		log.Panic("failed to create roles table")
	}
}

func createRoutes(db *sql.DB) {
	_, err := db.Exec(ROUTES_CREATE_TABLE)
	if err != nil {
		log.Panic("failed to create routes table")
	}
}
