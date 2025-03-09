# README

## Stateful Actors and State Machines
[<img src="https://raw.githubusercontent.com/restatedev/img/refs/heads/main/show-code.svg">](src/statefulactors/machineoperator.go)

This example implements a State Machine with a Virtual Object.

* The object holds the state of the state machine and defines the methods
  to transition between the states.
* The object's unique id identifies the state machine. Many parallel state
  machines exist, but only state machine (object) exists per id.

* The _single-writer-per-key_ characteristic of virtual objects ensures
  that one state transition per state machine is in progress at a time.
  Additional transitions are enqueued for that object, while a transition
  for a machine is still in progress.
* The state machine behaves like a **virtual stateful actor**.

* The state machine transitions (object methods) themselves run with
  _durable execution_, so they recover with all partial progress
  and intermediate state.

What you get by this are _linearized interactions_ with your state machine,
avoiding accidental state corruption and concurrency issues.

<details>
<summary><strong>Running the example</strong></summary>

1. [Start the Restate Server](https://docs.restate.dev/develop/local_dev) in a separate shell: `restate-server`
2. Start the service: `go run ./src/statefulactors`
3. Register the services (with `--force` to override the endpoint during **development**): `restate -y deployments register --force localhost:9080`

Invoke the state machine transitions like
```shell
curl -X POST localhost:8080/MachineOperator/my-machine/SetUp
```

To illustrate the concurrency safety here, send multiple requests without waiting on
results and see how they play out sequentially per object (state machine).
Copy all the curl command lines below and paste them to the terminal together.
You will see both from the later results (in the terminal with the curl commands) and in
the log of the service that the requests queue per object key and safely execute
unaffected by crashes and recoveries.

```shell
(curl -X POST localhost:8080/MachineOperator/a/SetUp &)
(curl -X POST localhost:8080/MachineOperator/a/TearDown &)
(curl -X POST localhost:8080/MachineOperator/b/SetUp &)
(curl -X POST localhost:8080/MachineOperator/b/SetUp &)
(curl -X POST localhost:8080/MachineOperator/b/TearDown &)
echo "executing..."
```

<details>
<summary>View logs</summary>

```
2025/01/07 15:43:39 WARN Accepting requests without validating request signatures; handler access must be restricted
2025/01/07 15:43:48 INFO Handling invocation method=MachineOperator/TearDown invocationID=inv_1dceKvwtEc2n73auyTOQxa4kxIlWRcptG9
2025/01/07 15:43:48 INFO Beginning transition to down: a
👻 A failure happened!2025/01/07 15:43:48 ERROR Invocation returned a non-terminal failure method=MachineOperator/TearDown invocationID=inv_1dceKvwtEc2n73auyTOQxa4kxIlWRcptG9 err="a failure happened"
2025/01/07 15:43:48 INFO Handling invocation method=MachineOperator/SetUp invocationID=inv_174rq2A9bm3T0atyOfXqUIVy47VcOx80Jb
2025/01/07 15:43:48 INFO Beginning transition to up: b
2025/01/07 15:43:48 INFO Handling invocation method=MachineOperator/TearDown invocationID=inv_1dceKvwtEc2n73auyTOQxa4kxIlWRcptG9
2025/01/07 15:43:48 INFO Beginning transition to down: a
2025/01/07 15:43:53 INFO Done transitioning to up: b
2025/01/07 15:43:53 INFO Invocation completed successfully method=MachineOperator/SetUp invocationID=inv_174rq2A9bm3T0atyOfXqUIVy47VcOx80Jb
2025/01/07 15:43:53 INFO Handling invocation method=MachineOperator/SetUp invocationID=inv_174rq2A9bm3T1DTVZFHn9ClEXySogMbf8J
2025/01/07 15:43:53 INFO Invocation completed successfully method=MachineOperator/SetUp invocationID=inv_174rq2A9bm3T1DTVZFHn9ClEXySogMbf8J
2025/01/07 15:43:53 INFO Handling invocation method=MachineOperator/TearDown invocationID=inv_174rq2A9bm3T3sOGmjdHa6cfEb2eFhNyaB
2025/01/07 15:43:53 INFO Beginning transition to down: b
👻 A failure happened!2025/01/07 15:43:53 ERROR Invocation returned a non-terminal failure method=MachineOperator/TearDown invocationID=inv_174rq2A9bm3T3sOGmjdHa6cfEb2eFhNyaB err="a failure happened"
2025/01/07 15:43:53 INFO Done transitioning to down: a
2025/01/07 15:43:53 INFO Invocation completed successfully method=MachineOperator/TearDown invocationID=inv_1dceKvwtEc2n73auyTOQxa4kxIlWRcptG9
2025/01/07 15:43:53 INFO Handling invocation method=MachineOperator/SetUp invocationID=inv_1dceKvwtEc2n4c2TwvTC3GhkUrqOH9PvCV
2025/01/07 15:43:53 INFO Beginning transition to up: a
2025/01/07 15:43:53 INFO Handling invocation method=MachineOperator/TearDown invocationID=inv_174rq2A9bm3T3sOGmjdHa6cfEb2eFhNyaB
2025/01/07 15:43:53 INFO Beginning transition to down: b
2025/01/07 15:43:58 INFO Done transitioning to up: a
2025/01/07 15:43:58 INFO Invocation completed successfully method=MachineOperator/SetUp invocationID=inv_1dceKvwtEc2n4c2TwvTC3GhkUrqOH9PvCV
2025/01/07 15:43:58 INFO Done transitioning to down: b
2025/01/07 15:43:58 INFO Invocation completed successfully method=MachineOperator/TearDown invocationID=inv_174rq2A9bm3T3sOGmjdHa6cfEb2eFhNyaB
```

</details>
</details>


