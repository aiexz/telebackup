package compress

import (
	"fmt"
	"log/slog"
	"os"
)

const MaxFileSize = 1900 * 1024 * 1024 // about 1.9GB

// SplitWriter wraps an io.Writer and splits output into multiple files when size exceeds MaxFileSize
type SplitWriter struct {
	baseFilePath string
	currentPart  int
	currentSize  int64
	currentFile  *os.File
	parts        []string // List of created part files
}

// NewSplitWriter creates a new SplitWriter
func NewSplitWriter(baseFilePath string) (*SplitWriter, error) {
	sw := &SplitWriter{
		baseFilePath: baseFilePath,
		currentPart:  0,
		currentSize:  0,
		parts:        make([]string, 0),
	}
	
	// Create the first file
	if err := sw.createNewPart(); err != nil {
		slog.Error("Failed to create initial part", "err", err, "baseFilePath", baseFilePath)
		return nil, err
	}
	
	return sw, nil
}

// createNewPart creates a new part file
func (sw *SplitWriter) createNewPart() error {
	if sw.currentFile != nil {
		if err := sw.currentFile.Close(); err != nil {
			slog.Error("Error closing current file", "err", err, "currentPart", sw.currentPart)
			return err
		}
	}
	
	var partPath string
	if sw.currentPart == 0 {
		partPath = sw.baseFilePath
	} else {
		partPath = fmt.Sprintf("%s.part%d", sw.baseFilePath, sw.currentPart)
	}
	
	slog.Debug("Creating part file", "path", partPath, "partIndex", sw.currentPart)
	file, err := os.OpenFile(partPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		slog.Error("Failed to open part file", "path", partPath, "err", err)
		return err
	}
	
	sw.currentFile = file
	sw.parts = append(sw.parts, partPath)
	sw.currentSize = 0
	
	return nil
}

// Write implements io.Writer interface
func (sw *SplitWriter) Write(p []byte) (n int, err error) {
	totalWritten := 0
	
	for len(p) > 0 {
		// Check if we need to create a new part
		if sw.currentSize >= MaxFileSize {
			slog.Debug("MaxFileSize reached, rotating to new part", "currentPart", sw.currentPart, "currentSize", sw.currentSize)
			sw.currentPart++
			if err := sw.createNewPart(); err != nil {
				slog.Error("Failed to create new part during Write", "err", err, "newPart", sw.currentPart)
				return totalWritten, err
			}
		}
		
		// Calculate how much we can write to current part
		remaining := MaxFileSize - sw.currentSize
		toWrite := int64(len(p))
		if toWrite > remaining {
			toWrite = remaining
		}
		
		// Write to current file
		written, err := sw.currentFile.Write(p[:toWrite])
		totalWritten += written
		sw.currentSize += int64(written)
		
		if err != nil {
			slog.Error("Error writing to part file", "err", err, "part", sw.currentPart, "attempted", toWrite, "written", written)
			return totalWritten, err
		}
		
		p = p[written:]
	}
	
	return totalWritten, nil
}

// Close closes the current file
func (sw *SplitWriter) Close() error {
	if sw.currentFile != nil {
		return sw.currentFile.Close()
	}
	return nil
}

// Parts returns the list of created part files
func (sw *SplitWriter) Parts() []string {
	return sw.parts
}

// IsSplit returns true if the output was split into multiple parts
func (sw *SplitWriter) IsSplit() bool {
	return len(sw.parts) > 1
}
