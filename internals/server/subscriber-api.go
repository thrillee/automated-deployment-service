package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thrillee/automated-deployment-service/internals/controllers"
)

type SubscriberAPIs struct {
	controller *controllers.Controller
}

func CreateSubscriberAPI(c *controllers.Controller) *SubscriberAPIs {
	return &SubscriberAPIs{
		controller: c,
	}
}

func (s *SubscriberAPIs) AddSubscriber(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	newSubscriberPayload := controllers.NewSubscriberPayload{}

	if err := decoder.Decode(&newSubscriberPayload); err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
	}

	controllerResponse := s.controller.AddSubscriber(ctx, newSubscriberPayload)
	responseWithJSON(w, 201, controllerResponse)
}
