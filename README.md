# Flowork

A workflow runner for simple, linear workflows.

## Installation

**Build from source:** `go build -o flowork ./cmd/cli` - this will produce
an executable called `flowork` in the root directory.

## Concepts

### Task

A task is a single, atomic computation, usually one program execution.
If a task is applied to multiple sets of inputs, then multiple instances
of the task may run in parallel.

### Workflow

A workflow is a sequence of tasks that are to be executed. The output of
each task is made available as the input of the next task in the sequence.
If there is more than one input available to a task, it will be fanned out
to allow parallel execution.

## Tutorial
