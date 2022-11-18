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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	// this should be used if [svc.Run] is not going to be invoked
	noProcess Process = nil
)

func TestWithName(t *testing.T) {
	tests := []struct {
		name     string
		with     string
		expected string
	}{
		{
			name:     "NoLeadingTrailingWhitespace",
			with:     "testing",
			expected: "testing",
		},
		{
			name:     "LeadingAndTrailingWhitespace",
			with:     "    testing     ",
			expected: "testing",
		},
		{
			name:     "BlankStringCurrentProcess",
			with:     "                 ",
			expected: "svc.test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(noProcess, WithName(tt.with))
			require.NoError(t, err)

			require.Equal(t, tt.expected, s.name)
		})
	}
}

func TestWithDescription(t *testing.T) {
	tests := []struct {
		name     string
		with     string
		expected string
	}{
		{
			name:     "NoLeadingTrailingWhitespace",
			with:     "a test description",
			expected: "a test description",
		},
		{
			name:     "LeadingAndTrailingWhitespace",
			with:     "    a test description     ",
			expected: "a test description",
		},
		{
			name:     "BlankStringDefaultDescription",
			with:     "                 ",
			expected: "Process svc.test wrapped using the tiny svc library by Purple Clay",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(noProcess, WithDescription(tt.with))
			require.NoError(t, err)

			require.Equal(t, tt.expected, s.description)
		})
	}
}

func TestWithExecs(t *testing.T) {
	tests := []struct {
		name     string
		with     []string
		expected []executable
	}{
		{
			name: "NoLeadingTrailingWhitespace",
			with: []string{
				"/path/to/executable1",
				"/path/to/executable2 --arg1 --arg2=value",
			},
			expected: []executable{
				{
					path: "/path/to",
					exec: "executable1",
				},
				{
					path:      "/path/to",
					exec:      "executable2",
					arguments: []string{"--arg1", "--arg2=value"},
				},
			},
		},
		{
			name: "LeadingAndTrailingWhitespace",
			with: []string{" /path/to/executable3 --arg1      "},
			expected: []executable{
				{
					path:      "/path/to",
					exec:      "executable3",
					arguments: []string{"--arg1"},
				},
			},
		},
		{
			name: "BlankExecPathIgnored",
			with: []string{"", "/path/to/executable4"},
			expected: []executable{
				{
					path: "/path/to",
					exec: "executable4",
				},
			},
		},
		{
			name: "BlankExecDefaultExecDetected",
			with: []string{""},
			expected: []executable{
				{
					path: workingDir(t),
					exec: "svc.test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := New(noProcess, WithExecs(tt.with))
			require.NoError(t, err)

			require.ElementsMatch(t, tt.expected, s.execs)
		})
	}
}

func workingDir(t *testing.T) string {
	t.Helper()

	dir, err := os.Executable()
	require.NoError(t, err)

	return filepath.Dir(dir)
}
