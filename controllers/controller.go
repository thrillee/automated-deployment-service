package controllers

import (
	"github.com/thrillee/automated-deployment-service/db"
)

type Controller struct {
	db *db.MongoDB
}

func CreateController(db *db.MongoDB) *Controller {
	return &Controller{
		db: db,
	}
}

func (c *Controller) HealthCheck() error {
	return c.db.PingDB()
}
