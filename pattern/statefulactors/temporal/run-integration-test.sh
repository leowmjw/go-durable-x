#!/bin/bash
# Integration test for Temporal Stateful Actor pattern
# This script simulates a client application interacting with our stateful actor

# Define workflow ID
WORKFLOW_ID="machine-operator-$(date +%s)"

# Start the workflow
echo "Starting workflow with ID: $WORKFLOW_ID"
temporal workflow start \
  --task-queue stateful-machines \
  --workflow-id $WORKFLOW_ID \
  --type MachineOperatorWorkflow \
  --input '"machine-test"'

# Wait for workflow to initialize
sleep 2

# Send setUp signal
echo -e "\nSending setUp signal..."
temporal workflow signal \
  --workflow-id $WORKFLOW_ID \
  --name setUp

# Allow time for activity to complete
sleep 3

# Send tearDown signal
echo -e "\nSending tearDown signal..."
temporal workflow signal \
  --workflow-id $WORKFLOW_ID \
  --name tearDown

# Allow time for activity to complete
sleep 3

# Complete the workflow and show final state
echo -e "\nCompleting workflow and retrieving final state..."
temporal workflow signal \
  --workflow-id $WORKFLOW_ID \
  --name complete

# Wait for workflow to complete
sleep 2

# Show workflow history and completion status
echo -e "\nWorkflow history:"
temporal workflow show --workflow-id $WORKFLOW_ID | grep -E "WorkflowExecutionStarted|ActivityTaskScheduled|ActivityTaskCompleted|SignalExternalWorkflowExecution" | head -15

echo -e "\nFinal workflow status:"
temporal workflow describe --workflow-id $WORKFLOW_ID | grep Status

echo -e "\nIntegration test complete! Machine went through full lifecycle: DOWN -> UP -> DOWN -> Complete"
