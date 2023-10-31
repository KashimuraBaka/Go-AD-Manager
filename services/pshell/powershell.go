package pshell

import (
	"errors"
	"regexp"
	"strings"

	PS "github.com/bhendo/go-powershell"
	"github.com/bhendo/go-powershell/backend"
	"github.com/bhendo/go-powershell/middleware"
)

type PowerShell struct {
	shell   PS.Shell
	session middleware.Middleware
}

func (ps *PowerShell) CreateSession(remoteAddr string, username string, password string) error {
	// 准备远程会话配置
	if ps.session != nil {
		ps.session.Exit()
	}

	config := middleware.NewSessionConfig()
	config.ComputerName = remoteAddr
	config.Credential = &middleware.UserPasswordCredential{
		Username: username,
		Password: password,
	}

	// 通过包装已存在的会话中间件来创建一个新的shell
	session, err := middleware.NewSession(ps.shell, config)
	if err != nil {
		return err
	}

	ps.session = session
	return nil
}

func (ps *PowerShell) CloseSession() {
	if ps.session != nil {
		ps.session.Exit()
		ps.session = nil
	}
}

func (ps *PowerShell) Execute(commands ...string) (string, string, error) {
	var stdout string
	var stderr string
	if ps.session == nil {
		return stdout, stderr, errors.New("powershell session is null")
	}
	if len(commands) > 0 {
		command := strings.Join(commands, " ")
		command = regexp.MustCompile(`\n\s*`).ReplaceAllString(command, " ")
		command = strings.Trim(command, " ")
		return ps.session.Execute(command)
	} else {
		return stdout, stderr, errors.New("command is null")
	}
}

func CreatePowershell() (*PowerShell, error) {
	// 挑选一个后台
	back := &backend.Local{}

	// 开启一个本地的 powershell 进程
	shell, err := PS.New(back)
	if err != nil {
		return nil, err
	}

	return &PowerShell{shell: shell}, nil
}
