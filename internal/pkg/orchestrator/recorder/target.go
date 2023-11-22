package recorder

import "github.com/glesica/flowork/internal/pkg/orchestrator"

type Target interface {
	Success(job *orchestrator.Job) error
	Failure(job *orchestrator.Job) error
}
