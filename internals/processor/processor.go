package processor

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/thrillee/automated-deployment-service/internals/db"
)

type DeploymentProcessor struct {
	db               *db.MongoDB
	processorFactory *ProcessorHandlerFactory
}

type processorHandlerFunc func(context.Context, interface{}) ProcessorResponse

type ProcessorHandlerFactory struct {
	handlers map[string]processorHandlerFunc
	db       *db.MongoDB
}

type ProcessorResponse struct {
	Success bool
	Message string
	Err     error
	Result  interface{}
}

func getProcessBatchId(ctx context.Context) any {
	return ctx.Value(processs_batch_id_key)
}

func NewProcessor(db *db.MongoDB, processorFactory *ProcessorHandlerFactory) *DeploymentProcessor {
	return &DeploymentProcessor{
		db:               db,
		processorFactory: processorFactory,
	}
}

func (p *ProcessorHandlerFactory) Register(key string, handler processorHandlerFunc) error {
	if _, ok := p.handlers[key]; ok {
		return fmt.Errorf(fmt.Sprintf("Processor Handler Func Registeration Failed: %s already registered", key))
	}

	p.handlers[key] = handler

	logrus.WithFields(logrus.Fields{
		"Key":     key,
		"Handler": handler,
	}).Info("Registering To Factory")

	return nil
}

func (p *ProcessorHandlerFactory) GetHandler(key string) (processorHandlerFunc, error) {
	handler, ok := p.handlers[key]
	if !ok {
		return nil, fmt.Errorf(fmt.Sprintf("%s handler not found", key))
	}
	return handler, nil
}

func CreateProcessorHandlerFactory(db *db.MongoDB) *ProcessorHandlerFactory {
	return &ProcessorHandlerFactory{
		db:       db,
		handlers: map[string]processorHandlerFunc{},
	}
}
