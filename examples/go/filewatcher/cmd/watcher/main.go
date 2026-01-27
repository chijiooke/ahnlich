package main

import (
	"log"

	"go/ahnlich/internals/fw"
	"go/ahnlich/internals/utils"

	"github.com/fsnotify/fsnotify"
)

func main() {
	dir := utils.PromptDirectory()
	absDir := utils.ResolveDirectory(dir)
	utils.ValidateDirectory(absDir)

	// Initial inode snapshot
	inodes := fw.GetINodes(absDir)
	state := make(map[string]fw.FileID, len(inodes))
	for k, v := range inodes {
		state[k] = v
	}

	// Start watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if err := watcher.Add(absDir); err != nil {
		log.Fatal(err)
	}

	log.Println("Watching:", absDir)
	fw.RunEventLoop(watcher, state)
}


// /Users/chijioke/Desktop/file-watcher-pdfs