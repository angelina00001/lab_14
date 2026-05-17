package collector

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"sports-stats-pipeline/internal/api"
	"sports-stats-pipeline/internal/model"
)

func TestRun_gracefulShutdownFlushesBuffer(t *testing.T) {
	dir := t.TempDir()
	sigCh := make(chan os.Signal, 1)

	var started atomic.Bool
	fetch := func(ctx context.Context, _ *http.Client, _ api.LeagueConfig) ([]model.Match, error) {
		started.Store(true)
		g1, g2 := 1, 0
		// Блокируемся, пока тест не отменит контекст (имитация долгого API).
		<-ctx.Done()
		return []model.Match{{
			MatchID: 99, League: "Test", HomeTeam: "A", AwayTeam: "B",
			HomeGoals: &g1, AwayGoals: &g2, Status: "finished",
		}}, nil
	}

	done := make(chan struct{})
	var outPath string
	var runErr error
	go func() {
		defer close(done)
		outPath, runErr = Run(context.Background(), Config{
			OutDir:        dir,
			BatchSize:     100,
			FlushInterval: time.Hour,
			ChannelBuf:    4,
			HTTPTimeout:   time.Second,
			Leagues:       []api.LeagueConfig{{ShortName: "t", Label: "Test", Season: "2024"}},
			Fetch:         fetch,
			SignalCh:      sigCh,
		})
	}()

	waitFor(t, &started, 2*time.Second)
	sigCh <- syscall.SIGTERM

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not finish after signal")
	}

	if runErr != nil {
		t.Fatalf("Run: %v", runErr)
	}
	lines, err := countLines(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if lines != 1 {
		t.Errorf("lines after shutdown=%d want 1", lines)
	}
}

func TestRun_writesAllMatchesOnNormalComplete(t *testing.T) {
	dir := t.TempDir()
	sigCh := make(chan os.Signal, 1)

	fetch := func(_ context.Context, _ *http.Client, _ api.LeagueConfig) ([]model.Match, error) {
		g1, g2 := 3, 3
		return []model.Match{
			{MatchID: 1, League: "L", HomeTeam: "H", AwayTeam: "A", HomeGoals: &g1, AwayGoals: &g2, Status: "finished"},
			{MatchID: 2, League: "L", HomeTeam: "X", AwayTeam: "Y", HomeGoals: &g1, AwayGoals: &g2, Status: "finished"},
		}, nil
	}

	outPath, err := Run(context.Background(), Config{
		OutDir:        dir,
		BatchSize:     10,
		FlushInterval: 0,
		ChannelBuf:    8,
		Leagues:       []api.LeagueConfig{{Label: "L"}},
		Fetch:         fetch,
		SignalCh:      sigCh,
	})
	if err != nil {
		t.Fatal(err)
	}
	n, err := countLines(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Errorf("lines=%d want 2", n)
	}
}

func waitFor(t *testing.T, flag *atomic.Bool, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if flag.Load() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("condition not met in time")
}

func countLines(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	if len(data) == 0 {
		return 0, nil
	}
	n := 0
	for _, line := range splitLines(string(data)) {
		if line == "" {
			continue
		}
		var check struct {
			MatchID int `json:"match_id"`
		}
		if err := json.Unmarshal([]byte(line), &check); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

func splitLines(s string) []string {
	return strings.Split(strings.TrimSuffix(s, "\n"), "\n")
}

func TestRun_outputInConfiguredDir(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "raw")
	fetch := func(_ context.Context, _ *http.Client, _ api.LeagueConfig) ([]model.Match, error) {
		return nil, nil
	}
	path, err := Run(context.Background(), Config{
		OutDir:    sub,
		BatchSize: 1,
		Leagues:   []api.LeagueConfig{{Label: "X"}},
		Fetch:     fetch,
		SignalCh:  make(chan os.Signal, 1),
	})
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Dir(path) != sub {
		t.Errorf("dir=%s want %s", filepath.Dir(path), sub)
	}
}
