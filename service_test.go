/*
Copyright (c) 2022 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package svc_test

import (
	"errors"
	"syscall"
	"testing"
	"time"

	"github.com/purpleclay/svc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type process struct {
	errRun       error
	errInterrupt error
}

func (p *process) Run() error {
	if p.errRun != nil {
		return p.errRun
	}

	// Process doesn't need to be a blocking operation
	return nil
}

func (p *process) Interrupt() error {
	return p.errInterrupt
}

func TestServiceRunProcessError(t *testing.T) {
	proc := &process{
		errRun: errors.New("process run error"),
	}
	service, err := svc.New(proc)
	require.NoError(t, err)

	err = service.Run()
	assert.EqualError(t, err, "process run error")
}

func TestServiceRunInterrupt(t *testing.T) {
	service, err := svc.New(&process{})
	require.NoError(t, err)

	raiseSignal(t, 200*time.Millisecond, syscall.SIGINT)
	err = service.Run()

	assert.NoError(t, err)
}

func TestServiceRunTerminate(t *testing.T) {
	service, err := svc.New(&process{})
	require.NoError(t, err)

	raiseSignal(t, 200*time.Millisecond, syscall.SIGTERM)
	err = service.Run()

	assert.NoError(t, err)
}

func raiseSignal(t *testing.T, after time.Duration, sig syscall.Signal) {
	t.Helper()

	ticker := time.NewTicker(after)
	go func() {
		for range ticker.C {
			syscall.Kill(syscall.Getpid(), sig)
		}
	}()
}

func TestServiceRunInterruptError(t *testing.T) {
	proc := &process{
		errInterrupt: errors.New("process interrupt error"),
	}
	service, err := svc.New(proc)
	require.NoError(t, err)

	raiseSignal(t, 200*time.Millisecond, syscall.SIGINT)
	err = service.Run()

	assert.EqualError(t, err, "process interrupt error")
}
