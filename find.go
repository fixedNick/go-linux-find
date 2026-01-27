package find

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Stats struct {
	d       time.Duration
	errors  int64
	files   int64
	dirs    int64
	bytes   int64
	matches int64
}

func (s *Stats) D() time.Duration { return s.d }

func (s *Stats) String() string {
	return fmt.Sprintf("[%.2fs] Files: %d / Dirs: %d / Errors: %d / Matches: %d / MB: %2.2f",
		s.d.Seconds(),
		s.files,
		s.dirs,
		s.errors,
		s.matches,
		float64(s.bytes)/1024/1024,
	)
}

type Find struct {
	wg         sync.WaitGroup
	maxWorkers int

	sem    chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
}

func NewFind(workers int) *Find {
	if workers < 1 {
		workers = 1
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Find{
		maxWorkers: workers,
		sem:        make(chan struct{}, workers),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// One-time call function.
// In future add restart method, if Find called again - cancel ctx, wait workers to finish current queue and renew channels.
func (f *Find) Find(root, target string) *Stats {

	var stats Stats
	startTime := time.Now()

	f.sem <- struct{}{}
	f.wg.Add(1)

	go func() {
		f.dir(root, target, &stats)
		<-f.sem
		f.wg.Done()
	}()

	f.wg.Wait()
	stats.d = time.Since(startTime)
	fmt.Println(stats.String())
	return &stats
}

func (f *Find) Stop() {
	f.cancel()
}

func (f *Find) dir(p string, target string, stats *Stats) {

	type localStats struct {
		errs    int64
		matches int64
		dirs    int64
		bytes   int64
		files   int64
	}

	local := localStats{}

	defer func(l *localStats) {
		atomic.AddInt64(&stats.errors, l.errs)
		atomic.AddInt64(&stats.matches, l.matches)
		atomic.AddInt64(&stats.dirs, l.dirs)
		atomic.AddInt64(&stats.bytes, l.bytes)
		atomic.AddInt64(&stats.files, l.files)
	}(&local)

	e, err := os.ReadDir(p)
	if err != nil {
		local.errs++
		return
	}

	for _, entry := range e {
		entryFp := filepath.Join(p, entry.Name())
		if strings.Contains(entry.Name(), target) {
			local.matches++
		}

		if entry.IsDir() {
			local.dirs++
			select {
			case <-f.ctx.Done():
				return
			case f.sem <- struct{}{}:
				f.wg.Add(1)
				go func() {
					defer func() {
						<-f.sem
						f.wg.Done()
					}()
					f.dir(entryFp, target, stats)
				}()
			default:
				f.dir(entryFp, target, stats)
			}

		} else if entry.Type().IsRegular() {
			fi, err := entry.Info()
			if err != nil {
				local.errs++
				continue
			}

			local.bytes += fi.Size()
			local.files++
		}
	}
}
