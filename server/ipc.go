package server

import (
	"fmt"
	"math/rand"
	"net"
	"os/exec"
	"time"

	"github.com/ii64/go-dlv-manager/log"
	"github.com/rs/zerolog"
)

type IPC struct {
	program     string
	programArgs []string
	logger      zerolog.Logger

	ipcProgramAddr string
	cmd            *exec.Cmd
	// stdoutBuffer   bytes.Buffer
	ipcConn net.Conn
}

func newIPC(program string, args []string) (*IPC, error) {
	ins := &IPC{
		program:     program,
		programArgs: args,
		logger:      log.Logger(),
	}

	port := (rand.Int() + 500) % int(^uint16(0))
	ins.ipcProgramAddr = fmt.Sprintf("127.0.0.1:%d", port)

	args = append(args,
		"--listen="+ins.ipcProgramAddr)
	ins.logger.Info().
		Str("p", program).
		Strs("args", args).
		Msg("ipc start")
	ins.cmd = exec.Command(program, args...)

	// ins.cmd.Stdout = &ins.stdoutBuffer
	return ins, ins.cmd.Start()
}

func (ins *IPC) Close() error {
	ins.logger.Info().Msg("ipc end")
	if conn := ins.ipcConn; conn != nil {
		conn.Close()
	}
	return ins.cmd.Process.Kill()
}

func (ins *IPC) Dial() (conn net.Conn, err error) {
	i := 0
	for {
		conn, err = net.Dial("tcp", ins.ipcProgramAddr)
		if err == nil || i > 20 {
			break
		}
		time.Sleep(time.Second * 1)
		i++
	}
	ins.ipcConn = conn
	return
}
