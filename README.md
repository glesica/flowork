# Flowork

A workflow runner for simple, linear workflows.

## Installation

**Build from source:** `go build -o flowork ./cmd/cli` - this will produce
an executable called `flowork` in the root directory.

## Concepts

Flowork attempts to be simple. Some use-cases won't fit into the Flowork
model, and that's fine. For those cases that do fit, users only need to
understand a few concepts in order to be productive.

### Orchestrator

The machine actually running the Flowork tool.

### Task

A task is a single, atomic computation, usually one program execution.
If a task is applied to multiple sets of inputs, then multiple instances
of the task may run in parallel.

### Task Runner

An implementation of the `task.Runner` interface that knows how to run a
single task on some sort of computing infrastructure. It is also
responsible for gathering input files and recovering output files based
on the storage mechanism used by its infrastructure.

Note that runners may run Flowork on remote machines, using Docker, in
order to copy files around to avoid files moving through the orchestrator
or placing requirements on worker machines beyond having Docker installed.

### Workflow

A workflow is a sequence of tasks that are to be executed. The output of
each task is made available as the input of the next task in the sequence.
If there is more than one input available to a task, it will be fanned out
to allow parallel execution.

## Tutorial
