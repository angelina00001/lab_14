package collector

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"sports-stats-pipeline/internal/api"
	"sports-stats-pipeline/internal/model"
	"sports-stats-pipeline/internal/writer"
)

// FetchFunc загружает матчи одной лиги (подменяется в тестах).
type FetchFunc func(ctx context.Context, client *http.Client, league api.LeagueConfig) ([]model.Match, error)

// Config параметры запуска сборщика.
type Config struct {
	OutDir        string
	BatchSize     int
	FlushInterval time.Duration
	ChannelBuf    int
	HTTPTimeout   time.Duration
	Leagues       []api.LeagueConfig
	Fetch         FetchFunc
	Client        *http.Client
	SignalCh      chan os.Signal
}

// Run запускает сбор: горутины → канал → пакетная JSONL, graceful shutdown по сигналу.
func Run(ctx context.Context, cfg Config) (outPath string, err error) {
	if cfg.Fetch == nil {
		cfg.Fetch = api.FetchLeagueMatches
	}
	if cfg.Client == nil {
		cfg.Client = &http.Client{Timeout: cfg.HTTPTimeout}
	}
	if cfg.SignalCh == nil {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		cfg.SignalCh = sigCh
	}
	if len(cfg.Leagues) == 0 {
		cfg.Leagues = api.DefaultLeagues
	}
	if err := os.MkdirAll(cfg.OutDir, 0o755); err != nil {
		return "", err
	}

	outPath = filepath.Join(cfg.OutDir, fmt.Sprintf("matches_%s.jsonl", time.Now().UTC().Format("20060102_150405")))
	bw, err := writer.NewBatchWriter(outPath, cfg.BatchSize, cfg.FlushInterval)
	if err != nil {
		return "", err
	}

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	records := make(chan model.Match, cfg.ChannelBuf)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for m := range records {
			if err := bw.Write(m); err != nil {
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var fetchWg sync.WaitGroup
		for _, lg := range cfg.Leagues {
			select {
			case <-runCtx.Done():
				return
			default:
			}
			fetchWg.Add(1)
			go func(league api.LeagueConfig) {
				defer fetchWg.Done()
				matches, fetchErr := cfg.Fetch(runCtx, cfg.Client, league)
				if fetchErr != nil {
					return
				}
				for _, m := range matches {
					if !sendMatch(runCtx, records, m) {
						return
					}
				}
			}(lg)
		}
		fetchWg.Wait()
		close(records)
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-cfg.SignalCh:
		cancel()
		<-done
	case <-runCtx.Done():
		cancel()
		<-done
	}

	if err := bw.Flush(); err != nil {
		_ = bw.Close()
		return outPath, err
	}
	if err := bw.Close(); err != nil {
		return outPath, err
	}
	return outPath, nil
}

// sendMatch передаёт матч в канал; при отмене контекста — best-effort доставка.
func sendMatch(ctx context.Context, records chan<- model.Match, m model.Match) bool {
	select {
	case records <- m:
		return true
	case <-ctx.Done():
		select {
		case records <- m:
			return true
		default:
			return false
		}
	}
}
