package main

import (
	"context"
	"go/ahnlich/internals/fw"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func fsnotifyAdapter(
	ctx context.Context,
	dir string,
	out chan<- fw.Event,
) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err := w.Add(dir); err != nil {
		return err
	}

	go func() {
		defer w.Close()

		for {
			select {
			case <-ctx.Done():
				return

			case ev := <-w.Events:
				out <- translate(ev)

			case err := <-w.Errors:
				log.Println("watcher error:", err)
			}
		}
	}()

	return nil
}

func translate(ev fsnotify.Event) fw.Event {
	name := filepath.Base(ev.Name)

	switch {
	case ev.Op&fsnotify.Create != 0:
		return fw.Event{Name: name, Path: ev.Name, Type: fw.Created}
	case ev.Op&fsnotify.Remove != 0:
		return fw.Event{Name: name, Path: ev.Name, Type: fw.Removed}
	case ev.Op&fsnotify.Rename != 0:
		return fw.Event{Name: name, Path: ev.Name, Type: fw.Renamed}
	default:
		return fw.Event{Name: name, Path: ev.Name, Type: fw.Modified}
	}
}
