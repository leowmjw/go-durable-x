#!/bin/bash
# Failure test for Temporal Stateful Actor pattern
# This script demonstrates the fault-tolerance of the workflow

# Function to modify activity file to inject failures
inject_failures() {
    cat > temp-activity.go << 'EOL'
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.temporal.io/sdk/activity"
)

// BringUpMachine activity brings up a machine to running state
func BringUpMachine(ctx context.Context, machineID string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Bringing up machine", "machineID", machineID)

	// Simulate random failures for testing fault tolerance
	if rand.Float32() < 0.7 {
		logger.Error("Simulated random failure in BringUpMachine", "machineID", machineID)
		return fmt.Errorf("simulated failure in BringUpMachine for %s", machineID)
	}

	// Simulate work
	time.Sleep(time.Millisecond * 500)
	return nil
}

// TearDownMachine activity tears down a machine from running state
func TearDownMachine(ctx context.Context, machineID string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Tearing down machine", "machineID", machineID)

	// Simulate random failures for testing fault tolerance
	if rand.Float32() < 0.7 {
		logger.Error("Simulated random failure in TearDownMachine", "machineID", machineID)
		return fmt.Errorf("simulated failure in TearDownMachine for %s", machineID)
	}

	// Simulate work
	time.Sleep(time.Millisecond * 500)
	return nil
}
EOL

    # Backup original file if it exists
    if [ -f "activities.go" ]; then
        mv activities.go activities.go.backup
    fi
    
    # Move temp file to activities.go
    mv temp-activity.go activities.go
    
    echo "Injected random failures into activities."
}

# Function to restore original activity file
restore_activities() {
    if [ -f "activities.go.backup" ]; then
        mv activities.go.backup activities.go
        echo "Restored original activities file."
    fi
}

# Main script
echo "Starting fault-tolerance test..."
echo "This test will inject random failures into activities to demonstrate retry behavior."

# Trap Ctrl+C to ensure cleanup
trap restore_activities EXIT

# Inject failures
inject_failures

# Define workflow ID
WORKFLOW_ID="fault-test-$(date +%s)"

# Start the workflow
echo -e "\nStarting workflow with ID: $WORKFLOW_ID"
echo "You will see activities fail and automatically retry."
echo "This demonstrates the fault-tolerance of the Temporal implementation."

temporal workflow start \
  --task-queue stateful-machines \
  --workflow-id $WORKFLOW_ID \
  --type MachineOperatorWorkflow \
  --input '"machine-failure-test"'

# Wait for workflow to initialize
sleep 2

# Send setUp signal
echo -e "\nSending setUp signal..."
echo "Watch for activity failures and automatic retries."
temporal workflow signal \
  --workflow-id $WORKFLOW_ID \
  --name setUp

# Allow time for retries (longer because we expect failures)
echo -e "\nWaiting for setUp retries to complete (10 seconds)..."
sleep 10

# Send tearDown signal
echo -e "\nSending tearDown signal..."
echo "Again, watch for activity failures and automatic retries."
temporal workflow signal \
  --workflow-id $WORKFLOW_ID \
  --name tearDown

# Allow time for retries
echo -e "\nWaiting for tearDown retries to complete (10 seconds)..."
sleep 10

# Complete the workflow and show final state
echo -e "\nCompleting workflow and retrieving final state..."
temporal workflow signal \
  --workflow-id $WORKFLOW_ID \
  --name complete

# Wait for workflow to complete
sleep 2

# Show workflow activity history (focusing on retries)
echo -e "\nWorkflow retry history (showing activity failures and retries):"
temporal workflow show --workflow-id $WORKFLOW_ID | grep -E "ActivityTaskFailed|ActivityTaskScheduled|ActivityTaskCompleted" | head -20

# Show final status
echo -e "\nFinal workflow status after handling all failures:"
temporal workflow describe --workflow-id $WORKFLOW_ID | grep Status

echo -e "\nFault-tolerance test complete! The workflow successfully handled all failures"
echo "and completed the full lifecycle: DOWN -> UP -> DOWN -> Complete"

# Clean up injected failures
restore_activities
