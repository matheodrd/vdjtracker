package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Watcher struct {
	FilePath    string
	PollingRate time.Duration
}

func NewWatcher(filePath string, pollingRate time.Duration) *Watcher {
	return &Watcher{FilePath: filePath, PollingRate: pollingRate}
}

// Watch watches a file and send its content in a channel when a change is detected.
// Content is space-trimmed.
func (w *Watcher) Watch() <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)

		var lastModTime time.Time
		var lastContent string

		ticker := time.NewTicker(w.PollingRate)
		defer ticker.Stop()

		for range ticker.C {
			info, err := os.Stat(w.FilePath)
			if err != nil {
				continue // file not yet created
			}

			// check if file has been modified
			if !info.ModTime().After(lastModTime) {
				continue
			}
			lastModTime = info.ModTime()

			data, err := os.ReadFile(w.FilePath)
			if err != nil {
				continue
			}

			content := strings.TrimSpace(string(data))

			if content != "" && content != lastContent {
				lastContent = content
				ch <- content
			}
		}
	}()

	return ch
}

func main() {
	watcher := NewWatcher("test.txt", 500*time.Millisecond)

	tracks := watcher.Watch()
	for track := range tracks {
		fmt.Println(track)
	}
}
