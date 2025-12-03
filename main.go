package main

import (
	"fmt"
	"time"

	"github.com/matheodrd/vdjtracker/internal/file"
)

func main() {
	watcher := file.NewWatcher("test.txt", 500*time.Millisecond)

	tracks := watcher.Watch()
	for track := range tracks {
		fmt.Println(track)
	}
}
