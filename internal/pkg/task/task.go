package task

import (
	"fmt"
	"log/slog"
)

// RunAll applies the given tasks, using the given runner, to the given
// volume, in the order that they are provided.
func RunAll(r Runner, tasks []*Instance, v Volume) error {
	for _, inst := range tasks {
		slog.Info("running task instance", "name", inst.Task.Name, "id", inst.ID, "volume", v)

		// TODO: We could copy output names to input names to make tasks easier to re-use

		err := r.Run(inst, v)
		if err != nil {
			slog.Error("task instance failed", "name", inst.Task.Name, "id", inst.ID, "volume", v)
			return fmt.Errorf("failed to run all tasks on %s: %w", v, err)
		}
	}

	return nil
}
