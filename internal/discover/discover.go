package discover

import (
	"os"
	"path/filepath"
	"strings"
)

// DiscoverFiles streams all .txt file paths from inputDir
func DiscoverFiles(inputDir string) <-chan string {
	out := make(chan string, 100)
	go func() {
		defer close(out)
		filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(path, ".txt") {
				out <- path
			}
			return nil
		})
	}()
	return out
}
