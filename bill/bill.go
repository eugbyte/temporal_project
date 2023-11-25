package billhandler

import (
	"context"
	"fmt"

	debug "encore.app/internal/logger"

	db "encore.app/internal/db/bill"
	temporalbill "encore.app/internal/temporal/bill"
	"encore.dev"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var envName = encore.Meta().Environment.Name
var taskQ = envName + "task-queue"

var logger = debug.Logger

type BillService interface {
	Create(billID string) (db.Bill, error)
	Add(billID string, billDetail db.TransactionDetail) (db.Bill, error)
	Close(billID string) error
	Get(billID string) (db.Bill, error)
}

//encore:service
type Handler struct {
	billService BillService
	client      client.Client
	worker      worker.Worker
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
	w.RegisterWorkflow(workflows.AddBill)

	w.RegisterActivity(activities.CreateBill)
	w.RegisterActivity(activities.AddBill)

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

func genWorkFlowID(billID string) string {
	randID, _ := gonanoid.New()
	return fmt.Sprintf("bill-%s-%s", billID, randID)
}
