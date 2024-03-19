package scan

import (
	"github.com/simonhylander/diskotective/internal/directory"
	"os"
)

type ScannedFile struct {
	Path string
	Size int64
}

type DiskScanEventType string

const (
	InitializedDiskScanEventType DiskScanEventType = "Initialized"
	FileDiskScanEvenType         DiskScanEventType = "File"
	CompletedDiskScanEventType   DiskScanEventType = "Completed"
)

type ScanEvent struct {
	Type        DiskScanEventType
	Items       []directory.Item
	ScannedFile ScannedFile
}

func ListFiles(path string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	return files, nil
}
