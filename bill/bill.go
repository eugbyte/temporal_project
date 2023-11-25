package billhandler

import (
	"context"
	"fmt"
	"time"

	debug "encore.app/internal/logger"

	db "encore.app/internal/db/bill"
	temporalbill "encore.app/internal/temporal/bill"
	"encore.dev"
	gonanoid "github.com/matoous/go-nanoid/v2"
	cmap "github.com/orcaman/concurrent-map/v2"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var envName = encore.Meta().Environment.Name
var taskQ = envName + "task-queue"

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
	// save workflows for signals to intercept
	workflowIDs cmap.ConcurrentMap[string, string]
}

// entry point, dependency injection
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

	w.RegisterWorkflow(workflows.CreateBill)
	w.RegisterActivity(activities.CreateBill)

	return &Handler{
		billService: billService,
		client:      c,
		worker:      w,
		workflowIDs: cmap.New[string](),
	}, nil
}

func (s *Handler) Shutdown(force context.Context) {
	s.client.Close()
	s.worker.Stop()
}

// action - "create" or "add"
func genWorkFlowID(action string, billID string) string {
	randID, _ := gonanoid.New()
	return fmt.Sprintf("%s-bill-%s-%s", action, billID, randID)
}
