package task

// Runner is the runner interface that allows different
// backends to execute workflow tasks.
type Runner interface {
	Name() string
	Run(t Instance) error
}
