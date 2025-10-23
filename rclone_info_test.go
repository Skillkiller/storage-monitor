package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestGetRcloneVersion_Success(t *testing.T) {
	mockExecutor := func(name string, arg ...string) ([]byte, error) {
		return []byte("rclone v1.60.1-DEV\n- os/version: debian 13.1 (64 bit)\n- os/kernel: 6.12.41+deb13-amd64 (x86_64)\n- os/type: linux\n- os/arch: amd64\n- go/version: go1.23.2\n- go/linking: dynamic\n- go/tags: none"), nil
	}

	version, err := getRcloneVersionWithExecutor("/usr/bin/rclone", mockExecutor)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "rclone v1.60.1-DEV"
	if version != expected {
		t.Errorf("Expected version %q, got %q", expected, version)
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

func TestGetRcloneVersion_MultilineOutput(t *testing.T) {
	mockExecutor := func(name string, arg ...string) ([]byte, error) {
		return []byte("rclone v1.62.2\nline 2\nline 3\n"), nil
	}

	version, err := getRcloneVersionWithExecutor("/usr/bin/rclone", mockExecutor)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "rclone v1.62.2"
	if version != expected {
		t.Errorf("Expected version %q, got %q", expected, version)
	}
}
