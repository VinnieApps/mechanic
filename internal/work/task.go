package work

// ParsedTask is a task that unmarshalled from a YAML and correctly
// parsed to a specific task.
type ParsedTask struct {
	Original map[string]interface{}
	Task     Task
}

// Task represents a task to be executed
type Task interface {
	// Execute should be an idempotent action
	Execute(packageBase string) error
}
