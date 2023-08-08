package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/SnoozeThis-org/logwait/scanners/common"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/exp/maps"
	"google.golang.org/grpc"
)

type File struct {
	Name         string
	OriginalName string
	HasData      chan struct{}
	Stop         chan struct{}
	useNotify    bool
	file         *os.File
}

var (
	srv   *common.Service
	files []*File
	mut   sync.Mutex
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] <file> [<file> ...]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
	}

	c, err := grpc.Dial(*common.ObserverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to observer at %q: %v", *common.ObserverAddress, err)
	}
	defer c.Close()

	filesToWatch := make(map[string]struct{}, flag.NArg())
	for _, fn := range flag.Args() {
		if _, err := os.Stat(fn); errors.Is(err, os.ErrNotExist) {
			log.Fatalf("failed to stat %q", fn)
		}
		filesToWatch[fn] = struct{}{}
	}

	watcher, err := fsnotify.NewWatcher()
	if err == nil {
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Has(fsnotify.Create) {
						for filename, _ := range filesToWatch {
							if event.Name == filename {
								TailFile(event.Name, true)
							}
						}
					}
					if event.Has(fsnotify.Remove) {
						mut.Lock()
						for i, f := range files {
							if f.Name == event.Name {
								f.Stop <- struct{}{}
								files[i] = files[len(files)-1]
								files = files[:len(files)-1]
							}
						}
						mut.Unlock()
					}
					if event.Has(fsnotify.Rename) {
						mut.Lock()
						for _, f := range files {
							if f.Name == event.Name {
								f.IsRenamed()
							}
						}
						mut.Unlock()
					}
					if event.Has(fsnotify.Write) {
						mut.Lock()
						for _, f := range files {
							if f.Name == event.Name {
								f.HasData <- struct{}{}
							}
						}
						mut.Unlock()
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("error:", err)
				}
			}
		}()
	}

	srv = common.NewService(c)
	srv.StartObserving = func() {
		for filename, _ := range filesToWatch {
			TailFile(filename, watcher != nil)
			if watcher != nil {
				watcher.Add(filepath.Dir(filename))
			}
		}
	}
	srv.StopObserving = func() {
		mut.Lock()
		for _, f := range files {
			f.Stop <- struct{}{}
		}
		files = nil
		mut.Unlock()
		if watcher != nil {
			for filename, _ := range filesToWatch {
				watcher.Remove(filepath.Dir(filename))
			}
		}
	}

	if watcher == nil {
		// Periodically check if file exists
		go func() {
			for {
				time.Sleep(5 * time.Second)
				unwatched := maps.Clone(filesToWatch)
				mut.Lock()
				for _, f := range files {
					if f.Name != "" {
						delete(unwatched, f.Name)
					}
				}
				mut.Unlock()
				for u, _ := range unwatched {
					if _, err := os.Stat(u); err == nil {
						TailFile(u, false)
					}
				}
			}
		}()
	}

	srv.ConnectLoop()
}

func TailFile(filename string, useNotify bool) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	f := &File{
		Name:         filename,
		OriginalName: filename,
		Stop:         make(chan struct{}),
		file:         file,
		useNotify:    useNotify,
	}
	if useNotify {
		f.HasData = make(chan struct{})
	}
	mut.Lock()
	files = append(files, f)
	mut.Unlock()

	go f.Tail()
	return nil
}

func (f *File) Tail() {
	defer f.file.Close()

	reader := bufio.NewReader(f.file)
readLoop:
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				select {
				case <-f.Stop:
					return
				default:
				}

				if curpos, err := f.file.Seek(0, io.SeekCurrent); err == nil {
					if info, err := f.file.Stat(); err == nil {
						if curpos > info.Size() {
							// File is truncated, reset file pointer
							f.file.Seek(0, io.SeekStart)
							continue
						}
					}
				}

				if f.useNotify {
					if _, ok := <-f.HasData; ok {
						continue readLoop
					}
					// The file was renamed, which closed the channel. We need to fall back to sleeping :(
				}
				time.Sleep(time.Second)
				continue
			}
			return
		}

		srv.MatchObservables(func(o common.Observable) bool {
			for field, regexp := range o.Regexps {
				switch field {
				case "message":
					if !regexp.MatchString(line) {
						return false
					}
				default:
					return false
				}
			}
			return true
		})
	}
}

func (f *File) IsRenamed() {
	// We don't know our new name
	f.Name = ""

	// We can no longer use notify because we don't know our new name
	close(f.HasData)
}
