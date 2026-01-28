package core

import (
	"context"
	"fmt"
	"time"
)

type StatsCollection struct {
	Errors  int32
	Files   int32
	Dirs    int32
	Bytes   int64
	Matches int32
}

type Stats struct {
	startTime time.Time
	Duration  time.Duration
	Errors    int32
	Files     int32
	Dirs      int32
	Bytes     int64
	Matches   int32
}

func NewStats() *Stats {
	return &Stats{
		startTime: time.Now(),
	}
}

func (s *Stats) Collect(ctx context.Context, in <-chan *StatsCollection) {
	defer func() {
		s.End()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case sc, ok := <-in:
			if !ok {
				return
			}
			s.Bytes += sc.Bytes
			s.Dirs += sc.Dirs
			s.Errors += sc.Errors
			s.Files += sc.Files
			s.Matches += sc.Matches
		}
	}
}

func (s *Stats) String() string {
	return fmt.Sprintf("[%.2fs] Files: %d / Dirs: %d / Errors: %d / Matches: %d / MB: %.2f",
		s.Duration.Seconds(),
		s.Files,
		s.Dirs,
		s.Errors,
		s.Matches,
		float64(s.Bytes)/1024/1024,
	)
}

func (s *Stats) End() { s.Duration = time.Since(s.startTime) }
