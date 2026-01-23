package runtime

import "time"

type ExecutionMode string

type ProcessState string

type HealthState string

const (
	ModeContainer ExecutionMode = "container"
	ModeNative    ExecutionMode = "native"
)

const (
	StateStarting ProcessState = "starting"
	StateRunning  ProcessState = "running"
	StateStopped  ProcessState = "stopped"
	StateFailed   ProcessState = "failed"
)

const (
	HealthUnknown   HealthState = "unknown"
	HealthHealthy   HealthState = "healthy"
	HealthUnhealthy HealthState = "unhealthy"
)

type HealthCheck struct {
	Command  []string
	Interval time.Duration
	Timeout  time.Duration
}

type ModuleSpec struct {
	Name             string
	Mode             ExecutionMode
	Image            string
	Command          []string
	Args             []string
	Env              map[string]string
	WorkDir          string
	ContainerName    string
	ContainerRuntime string
	HealthCheck      HealthCheck
}

type Status struct {
	Name        string
	Mode        ExecutionMode
	PID         int
	ContainerID string
	State       ProcessState
	Health      HealthState
	StartedAt   time.Time
	FinishedAt  *time.Time
	LogPath     string
	LastError   string
}
