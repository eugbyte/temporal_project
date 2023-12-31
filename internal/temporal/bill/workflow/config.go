package workflows

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const SignalChannel = "confirm-invoice"

var retrypolicy = &temporal.RetryPolicy{
	InitialInterval:        time.Second,
	BackoffCoefficient:     2.0,
	MaximumInterval:        100 * time.Second,
	MaximumAttempts:        0, // unlimited retries
	NonRetryableErrorTypes: []string{"NonRetryError"},
}

var options = workflow.ActivityOptions{
	// Timeout options specify when to automatically timeout Activity functions.
	StartToCloseTimeout: time.Minute,
	// Optionally provide a customized RetryPolicy.
	// Temporal retries failed Activities by default.
	RetryPolicy: retrypolicy,
}
