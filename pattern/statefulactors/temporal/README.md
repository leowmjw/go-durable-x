# Temporal Stateful Actor Pattern

This project demonstrates implementing a stateful actor pattern using Temporal, showcasing fault-tolerance and robust error handling. The implementation simulates a machine operator that can transition between UP and DOWN states.

## Prerequisites

- [Go](https://golang.org/dl/) (version 1.21+)
- [Temporal Server](https://docs.temporal.io/self-hosted-guide) (local deployment or cloud offering)
- [Temporal CLI](https://docs.temporal.io/cli) (for running integration tests)

## Quick Start

### 1. Install Dependencies

```bash
# Install Go dependencies
go mod tidy
```

### 2. Run Temporal Server Locally

If you don't have a Temporal server running, you can start one locally using Docker:

```bash
# Pull and start the Temporal development server
docker run --rm -p 7233:7233 temporalio/temporal:latest
```

### 3. Run Test Suite

```bash
# Run the test suite
go test -v
```

### 4. Run as a Worker (Integration Test)

```bash
# In one terminal, start the worker
go run .

# In another terminal, execute a workflow using the Temporal CLI
temporal workflow start \
  --task-queue stateful-machines \
  --workflow-id machine-operator-1 \
  --type MachineOperatorWorkflow \
  --input '"machine1"'
```

## Understanding the Implementation

This implementation demonstrates a stateful actor pattern with the following components:

### State Machine

The machine can be in one of two states:
- **DOWN**: Initial state, machine is not running
- **UP**: Machine is running

### Signals

The workflow responds to three signals:
1. **setUp**: Transitions the machine from DOWN to UP
2. **tearDown**: Transitions the machine from UP to DOWN
3. **complete**: Completes the workflow

### Activities

Two main activities are implemented:
- **BringUpMachine**: Brings the machine to the UP state
- **TearDownMachine**: Transitions the machine to the DOWN state

## Integration Testing

### Basic Workflow Control

After starting the worker, you can send signals to control the workflow:

```bash
# Signal the workflow to set up the machine
temporal workflow signal \
  --workflow-id machine-operator-1 \
  --name setUp

# Check the current status
temporal workflow query \
  --workflow-id machine-operator-1 \
  --query-type getStatus

# Signal the workflow to tear down the machine
temporal workflow signal \
  --workflow-id machine-operator-1 \
  --name tearDown

# Check the current status again
temporal workflow query \
  --workflow-id machine-operator-1 \
  --query-type getStatus

# Complete the workflow
temporal workflow signal \
  --workflow-id machine-operator-1 \
  --name complete
```

### Simulating the Original Restate Example

To recreate the typical Restate stateful actor pattern experience, you can use the following shell script. This script simulates a client application interacting with our stateful actor:

```bash
#!/bin/bash
# save as run-integration-test.sh and make executable with: chmod +x run-integration-test.sh

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
echo "\nInitial state check:"
temporal workflow query \
  --workflow-id $WORKFLOW_ID \
  --query-type getStatus

# Send setUp signal
echo "\nSending setUp signal..."
temporal workflow signal \
  --workflow-id $WORKFLOW_ID \
  --name setUp

# Allow time for activity to complete
sleep 3

# Check state after setUp
echo "\nState after setUp:"
temporal workflow query \
  --workflow-id $WORKFLOW_ID \
  --query-type getStatus

# Send tearDown signal
echo "\nSending tearDown signal..."
temporal workflow signal \
  --workflow-id $WORKFLOW_ID \
  --name tearDown

# Allow time for activity to complete
sleep 3

# Check state after tearDown
echo "\nState after tearDown:"
temporal workflow query \
  --workflow-id $WORKFLOW_ID \
  --query-type getStatus

# Complete the workflow
echo "\nCompleting workflow..."
temporal workflow signal \
  --workflow-id $WORKFLOW_ID \
  --name complete

# Wait for workflow to complete
sleep 2

# Check workflow status
echo "\nWorkflow completion status:"
temporal workflow describe --workflow-id $WORKFLOW_ID | grep Status
```

### Testing Failure Scenarios

To simulate failures and observe the workflow's fault-tolerance, you can modify the main.go file to introduce artificial failures in the activities:

```go
// In the BringUpMachine or TearDownMachine activities,
// add this code to simulate random failures
if rand.Float32() < 0.5 {
    return fmt.Errorf("simulated random failure")
}
```

Then run the integration test script again to observe how the system handles failures and automatically retries the activities.

## Comparison with Restate

This implementation is inspired by the stateful actor model in Restate but implemented using Temporal. Key differences:

1. **Programming Model**: 
   - Temporal uses a workflow-centric approach with explicit signals and activities
   - Restate uses an actor-centric approach with methods invoked directly on actors

2. **State Management**:
   - Temporal requires explicit state management in workflow code
   - Restate handles state persistence automatically

3. **Durability**:
   - Both systems provide strong consistency and durability guarantees
   - Both can recover from failures and continue execution

For more details on the comparison, see the [LEARNINGS.md](LEARNINGS.md) file.

## Error Handling

The implementation includes sophisticated error handling with:
- Retry policies for transient failures
- State preservation during retries
- Idempotent operations to handle activity retries safely

## Visualizing Workflows

After running the workflow, you can visualize its execution using the Temporal Web UI:

```bash
# If using Docker, open http://localhost:8233 in your browser
```

## Additional Resources

- [Temporal Documentation](https://docs.temporal.io/)
- [Go SDK Documentation](https://docs.temporal.io/dev-guide/go)
- [Full LEARNINGS document](LEARNINGS.md) for detailed insights about this implementation
