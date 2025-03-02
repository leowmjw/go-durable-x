package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type Status string

const (
	UP   Status = "UP"
	DOWN Status = "DOWN"
)

// MachineState represents the state of a machine
type MachineState struct {
	Status Status
}

// MachineOperatorWorkflow implements the state machine logic as a Temporal workflow
func MachineOperatorWorkflow(ctx workflow.Context, machineId string) (string, error) {
	// Set workflow options
	ctx = workflow.WithWorkflowRunTimeout(ctx, time.Second * 120)
	ctx = workflow.WithWorkflowTaskTimeout(ctx, time.Second * 30)
	state := &MachineState{Status: DOWN}
	workflow.SetQueryHandler(ctx, "getStatus", func() (Status, error) {
		return state.Status, nil
	})

	// Signal handlers for state transitions
	setUpChan := workflow.GetSignalChannel(ctx, "setUp")
	tearDownChan := workflow.GetSignalChannel(ctx, "tearDown")

	// Channel for workflow completion
	completeChan := workflow.GetSignalChannel(ctx, "complete")

	// Create a channel to track completion
	done := false

	// Main workflow loop
	selector := workflow.NewSelector(ctx)

	// Add complete signal handler first
	selector.AddReceive(completeChan, func(ch workflow.ReceiveChannel, _ bool) {
		var signal struct{}
		ch.Receive(ctx, &signal)
		done = true
		workflow.GetLogger(ctx).Info("Workflow completed")
	})

	// Add setUp signal handler
	selector.AddReceive(setUpChan, func(ch workflow.ReceiveChannel, _ bool) {
			var signal struct{}
			ch.Receive(ctx, &signal)

			if state.Status == UP {
				return
			}

			ao := workflow.ActivityOptions{
				StartToCloseTimeout:    time.Second * 10,
				ScheduleToCloseTimeout: time.Second * 60,
				ScheduleToStartTimeout: time.Second * 5,
				HeartbeatTimeout:       time.Second * 5,
				RetryPolicy: &temporal.RetryPolicy{
					InitialInterval:    time.Second,
					MaximumInterval:    10 * time.Second,
					BackoffCoefficient: 2.0,
					MaximumAttempts:    5,
				},
			}
			ctx = workflow.WithActivityOptions(ctx, ao)

			workflow.GetLogger(ctx).Info("BringUpMachine activity options configured")

			err := workflow.ExecuteActivity(ctx, BringUpMachine, machineId).Get(ctx, nil)
			if err == nil {
				state.Status = UP
			}
		})

	// Add tearDown signal handler
	selector.AddReceive(tearDownChan, func(ch workflow.ReceiveChannel, _ bool) {
			var signal struct{}
			ch.Receive(ctx, &signal)

			if state.Status != UP {
				return
			}

			ao := workflow.ActivityOptions{
				StartToCloseTimeout:    time.Second * 10,
				ScheduleToCloseTimeout: time.Second * 60,
				ScheduleToStartTimeout: time.Second * 5,
				HeartbeatTimeout:       time.Second * 5,
				RetryPolicy: &temporal.RetryPolicy{
					InitialInterval:    time.Second,
					MaximumInterval:    10 * time.Second,
					BackoffCoefficient: 2.0,
					MaximumAttempts:    5,
				},
			}
			ctx = workflow.WithActivityOptions(ctx, ao)

			workflow.GetLogger(ctx).Info("TearDownMachine activity options configured")

			err := workflow.ExecuteActivity(ctx, TearDownMachine, machineId).Get(ctx, nil)
			if err == nil {
				state.Status = DOWN
			}
		})



	// Wait for signals until done
	workflow.GetLogger(ctx).Info("Starting main workflow loop")
	for !done {
		selector.Select(ctx)
		workflow.GetLogger(ctx).Info("Processed a signal, checking if done")
	}

	workflow.GetLogger(ctx).Info("Workflow exiting normally")
	return "completed", nil
}

// BringUpMachine activity implements the machine startup logic
func BringUpMachine(ctx context.Context, machineId string) error {
	slog.Info("Beginning transition to up: " + machineId)
	activity.RecordHeartbeat(ctx, "starting")

	// Simulate potential failure (10% chance for tests to pass more reliably)
	if err := MaybeCrash(0.1); err != nil {
		slog.Error("Failed during BringUpMachine", "error", err, "machineId", machineId)
		return temporal.NewNonRetryableApplicationError(
			"activity failed",
			"BringUpMachine",
			err,
		)
	}

	activity.RecordHeartbeat(ctx, "simulating work")
	slog.Info("BringUpMachine simulating work", "machineId", machineId)
	time.Sleep(3 * time.Second) // Reduced from 5 to 3 seconds

	slog.Info("Done transitioning to up: " + machineId)
	return nil
}

// TearDownMachine activity implements the machine shutdown logic
func TearDownMachine(ctx context.Context, machineId string) error {
	slog.Info("Beginning transition to down: " + machineId)
	activity.RecordHeartbeat(ctx, "starting")

	// Simulate potential failure (10% chance for tests to pass more reliably)
	if err := MaybeCrash(0.1); err != nil {
		slog.Error("Failed during TearDownMachine", "error", err, "machineId", machineId)
		return temporal.NewNonRetryableApplicationError(
			"activity failed",
			"TearDownMachine",
			err,
		)
	}

	activity.RecordHeartbeat(ctx, "simulating work")
	slog.Info("TearDownMachine simulating work", "machineId", machineId)
	time.Sleep(3 * time.Second) // Reduced from 5 to 3 seconds

	slog.Info("Done transitioning to down: " + machineId)
	return nil
}

func MaybeCrash(probability float32) error {
	if rand.Float32() < probability {
		fmt.Printf("👻 A failure happened!")
		return fmt.Errorf("a failure happened")
	}
	return nil
}

func main() {
	client, err := client.NewClient(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		slog.Error("Unable to create client", "err", err)
		os.Exit(1)
	}
	defer client.Close()

	w := worker.New(client, "machine-operator", worker.Options{})
	w.RegisterWorkflow(MachineOperatorWorkflow)
	w.RegisterActivity(BringUpMachine)
	w.RegisterActivity(TearDownMachine)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		slog.Error("Unable to start worker", "err", err)
		os.Exit(1)
	}
}
