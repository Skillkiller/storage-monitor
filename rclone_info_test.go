package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestGetRcloneVersion_Success(t *testing.T) {
	tests := []struct {
		name            string
		output          string
		expectedVersion string
	}{
		{
			name: "rclone v1.60.1-DEV",
			output: `rclone v1.60.1-DEV
- os/version: debian 13.1 (64 bit)
- os/kernel: 6.12.41+deb13-amd64 (x86_64)
- os/type: linux
- os/arch: amd64
- go/version: go1.23.2
- go/linking: dynamic
- go/tags: none`,
			expectedVersion: "rclone v1.60.1-DEV",
		},
		{
			name: "rclone v1.71.1",
			output: `rclone v1.71.1
- os/version: debian 13.1 (64 bit)
- os/kernel: 6.12.41+deb13-arm64 (aarch64)
- os/type: linux
- os/arch: arm64 (ARMv8 compatible)
- go/version: go1.25.1
- go/linking: static
- go/tags: none`,
			expectedVersion: "rclone v1.71.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExecutor := func(name string, arg ...string) ([]byte, error) {
				return []byte(tt.output), nil
			}

			version, err := getRcloneVersionWithExecutor("/usr/bin/rclone", mockExecutor)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if version != tt.expectedVersion {
				t.Errorf("Expected version %q, got %q", tt.expectedVersion, version)
			}
		})
	}
}
func TestGetRcloneVersion_NotRclone(t *testing.T) {
	mockExecutor := func(name string, arg ...string) ([]byte, error) {
		return []byte("restic 0.18.1 compiled with go1.25.1 on linux/arm64"), nil
	}

	_, err := getRcloneVersionWithExecutor("/usr/bin/fake", mockExecutor)

	if err == nil {
		t.Fatal("Expected error when executable is not rclone")
	}

	expected := "rclone executable is not rclone"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error containing %q, got %q", expected, err.Error())
	}
}

func TestGetRcloneVersion_ExitError(t *testing.T) {
	mockExecutor := func(name string, arg ...string) ([]byte, error) {
		return nil, &exec.ExitError{}
	}

	_, err := getRcloneVersionWithExecutor("/usr/bin/rclone", mockExecutor)

	if err == nil {
		t.Fatal("Expected error when rclone exits with error")
	}

	expected := "rclone exited with code"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error containing %q, got %q", expected, err.Error())
	}
}

func TestGetRcloneVersion_GenericError(t *testing.T) {
	mockExecutor := func(name string, arg ...string) ([]byte, error) {
		return nil, fmt.Errorf("command not found")
	}

	_, err := getRcloneVersionWithExecutor("/usr/bin/nonexistent", mockExecutor)

	if err == nil {
		t.Fatal("Expected error when command fails")
	}

	expected := "failed to execute rclone"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error containing %q, got %q", expected, err.Error())
	}
}

func TestGetRcloneVersion_EmptyOutput(t *testing.T) {
	mockExecutor := func(name string, arg ...string) ([]byte, error) {
		return []byte(""), nil
	}

	_, err := getRcloneVersionWithExecutor("/usr/bin/rclone", mockExecutor)

	if err == nil {
		t.Fatal("Expected error when output is empty")
	}

	expected := "rclone executable is not rclone"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error containing %q, got %q", expected, err.Error())
	}
}
