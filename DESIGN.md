# Background Job Processing System

## 1. Overview

The goal of this project is to build a reliable background job processing system using Go and PostgreSQL.

The system allows clients to submit jobs through an HTTP API. Jobs are persisted in PostgreSQL and executed asynchronously by a configurable pool of workers.

The system should support concurrent execution, automatic retries, failure recovery, and graceful shutdown.

### Initial job types

- `email` – simulate sending an email.
- `report` – simulate generating a report.

## 2. Requirements

### Functional requirements

- Clients can submit jobs through an HTTP API.
- The API validates job payloads according to their type.
- Jobs are persisted in PostgreSQL.
- The API immediately returns a job ID without waiting for execution.
- Clients can retrieve job status using the job ID.
- Multiple workers can execute jobs concurrently.
- Workers support different job types through handlers.
- Failed jobs can be retried automatically.
- Jobs that are abandoned due to worker crashes can be recovered.

### Non-functional requirements

- Reliable job processing.
- At-least-once execution semantics.
- Idempotent handling of operations where possible.
- Configurable worker count.
- Graceful application shutdown.
- Structured logging and error handling.

## 3. High-Level Architecture

The system consists of the following components:

### HTTP API
- Accepts job creation requests.
- Validates incoming requests.
- Exposes job status endpoints.
- Returns HTTP 202 Accepted when a job is successfully submitted.

### Job Service
- Handles job creation.
- Applies job-related business rules.
- Coordinates with the Queue Service.

### Queue Service
- Enqueues jobs.
- Claims available jobs.
- Acknowledges successful executions.
- Reschedules failed jobs.
- Recovers jobs with expired leases.

### Worker Pool / Orchestrator
- Starts and manages workers.
- Controls the number of concurrent workers.
- Coordinates execution results.
- Applies retry policies.
- Handles worker lifecycle and graceful shutdown.

### Worker
- Requests available jobs from the Queue Service.
- Executes jobs using the appropriate handler.
- Reports execution results.
- Does not contain job-specific business logic.

### Handler Registry
- Resolves the appropriate handler based on job type.
- Supports registering different job handlers.

Initial handlers:
- `EmailHandler`
- `ReportHandler`

### Repository
- Provides database access.
- Executes queries and transactions.
- Is used by the Queue Service.

### PostgreSQL
- Persists jobs and their execution state.
- Provides atomic job-claiming operations.
- Acts as the durable storage layer for the queue.

## 4. Job Lifecycle

Supported job statuses:

- `pending`
- `running`
- `completed`
- `failed`

Expected lifecycle:

    Pending
       |
       v
     Running
      /   \
     /     \
    v       v
Completed  Failed
              |
              v
        Retry available?
           /      
          v        v
       Pending   Failed

A job is considered completed only after successful execution has been acknowledged.

A failed job is rescheduled if retry attempts remain. Otherwise, its final status becomes `failed`.

## 5. Database Design

Initial `jobs` table:

| Column | Type | Description |
|---|---|---|
| id | UUID | Unique job identifier |
| type | VARCHAR | Job type |
| payload | JSONB | Job-specific data |
| status | VARCHAR | Current job status |
| attempts | INTEGER | Number of execution attempts |
| max_attempts | INTEGER | Maximum retry attempts |
| scheduled_at | TIMESTAMPTZ | Earliest execution time |
| locked_by | VARCHAR | Worker currently owning the job |
| locked_until | TIMESTAMPTZ | Lease expiration time |
| last_error | TEXT | Most recent execution error |
| created_at | TIMESTAMPTZ | Job creation time |
| updated_at | TIMESTAMPTZ | Last update time |
| completed_at | TIMESTAMPTZ | Completion time |

The database will enforce constraints for valid job types, statuses, and attempt counts.

The payload will use JSONB because different job types require different data structures. Payload validation will be handled according to the job type at the API boundary, with safe decoding and validation during execution as well.

## 6. Job Execution and Concurrency

Each worker independently requests the next available job through the Queue Service.

The Queue Service uses PostgreSQL transactions and `SELECT ... FOR UPDATE SKIP LOCKED` to ensure that multiple workers cannot claim the same job simultaneously.

Once a job is claimed:

1. Its status changes to `running`.
2. The worker receives the job.
3. The Handler Registry selects the appropriate handler.
4. The handler executes the job.
5. The execution result is reported.
6. The job is marked completed or rescheduled.

The number of workers will be configurable.

## 7. Retry Mechanism

The Orchestrator coordinates retry decisions.

When a worker reports an execution failure:

1. The Orchestrator evaluates the retry policy.
2. It checks whether attempts remain.
3. If attempts remain, it requests the Queue Service to reschedule the job.
4. The Queue Service updates the attempt count and scheduled execution time.
5. If no attempts remain, the job is marked as failed.

Initial retry strategy: exponential backoff.

Example:

| Attempt | Retry delay |
|---|---|
| 1 | 5 seconds |
| 2 | 10 seconds |
| 3 | 20 seconds |
| 4 | 40 seconds |

Retry configuration should be adjustable.

## 8. Worker Recovery and Job Leases

When a worker claims a job, the Queue Service assigns:

- `locked_by`
- `locked_until`

Workers periodically renew their leases while processing long-running jobs.

If a worker crashes or stops renewing its lease, the job becomes eligible for recovery after the lease expires.

A lease token or fencing mechanism should be considered to prevent an outdated worker from modifying a job after its lease has been reassigned.

Worker health monitoring and job lease monitoring are separate concerns.

## 9. Idempotency and Delivery Guarantees

The system provides at-least-once execution semantics.

A job may execute more than once, for example, when a worker successfully performs an operation but crashes before acknowledging completion.

Handlers should be designed to be idempotent wherever possible.

The job ID can be used as an idempotency key.

For external operations, such as sending emails, provider-supported idempotency mechanisms or deduplication strategies may be required.

Exactly-once execution is not guaranteed.

## 10. Graceful Shutdown

When the application receives a shutdown signal:

1. The Orchestrator stops assigning new work to workers.
2. Workers are allowed to finish their currently executing jobs within a configured grace period.
3. Remaining jobs must not be acknowledged as completed unless execution succeeded.
4. Unfinished jobs must become recoverable through lease expiration or explicit release.
5. Database connections and other resources are closed cleanly.

## 11. Initial API Endpoints

### Create a job

`POST /jobs`

Example request:

    {
      "type": "email",
      "payload": {
        "to": "user@example.com",
        "subject": "Welcome",
        "body": "Welcome to our platform!"
      }
    }

Example response:

    {
      "id": "job-123",
      "status": "pending"
    }

HTTP status: `202 Accepted`

### Get job status

`GET /jobs/{id}`

Example response:

    {
      "id": "job-123",
      "type": "email",
      "status": "completed",
      "attempts": 1
    }

## 12. Initial Project Scope

The first version will include:

- Go HTTP API.
- PostgreSQL persistence.
- Generic job model with JSONB payload.
- Email and report handlers.
- Configurable worker pool.
- Concurrent job execution.
- Atomic job claiming.
- Retry with exponential backoff.
- Job leases and recovery.
- Graceful shutdown.

Future improvements:

- Job cancellation.
- Dead-letter queue.
- Metrics and monitoring.
- Multiple application instances.
- Advanced idempotency and deduplication.
- Job prioritization.
