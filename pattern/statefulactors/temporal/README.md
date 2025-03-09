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

# Signal the workflow to tear down the machine
temporal workflow signal \
  --workflow-id machine-operator-1 \
  --name tearDown

# Complete the workflow and check final status
temporal workflow signal \
  --workflow-id machine-operator-1 \
  --name complete

# View workflow history and final status
temporal workflow show --workflow-id machine-operator-1
temporal workflow describe --workflow-id machine-operator-1 | grep Status
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

# Send setUp signal
echo "\nSending setUp signal..."
temporal workflow signal \
  --workflow-id $WORKFLOW_ID \
  --name setUp

# Allow time for activity to complete
sleep 3

# Send tearDown signal
echo "\nSending tearDown signal..."
temporal workflow signal \
  --workflow-id $WORKFLOW_ID \
  --name tearDown

# Allow time for activity to complete
sleep 3

# Complete the workflow and show final state
echo "\nCompleting workflow and retrieving final state..."
temporal workflow signal \
  --workflow-id $WORKFLOW_ID \
  --name complete

# Wait for workflow to complete
sleep 2

# Show workflow history and completion status
echo "\nWorkflow history:"
temporal workflow show --workflow-id $WORKFLOW_ID | grep -E "WorkflowExecutionStarted|ActivityTaskCompleted" | head -15

echo "\nFinal workflow status:"
temporal workflow describe --workflow-id $WORKFLOW_ID | grep Status
```

### Integration Testing Scripts

This project includes three pre-configured integration test scripts that match the testing scenarios in the original Restate example:

#### 1. Basic Integration Test

The `run-integration-test.sh` script demonstrates the normal operation of the workflow without failures:

```bash
# Make the script executable if needed
chmod +x run-integration-test.sh

# Run the integration test
./run-integration-test.sh
```

This test will take a workflow through the complete state transition sequence:
- Start in DOWN state
- Transition to UP (via setUp signal)
- Transition to DOWN (via tearDown signal)
- Complete the workflow

#### 2. Fault-Tolerance Test

The `run-failure-test.sh` script demonstrates how the workflow handles failures and automatic retries:

```bash
# Make the script executable if needed
chmod +x run-failure-test.sh

# Run the fault-tolerance test
./run-failure-test.sh
```

This script:
1. Temporarily injects random failures into the activities (70% failure rate)
2. Executes the workflow with these failures
3. Shows how Temporal automatically retries failed activities
4. Demonstrates that the workflow eventually reaches the correct state despite failures
5. Restores the original activity implementations afterward

#### 3. Concurrent Request Test

The `run-concurrent-test.sh` script demonstrates how multiple concurrent state transitions are properly linearized per machine ID, exactly matching the testing scenario in the Restate example:

```bash
# Make the script executable if needed
chmod +x run-concurrent-test.sh

# Run the concurrent request test
./run-concurrent-test.sh
```

This script:
1. Creates two separate workflows (Machine A and Machine B)
2. Sends multiple concurrent signals to both workflows:
   - Machine A: SetUp, TearDown
   - Machine B: SetUp, SetUp (duplicate), TearDown
3. Demonstrates that signals are processed sequentially per workflow
4. Shows that duplicate signals have the expected idempotent behavior
5. Displays workflow history to visualize linearized execution

This test directly mirrors the concurrent request test in the Restate example, showing how Temporal provides the same _single-writer-per-key_ guarantees that ensure state transitions are processed one at a time per machine ID.

## Comparison with Restate

### Matching the Restate Example

This implementation directly parallels the Restate stateful actor example found in `/pattern/statefulactors/restate/`. The following table shows how our testing scenarios map to the Restate example:

| Restate Example | Temporal Implementation |
|-----------------|------------------------|
| State Machine (UP/DOWN) | Same state machine with identical states |
| SetUp/TearDown methods | Implemented as activities triggered by signals |
| Single machine per ID | One workflow instance per machine ID |
| Linearized operations | Operations are processed sequentially per workflow |
| Automatic retries on failure | Retry policies on activities with similar behavior |
| Concurrent requests test | `run-concurrent-test.sh` script |
| Random failures for testing | `run-failure-test.sh` script |

### Equivalence of Guarantees

Both implementations provide these important guarantees:

1. **Single-writer-per-key**: In Restate, this is provided by the actor model. In Temporal, this is provided by having one workflow instance per machine ID, with signals processed sequentially.

2. **Durable execution**: Both systems recover with all partial progress and intermediate state after failures.

3. **Linearized interactions**: Both systems ensure state transitions happen one at a time per machine, avoiding accidental state corruption and concurrency issues.

### Simplified Interface Approach

The implementation has been simplified to match the Restate pattern more closely:

1. **Focus on final state**: Instead of querying state during transitions, we focus on the final workflow status upon completion, matching Restate's approach of returning the final state when operations complete.

2. **Signal-driven operations**: Operations are fully signal-driven, with the workflow execution history showing the linearized sequence of events.

3. **Workflow history as proof**: The workflow history serves as proof of linearization, showing that operations are processed one at a time in the order they were received.

### Key Implementation Differences

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

## Comprehensive Testing Guide

This project provides a progression of tests from unit tests to full integration tests, allowing developers to understand the stateful actor pattern from different perspectives.

### Testing Progression

1. **Unit Tests** - Validate core logic without Temporal machinery
   ```bash
   # Run the StateMachineLogic test only
   go test -run TestStateMachineLogic
   ```

2. **Workflow Tests** - Test workflow behavior with mocked activities
   ```bash
   # Run workflow completion test
   go test -run TestWorkflowCompletion
   ```

3. **Activity Tests** - Test activity implementations
   ```bash
   # Run activity implementation test
   go test -run TestActivityImplementations
   ```

4. **Error Handling Tests** - Validate fault-tolerance
   ```bash
   # Run error handling tests
   go test -run TestWorkflowWithErrors
   ```

5. **Basic Integration Test** - End-to-end test with a real Temporal server
   ```bash
   # Start Temporal server if not running
   docker run --rm -d -p 7233:7233 -p 8233:8233 temporalio/temporal:latest
   
   # In one terminal, start the worker
   go run .
   
   # In another terminal, run the integration test
   ./run-integration-test.sh
   ```

6. **Fault-Tolerance Integration Test** - Test actual retry behavior
   ```bash
   # With worker and Temporal server running
   ./run-failure-test.sh
   ```

7. **Concurrent Requests Test** - Demonstrate linearized execution (matches Restate example)
   ```bash
   # With worker and Temporal server running
   ./run-concurrent-test.sh
   ```
   
   This test directly corresponds to the concurrent request example in the Restate README, where multiple operations are submitted simultaneously to show how they are linearized per machine ID.

### Visualization

After running the integration tests, view the workflows in the Temporal Web UI:

1. Open [http://localhost:8233](http://localhost:8233) in your browser
2. Navigate to the "Default" namespace
3. Find your workflow execution by ID
4. Explore the execution history, which shows:
   - Activity executions and retries
   - Signal receipts
   - State transitions
   - Workflow completion

This visualization is particularly valuable for understanding how Temporal provides durability and visibility into distributed processes.
