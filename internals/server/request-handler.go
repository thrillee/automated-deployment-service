package server

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type APIFunc func(context.Context, http.ResponseWriter, *http.Request)

func makeHTTPHandlerFunc(apiFunc APIFunc) http.HandlerFunc {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "requestID", uuid.New().String())

	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"requestId": ctx.Value("requestID"),
			"took":      time.Since(begin),
		}).Info("fetchPrice")
	}(time.Now())

	return func(w http.ResponseWriter, r *http.Request) {
		apiFunc(ctx, w, r)
	}
}
