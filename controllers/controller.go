package controllers

import "github.com/thrillee/automated-deployment-service/db"

type Controller struct {
	db *db.MongoDB
}
