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

// Process ...
type Process interface {
	Run() error
	Interrupt() error
}

// Service ...
type Service struct {
	proc Process
	errs chan error
}

// New ...
func New(proc Process) *Service {
	return &Service{
		proc: proc,
	}
}

// Run ...
func (s *Service) Run() error {
	sig := make(chan os.Signal, 1)
	s.errs = make(chan error)

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
