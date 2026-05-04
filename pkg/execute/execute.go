package execute

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"syscall"

	"github.com/puzpuzpuz/xsync/v3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type CommandExecutor struct {
	ctx      context.Context
	commands *xsync.MapOf[string, *Command]
}

type CommandResult struct {
	Lapse    float32 `json:"lapse"`
	ExitCode int     `json:"exitCode"`
	Stdout   string  `json:"stdout"`
	Stderr   string  `json:"stderr"`
	Error    string  `json:"error"`
	Aborted  bool    `json:"aborted"`
}

func (ce *CommandExecutor) SetContext(ctx context.Context) {
	ce.ctx = ctx
	ce.commands = xsync.NewMapOf[string, *Command]()
}

func (ce *CommandExecutor) Run(program string, options []string) string {
	id := ce.generateId()
	ce.commands.Store(id, NewCommand(program, options))

	go ce.dispatch(id)

	return id
}

func (ce *CommandExecutor) RunAndOutput(program string, options []string, hideWindow bool) CommandResult {
	var (
		errMsg  string
		command = NewCommand(program, options)
	)

	if hideWindow {
		command.cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x08000000,
		}
	}

	if err := command.Run(); err != nil {
		errMsg = err.Error()
	}

	return CommandResult{
		command.Lapse(),
		command.cmd.ProcessState.ExitCode(),
		command.stdout.String(),
		command.stderr.String(),
		errMsg,
		command.stopped,
	}
}

func (ce *CommandExecutor) Abort(id string) error {
	if task, ok := ce.commands.Load(id); !ok {
		return errors.New("execute: id not found")
	} else {
		if err := task.Stop(); err != nil {
			return errors.Join(err, errors.New("execute: abort failed"))
		}
		return nil
	}
}

func (ce *CommandExecutor) dispatch(id string) {
	command, ok := ce.commands.Load(id)
	if !ok {
		runtime.EventsEmit(ce.ctx, "execute:exited", id, CommandResult{
			Error: "execute: id not found",
		})
		return
	}

	var errMsg string
	if err := command.Run(); err != nil {
		errMsg = err.Error()
	}

	runtime.EventsEmit(ce.ctx, "execute:exited", id, CommandResult{
		command.Lapse(),
		command.cmd.ProcessState.ExitCode(),
		command.DecodeStdout(),
		command.DecodeStderr(),
		errMsg,
		command.stopped,
	})
}

func (ce CommandExecutor) generateId() string {
	id := ""
	for id == "" {
		b := make([]byte, 4)
		if _, err := rand.Read(b); err != nil {
			continue
		}

		tmpId := hex.EncodeToString(b)
		if _, ok := ce.commands.Load(tmpId); ok {
			continue
		}

		id = tmpId
	}
	return id
}
