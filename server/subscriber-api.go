package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type newSubscriberPayload struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Reference   string `json:"reference" validate:"required"`
	Subdomain   string `json:"subdomain" validate:"required"`
}

func AddSubscriber(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	newSubscriberPayload := newSubscriberPayload{}

	if err := decoder.Decode(&newSubscriberPayload); err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
	}
}
