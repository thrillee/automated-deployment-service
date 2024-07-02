package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/thrillee/automated-deployment-service/controllers"
	"github.com/thrillee/automated-deployment-service/db"
)

type HttpAPIServer struct {
	ListenAddr string
	dbCon      *db.MongoDB
}

func (h *HttpAPIServer) Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(
		cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		},
	))

	controller := controllers.CreateController(h.dbCon)

	subAPIs := CreateSubscriberAPI(controller)

	router := chi.NewRouter()
	router.Post("/new-sub", makeHTTPHandlerFunc(subAPIs.AddSubscriber))

	r.Mount("/api", router)

	r.Post("/health-check", func(w http.ResponseWriter, r *http.Request) {
		err := controller.HealthCheck()
		if err != nil {
			responseWithError(w, 500, fmt.Sprintf("Health Check Failed: %v", err))
		}
		responseWithError(w, 200, "System is fine")
	})

	srv := &http.Server{
		Handler: r,
		Addr:    h.ListenAddr,
	}
	fmt.Printf("Starting HTTP Server ON %s...", h.ListenAddr)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
