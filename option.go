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
	"path/filepath"
	"strings"
)

// ServiceOption provides a utility for setting service specific options during
// initialisation. A service will always be created with sensible default
// values.
type ServiceOption func(*Service)

// WithExecs defines any number of executables that will be wrapped as a service
// and managed by the service manager on the OS. If no executable is provided, one
// will be resolved using [os.Executable], ultimately identifying where [svc.New]
// was invoked. If this behaviour is not desirable, then setting this option is
// paramount.
//
// When defining an executable path, all arguments must be included e.g.
//
//	path/to/executable --arg1 --arg2=value
func WithExecs(paths []string) ServiceOption {
	return func(s *Service) {
		execs := make([]executable, 0, len(paths))

		for _, path := range paths {
			trimmedPath := strings.Trim(path, " ")
			if trimmedPath == "" {
				continue
			}

			execs = append(execs, buildExecPath(trimmedPath))
		}

		s.execs = execs
	}
}

func buildExecPath(path string) executable {
	exePath, args, _ := strings.Cut(path, " ")

	exe := executable{
		exec: filepath.Base(exePath),
		path: filepath.Dir(exePath),
	}

	if args != "" {
		exe.arguments = strings.Split(args, " ")
	}

	return exe
}

// WithName sets a friendly name for the service. This name will ultimately be
// used when creating the service definition. If no name is provided, the first
// targeted executable will be used.
func WithName(name string) ServiceOption {
	return func(s *Service) {
		s.name = strings.Trim(name, " ")
	}
}

// WithDescription sets a custom description for the service when building the
// service definition file during installation. If no description is provided,
// a suitable default will be used.
//
//	Process <EXECUTABLE> wrapped using the tiny svc library by Purple Clay
func WithDescription(desc string) ServiceOption {
	return func(s *Service) {
		s.description = strings.Trim(desc, " ")
	}
}
