package runtime

import (
	"context"
	"fmt"
	"os/exec"
)

func (m *Manager) startNative(_ context.Context, spec ModuleSpec, proc *process) error {
	command, err := joinCommand(spec.Command, spec.Args)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	proc.cancelRun = cancel

	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	cmd.Stdout = proc.logFile
	cmd.Stderr = proc.logFile
	if spec.WorkDir != "" {
		cmd.Dir = spec.WorkDir
	}
	cmd.Env = buildEnv(spec.Env)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start native process: %w", err)
	}
	proc.cmd = cmd
	proc.status.State = StateRunning
	proc.status.PID = cmd.Process.Pid

	go m.waitForExit(proc, cmd.Wait)
	return nil
}

func (m *Manager) stopNative(ctx context.Context, proc *process) error {
	if proc.cancelRun != nil {
		proc.cancelRun()
	}
	if proc.cmd == nil {
		return nil
	}

	done := make(chan error, 1)
	go func() {
		done <- proc.cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("native process exit: %w", err)
		}
		m.markStopped(proc, nil)
	case <-ctx.Done():
		if proc.cmd.Process != nil {
			if killErr := proc.cmd.Process.Kill(); killErr != nil {
				return fmt.Errorf("kill native process: %w", killErr)
			}
		}
		return fmt.Errorf("timeout stopping native process")
	}

	return nil
}
