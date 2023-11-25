package billhandler

import (
	"context"
	"fmt"
	"time"

	debug "encore.app/internal/logger"

	db "encore.app/internal/db/bill"
	temporalbill "encore.app/internal/temporal/bill"
	"encore.dev"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var envName = encore.Meta().Environment.Name
var taskQ = envName + "task-queue"

const workFlowID = "bill-id"

var logger = debug.Logger

type BillService interface {
	Create(billID string) (db.Bill, error)
	Add(billID string, date time.Time, item string, amount float64) (db.Bill, error)
	Close(billID string) error
	Get(billID string) (db.Bill, error)
}

//encore:service
type Handler struct {
	billService BillService
	client      client.Client
	worker      worker.Worker
}

// dependency injection
func initHandler() (*Handler, error) {
	billService := db.New()

	c, err := client.Dial(client.Options{})
	if err != nil {
		return nil, fmt.Errorf("create temporal client: %v", err)
	}

	w := worker.New(c, taskQ, worker.Options{})
	err = w.Start()
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("start temporal worker: %v", err)
	}

	workflows := temporalbill.NewWorkFlow(billService)
	activities := temporalbill.NewActivity(billService)

	w.RegisterWorkflow(workflows.Create)
	w.RegisterActivity(activities.Create)

	return &Handler{
		billService: billService,
		client:      c,
		worker:      w,
	}, nil
}

func (s *Handler) Shutdown(force context.Context) {
	s.client.Close()
	s.worker.Stop()
}
