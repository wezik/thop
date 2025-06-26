package executor

import (
	"log/slog"
	"os"
	"os/exec"
)

type CommandExecutor interface {
	Execute(cmd *exec.Cmd) (string, int, error)
	ExecuteInteractive(cmd *exec.Cmd) (int, error)
}

type ShellExecutor struct {
	Logger *slog.Logger
}

func (s *ShellExecutor) Execute(cmd *exec.Cmd) (string, int, error) {
	s.Logger.Info(cmd.String())
	res, err := cmd.Output()
	return string(res), cmd.ProcessState.ExitCode(), err
}

func (s *ShellExecutor) ExecuteInteractive(cmd *exec.Cmd) (int, error) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return cmd.ProcessState.ExitCode(), err
}
