package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/ii64/go-dlv-manager/conf"
	"github.com/ii64/go-dlv-manager/log"
	"github.com/ii64/go-dlv-manager/server"
)

// -gcflags="all=-N -l" on Go 1.10 or later, -gcflags="-N -l"
// dlv dap --check-go-version=false --listen=127.0.0.1:9999 --log=true --log-output=debugger,debuglineerr,gdbwire,lldout,rpc
// dlv dap --check-go-version=false --listen=127.0.0.1:9999 --log=true --log-output=debugger,debuglineerr,gdbwire,lldout,rpc

func main() {
	var exitCode = 0
	defer func() {
		os.Exit(exitCode)
	}()

	conf.Init()
	defer conf.Close()

	log.Init()
	defer log.Close()

	logger := log.Logger()
	logger.Info().Msg("init")
	defer logger.Info().Msg("exited")

	srv, err := server.New(&server.Option{
		Program:     conf.Default.Program,
		ProgramArgs: conf.Default.ProcessArgs,
	})
	if err != nil {
		logger.Err(err).Msg("server init")
	}

	// ---

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	var task taskc

	task.Task(func(wg *sync.WaitGroup) {
		defer cancel() // abort
		defer wg.Done()
		err = srv.ListenAndServe(ctx, conf.Default.Addr)
		if err != nil {
			logger.Err(err).Msg("server closed")
		}
	})

	task.Task(func(wg *sync.WaitGroup) {
		defer wg.Done()
		<-ctx.Done()

		// release resource
		var err error
		err = srv.Close()
		if err != nil {
			logger.Err(err).Msg("server close")
		}
	})

	finishChan := make(chan struct{}, 1)
	go func() {
		task.Wait()
		finishChan <- struct{}{}
	}()

	intrCnt := 0
	for {
		select {
		case <-finishChan:
			exitCode = 0
			goto exitLoop
		case s := <-sig:
			if intrCnt > 1 {
				exitCode = 1
				logger.Error().Msg("abort graceful shutdown")
				goto exitLoop
			}
			fmt.Println("received signal:", s)
			cancel()
			intrCnt++
		}
		continue
	exitLoop:
		break
	}
}

type taskc struct {
	wg sync.WaitGroup
}

func (t *taskc) Task(f func(wg *sync.WaitGroup)) {
	t.wg.Add(1)
	go f(&t.wg)
}

func (t *taskc) Wait() {
	t.wg.Wait()
}
