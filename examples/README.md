# Workflow Examples

This directory contains example workflows and scripts to demonstrate how to use the workflow system.

## Housekeeping Workflow

The housekeeping workflow is a simple example that performs basic housekeeping tasks like cleaning up temporary files and logs. It's implemented in two ways:

1. **Scheduled Workflow**: This workflow is scheduled to run at a specific time (midnight the next day).
2. **Direct Workflow**: This workflow is triggered immediately via the API.

Both workflows perform the same tasks, but they are triggered differently.

## Files

- `housekeeping_workflow.go`: Contains the implementation of the housekeeping workflow.
- `register_scheduled_workflow.go`: Script to register and schedule the housekeeping workflow.
- `trigger_direct_workflow.go`: Script to register and trigger the housekeeping workflow directly via the API.

## How to Use

### Scheduled Workflow

To register and schedule the housekeeping workflow:

```bash
go run examples/register_scheduled_workflow.go
```

This will:
1. Register a new workflow called "Housekeeping (Scheduled)"
2. Set its status to active
3. Schedule it to run at midnight the next day

### Direct Workflow

To register and trigger the housekeeping workflow directly:

```bash
go run examples/trigger_direct_workflow.go
```

This will:
1. Register a new workflow called "Housekeeping (Direct)"
2. Set its status to active
3. Trigger it immediately via the API

## Notes

- Both scripts require a running PostgreSQL database with the correct schema.
- The direct workflow script requires the API server to be running on localhost:3000.
- You'll need to replace the token in the direct workflow script with a valid token.
- The actual workflow execution is simulated in the current implementation.