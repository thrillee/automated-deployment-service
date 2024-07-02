package controllers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/thrillee/automated-deployment-service/models"
)

type NewSubscriberPayload struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Reference   string `json:"reference" validate:"required"`
	Subdomain   string `json:"subdomain" validate:"required"`
}

func (c *Controller) AddSubscriber(data NewSubscriberPayload) ControllerResponse {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newSub := models.Subscriber{
		Title:       data.Title,
		Description: data.Description,
		Reference:   uuid.New().String(),
		CallerRef:   data.Reference,
	}
	newSub.Insert(ctx, c.db)

	props := models.SubscriberProp{
		SubscriberID: newSub.ID,
		NKey:         "sub_domain",
		NValue:       data.Subdomain,
	}
	props.Insert(ctx, c.db)

	return ControllerResponse{
		Success: true,
		Message: "Subscriber Added Successfully",
	}
}
