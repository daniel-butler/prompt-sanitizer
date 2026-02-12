package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestStdinMode(t *testing.T) {
	testInput := "untrusted data from stdin"
	stdin := strings.NewReader(testInput)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	
	args := []string{"prompt-sanitizer", "--source", "Web Search"}
	
	err := run(args, stdin, stdout, stderr)
	if err != nil {
		t.Fatalf("run() error: %v", err)
	}
	
	output := stdout.String()
	
	if !strings.Contains(output, "<<<EXTERNAL_UNTRUSTED_CONTENT>>>") {
		t.Errorf("Output missing start marker")
	}
	if !strings.Contains(output, "Source: Web Search") {
		t.Errorf("Output missing source label")
	}
	if !strings.Contains(output, testInput) {
		t.Errorf("Output missing input content")
	}
	if !strings.Contains(output, "<<<END_EXTERNAL_UNTRUSTED_CONTENT>>>") {
		t.Errorf("Output missing end marker")
	}
}

func TestFileMode(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	
	testContent := "File content to wrap"
	tmpFile.WriteString(testContent)
	tmpFile.Close()
	
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	
	args := []string{"prompt-sanitizer", "--source", "email", "--file", tmpFile.Name()}
	
	err = run(args, stdin, stdout, stderr)
	if err != nil {
		t.Fatalf("run() error: %v", err)
	}
	
	output := stdout.String()
	
	if !strings.Contains(output, "<<<EXTERNAL_UNTRUSTED_CONTENT>>>") {
		t.Errorf("Output missing start marker")
	}
	if !strings.Contains(output, "Source: email") {
		t.Errorf("Output missing source label")
	}
	if !strings.Contains(output, testContent) {
		t.Errorf("Output missing file content")
	}
	if !strings.Contains(output, "<<<END_EXTERNAL_UNTRUSTED_CONTENT>>>") {
		t.Errorf("Output missing end marker")
	}
}

func TestCommandMode(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	
	args := []string{"prompt-sanitizer", "--source", "curl", "--", "echo", "command output"}
	
	err := run(args, stdin, stdout, stderr)
	if err != nil {
		t.Fatalf("run() error: %v", err)
	}
	
	output := stdout.String()
	
	if !strings.Contains(output, "<<<EXTERNAL_UNTRUSTED_CONTENT>>>") {
		t.Errorf("Output missing start marker")
	}
	if !strings.Contains(output, "Source: curl") {
		t.Errorf("Output missing source label")
	}
	if !strings.Contains(output, "command output") {
		t.Errorf("Output missing command output")
	}
	if !strings.Contains(output, "<<<END_EXTERNAL_UNTRUSTED_CONTENT>>>") {
		t.Errorf("Output missing end marker")
	}
}

func TestEmptyInput(t *testing.T) {
	stdin := strings.NewReader("")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	
	args := []string{"prompt-sanitizer", "--source", "Test"}
	
	err := run(args, stdin, stdout, stderr)
	if err != nil {
		t.Fatalf("run() error: %v", err)
	}
	
	output := stdout.String()
	
	// Should still wrap even empty content
	if !strings.Contains(output, "<<<EXTERNAL_UNTRUSTED_CONTENT>>>") {
		t.Errorf("Output missing start marker")
	}
	if !strings.Contains(output, "<<<END_EXTERNAL_UNTRUSTED_CONTENT>>>") {
		t.Errorf("Output missing end marker")
	}
}

func TestVeryLongContent(t *testing.T) {
	// Generate 1MB of content
	longContent := strings.Repeat("A", 1024*1024)
	stdin := strings.NewReader(longContent)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	
	args := []string{"prompt-sanitizer", "--source", "Large File"}
	
	err := run(args, stdin, stdout, stderr)
	if err != nil {
		t.Fatalf("run() error: %v", err)
	}
	
	output := stdout.String()
	
	if !strings.Contains(output, "<<<EXTERNAL_UNTRUSTED_CONTENT>>>") {
		t.Errorf("Output missing start marker")
	}
	if !strings.Contains(output, longContent) {
		t.Errorf("Output missing long content")
	}
	if !strings.Contains(output, "<<<END_EXTERNAL_UNTRUSTED_CONTENT>>>") {
		t.Errorf("Output missing end marker")
	}
}

func TestBinaryData(t *testing.T) {
	// Binary data with null bytes and special characters
	binaryData := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 'H', 'i', 0x00}
	stdin := bytes.NewReader(binaryData)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	
	args := []string{"prompt-sanitizer", "--source", "Binary"}
	
	err := run(args, stdin, stdout, stderr)
	if err != nil {
		t.Fatalf("run() error: %v", err)
	}
	
	output := stdout.String()
	
	// Should still wrap binary data
	if !strings.Contains(output, "<<<EXTERNAL_UNTRUSTED_CONTENT>>>") {
		t.Errorf("Output missing start marker")
	}
	if !strings.Contains(output, "<<<END_EXTERNAL_UNTRUSTED_CONTENT>>>") {
		t.Errorf("Output missing end marker")
	}
}

func TestNonExistentFile(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	
	args := []string{"prompt-sanitizer", "--source", "test", "--file", "/tmp/nonexistent-file-12345.txt"}
	
	err := run(args, stdin, stdout, stderr)
	if err == nil {
		t.Errorf("Expected error for non-existent file, got nil")
	}
}

func TestFailingCommand(t *testing.T) {
	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	
	args := []string{"prompt-sanitizer", "--source", "test", "--", "false"}
	
	err := run(args, stdin, stdout, stderr)
	if err == nil {
		t.Errorf("Expected error for failing command, got nil")
	}
}
