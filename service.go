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

package svc

import (
	"os"
	"os/signal"
	"syscall"
)

// Process is simply a program (or executable) that is to be managed by
// a service manager within the given OS.
type Process interface {
	// Run provides an entry point for a program that is executed upon
	// startup by the service manager. Any error raised, will result in
	// the process terminating and the service manager handling it in
	// accordance with the service definition.
	//
	// A process entry point doesn't need to be blocking operation, as
	// [svc.Run] will automatically wrap this process within a blocking
	// signal handling loop.
	Run() error

	// Interrupt is called when either a service manager attempts to stop
	// the running process, or the process is intentionally killed by a user.
	// Any process tidying should be carried out here, before the service
	// manager responds.
	Interrupt() error
}

// Service is a process that is designed to run in the background without any
// direct user interaction. Once installed, a process is managed through a
// service manager provided by the given OS.
type Service struct {
	proc Process
	errs chan error
}

// New creates a new service for the given process.
func New(proc Process) *Service {
	return &Service{
		proc: proc,
	}
}

// Run initialises the service and executes the process before blocking
// and waiting for signals to be raised by the service manager.
func (s *Service) Run() error {
	sig := make(chan os.Signal, 1)
	s.errs = make(chan error)

	// Handle both SIGINT (interrupt) and SIGTERM (terminate) signals that will
	// be raised by either a service manager or from a user intentionally killing
	// the process
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		signal.Stop(sig)
		close(sig)
		close(s.errs)
	}()

	go func() {
		if err := s.proc.Run(); err != nil {
			s.errs <- err
		}
	}()

	for {
		select {
		case <-sig:
			return s.proc.Interrupt()
		case err := <-s.errs:
			return err
		}
	}
}
