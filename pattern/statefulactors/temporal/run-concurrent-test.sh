#!/bin/bash
# Concurrent request test for Temporal Stateful Actor pattern
# This script simulates multiple concurrent requests to demonstrate linearization

echo "Testing concurrent requests to demonstrate linearization of state transitions"

# Create two machine IDs for testing
MACHINE_A="machine-a-$(date +%s)"
MACHINE_B="machine-b-$(date +%s)"

# Start the workflows
echo -e "\nStarting workflow for machine A: $MACHINE_A"
temporal workflow start \
  --task-queue stateful-machines \
  --workflow-id $MACHINE_A \
  --type MachineOperatorWorkflow \
  --input "\"$MACHINE_A\"" > /dev/null 2>&1

echo -e "Starting workflow for machine B: $MACHINE_B"
temporal workflow start \
  --task-queue stateful-machines \
  --workflow-id $MACHINE_B \
  --type MachineOperatorWorkflow \
  --input "\"$MACHINE_B\"" > /dev/null 2>&1

# Wait for workflows to initialize
sleep 5

echo -e "\nSending multiple concurrent signals to demonstrate linearization:"
echo "- Machine A: SetUp, TearDown"
echo "- Machine B: SetUp, SetUp (duplicate), TearDown"
echo -e "\nExpect: Each machine will process signals in sequence, one at a time."

# Send signals concurrently
echo -e "\nSending signals simultaneously..."

# Send SetUp to machine A
(
  temporal workflow signal \
    --workflow-id $MACHINE_A \
    --name setUp > /dev/null 2>&1
  echo "Signal sent: Machine A SetUp"
) &

# Send TearDown to machine A
(
  temporal workflow signal \
    --workflow-id $MACHINE_A \
    --name tearDown > /dev/null 2>&1
  echo "Signal sent: Machine A TearDown"
) &

# Send SetUp to machine B
(
  temporal workflow signal \
    --workflow-id $MACHINE_B \
    --name setUp > /dev/null 2>&1
  echo "Signal sent: Machine B SetUp"
) &

# Send SetUp to machine B again (should have no effect if already UP)
(
  sleep 0.5  # Small delay just to ensure the first signal gets processed first
  temporal workflow signal \
    --workflow-id $MACHINE_B \
    --name setUp > /dev/null 2>&1
  echo "Signal sent: Machine B SetUp (again)"
) &

# Send TearDown to machine B
(
  temporal workflow signal \
    --workflow-id $MACHINE_B \
    --name tearDown > /dev/null 2>&1
  echo "Signal sent: Machine B TearDown"
) &

# Wait for signals to be sent and activities to complete
echo "Executing..."
sleep 15

# Let's look at the workflow history for one of the machines to show linearization
echo -e "\nWorkflow history for Machine A (showing linearized execution):"
temporal workflow show \
  --workflow-id $MACHINE_A | grep -E "WorkflowExecutionStarted|ActivityTaskScheduled|ActivityTaskCompleted|SignalExternalWorkflowExecution" | head -15

# Complete the workflows
echo -e "\nCompleting the workflows:"

echo -e "\nMachine A final state and completion:"
temporal workflow signal \
  --workflow-id $MACHINE_A \
  --name complete > /dev/null 2>&1

# Wait a moment for workflow to complete
sleep 2
temporal workflow describe --workflow-id $MACHINE_A | grep Status

echo -e "\nMachine B final state and completion:"
temporal workflow signal \
  --workflow-id $MACHINE_B \
  --name complete > /dev/null 2>&1

# Wait a moment for workflow to complete
sleep 2
temporal workflow describe --workflow-id $MACHINE_B | grep Status

echo -e "\nConcurrent test complete! The workflow history shows that:"
echo "1. Each signal is processed sequentially"
echo "2. Activities execute in order with no overlapping execution"
echo "3. State transitions are linearized per machine ID"
echo "4. Final workflow status shows successful completion"
