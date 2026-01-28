package find

import (
	"context"
	"main/stuff/find/core"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Finder struct {
	sem chan struct{}
}

func NewFinder(workers int) *Finder {
	if workers < 1 {
		workers = 1
	}

	return &Finder{
		sem: make(chan struct{}, workers),
	}
}

func (f *Finder) Run(ctx context.Context, root, target string) *core.Stats {

	statsChan := make(chan *core.StatsCollection, 100)

	var wg sync.WaitGroup

	f.sem <- struct{}{}
	wg.Add(1)
	go func() {
		f.walk(ctx, root, target, &wg, statsChan)
		<-f.sem
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(statsChan)
	}()

	stats := core.NewStats()
	stats.Collect(ctx, statsChan)

	return stats
}

func (f *Finder) walk(ctx context.Context, p, target string, wg *sync.WaitGroup, stats chan<- *core.StatsCollection) {

	local := &core.StatsCollection{}
	defer func(local *core.StatsCollection) {
		select {
		case <-ctx.Done():
		case stats <- local:
		}
	}(local)

	e, err := os.ReadDir(p)
	if err != nil {
		local.Errors++
		return
	}

	for _, entry := range e {
		fullPath := filepath.Join(p, entry.Name())

		if strings.Contains(entry.Name(), target) {
			local.Matches++
		}

		if entry.IsDir() {
			local.Dirs++
			select {
			case <-ctx.Done():
				return
			case f.sem <- struct{}{}:
				wg.Add(1)
				go func() {
					f.walk(ctx, fullPath, target, wg, stats)
					wg.Done()
					<-f.sem
				}()
			default:
				f.walk(ctx, fullPath, target, wg, stats)
			}
		}

		if entry.Type().IsRegular() {
			fi, err := entry.Info()
			if err != nil {
				local.Errors++
				continue
			}

			local.Bytes += fi.Size()
			local.Files++
		}
	}
}
