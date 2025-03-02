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

# Initial state check
echo -e "\nInitial state check:"
temporal workflow query \
  --workflow-id $WORKFLOW_ID \
  --query-type getStatus

# Send setUp signal
echo -e "\nSending setUp signal..."
temporal workflow signal \
  --workflow-id $WORKFLOW_ID \
  --name setUp

# Allow time for activity to complete
sleep 3

# Check state after setUp
echo -e "\nState after setUp:"
temporal workflow query \
  --workflow-id $WORKFLOW_ID \
  --query-type getStatus

# Send tearDown signal
echo -e "\nSending tearDown signal..."
temporal workflow signal \
  --workflow-id $WORKFLOW_ID \
  --name tearDown

# Allow time for activity to complete
sleep 3

# Check state after tearDown
echo -e "\nState after tearDown:"
temporal workflow query \
  --workflow-id $WORKFLOW_ID \
  --query-type getStatus

# Complete the workflow
echo -e "\nCompleting workflow..."
temporal workflow signal \
  --workflow-id $WORKFLOW_ID \
  --name complete

# Wait for workflow to complete
sleep 2

# Check workflow status
echo -e "\nWorkflow completion status:"
temporal workflow describe --workflow-id $WORKFLOW_ID | grep Status

echo -e "\nIntegration test complete!"
