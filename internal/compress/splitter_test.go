package compress

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSplitWriter_SingleFile(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.tar.gz")
	
	sw, err := NewSplitWriter(testFile)
	if err != nil {
		t.Fatalf("Failed to create SplitWriter: %v", err)
	}
	defer sw.Close()
	
	// Write less than 2GB
	data := make([]byte, 1024*1024) // 1MB
	for i := 0; i < 100; i++ { // 100MB total
		_, err := sw.Write(data)
		if err != nil {
			t.Fatalf("Failed to write: %v", err)
		}
	}
	
	if err := sw.Close(); err != nil {
		t.Fatalf("Failed to close: %v", err)
	}
	
	// Should create only one file
	parts := sw.Parts()
	if len(parts) != 1 {
		t.Errorf("Expected 1 part, got %d", len(parts))
	}
	
	if !sw.IsSplit() {
		// Good, single file
	}
	
	// Verify file exists
	if _, err := os.Stat(parts[0]); os.IsNotExist(err) {
		t.Errorf("Part file does not exist: %s", parts[0])
	}
}

func TestSplitWriter_MultipleFiles(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.tar.gz")
	
	sw, err := NewSplitWriter(testFile)
	if err != nil {
		t.Fatalf("Failed to create SplitWriter: %v", err)
	}
	defer sw.Close()
	
	// Write more than 2GB to trigger splitting
	// We'll write in 100MB chunks to make it faster
	chunkSize := 100 * 1024 * 1024 // 100MB
	data := make([]byte, chunkSize)
	
	// Write 2.5GB worth of data (25 chunks of 100MB)
	for i := 0; i < 25; i++ {
		_, err := sw.Write(data)
		if err != nil {
			t.Fatalf("Failed to write chunk %d: %v", i, err)
		}
	}
	
	if err := sw.Close(); err != nil {
		t.Fatalf("Failed to close: %v", err)
	}
	
	// Should create multiple files
	parts := sw.Parts()
	if len(parts) < 2 {
		t.Errorf("Expected at least 2 parts, got %d", len(parts))
	}
	
	if !sw.IsSplit() {
		t.Error("Expected IsSplit to be true")
	}
	
	// Verify all files exist and check sizes
	for i, part := range parts {
		info, err := os.Stat(part)
		if os.IsNotExist(err) {
			t.Errorf("Part file %d does not exist: %s", i, part)
			continue
		}
		
		// All parts except the last should be exactly 2GB
		if i < len(parts)-1 {
			if info.Size() != MaxFileSize {
				t.Errorf("Part %d size is %d, expected %d", i, info.Size(), MaxFileSize)
			}
		} else {
			// Last part should be less than 2GB
			if info.Size() >= MaxFileSize {
				t.Errorf("Last part size is %d, should be less than %d", info.Size(), MaxFileSize)
			}
		}
	}
}

func TestSplitWriter_ExactBoundary(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.tar.gz")
	
	sw, err := NewSplitWriter(testFile)
	if err != nil {
		t.Fatalf("Failed to create SplitWriter: %v", err)
	}
	defer sw.Close()
	
	// Write exactly 1.9GB
	chunkSize := 100 * 1024 * 1024 // 100MB
	data := make([]byte, chunkSize)

	// Write exactly 1.9GB (19 chunks of 100MB)
	for i := 0; i < 19; i++ {
		_, err := sw.Write(data)
		if err != nil {
			t.Fatalf("Failed to write chunk %d: %v", i, err)
		}
	}
	
	if err := sw.Close(); err != nil {
		t.Fatalf("Failed to close: %v", err)
	}
	
	// Should create exactly 1 file (2GB is at the boundary)
	parts := sw.Parts()
	if len(parts) != 1 {
		t.Errorf("Expected 1 part for exactly 1.9GB, got %d", len(parts))
	}
}
