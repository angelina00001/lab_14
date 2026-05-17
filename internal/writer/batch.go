package writer

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"sports-stats-pipeline/internal/model"
)

// BatchWriter — пакетная запись JSONL (N записей или интервал T).
type BatchWriter struct {
	path      string
	batchSize int
	flushDur  time.Duration

	mu     sync.Mutex
	file   *os.File
	buffer []model.Match
	timer  *time.Timer
	closed bool
}

func NewBatchWriter(path string, batchSize int, flushInterval time.Duration) (*BatchWriter, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	w := &BatchWriter{
		path:      path,
		batchSize: batchSize,
		flushDur:  flushInterval,
		file:      f,
		buffer:    make([]model.Match, 0, batchSize),
	}
	w.resetTimer()
	return w, nil
}

func (w *BatchWriter) resetTimer() {
	if w.flushDur <= 0 {
		return
	}
	if w.timer != nil {
		w.timer.Stop()
	}
	w.timer = time.AfterFunc(w.flushDur, func() {
		if err := w.Flush(); err != nil {
			slog.Error("batch timer flush", "path", w.path, "err", err)
		}
	})
}

// Write добавляет запись в буфер; при достижении batchSize сбрасывает на диск.
func (w *BatchWriter) Write(m model.Match) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return fmt.Errorf("writer closed")
	}
	w.buffer = append(w.buffer, m)
	if len(w.buffer) >= w.batchSize {
		return w.flushLocked()
	}
	return nil
}

// Flush записывает накопленный буфер в файл.
func (w *BatchWriter) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.flushLocked()
}

func (w *BatchWriter) flushLocked() error {
	if len(w.buffer) == 0 {
		w.resetTimer()
		return nil
	}
	enc := json.NewEncoder(w.file)
	for _, m := range w.buffer {
		if err := enc.Encode(m); err != nil {
			return err
		}
	}
	if err := w.file.Sync(); err != nil {
		return err
	}
	w.buffer = w.buffer[:0]
	w.resetTimer()
	return nil
}

// Close сбрасывает буфер и закрывает файл.
func (w *BatchWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return nil
	}
	w.closed = true
	if w.timer != nil {
		w.timer.Stop()
	}
	if err := w.flushLocked(); err != nil {
		_ = w.file.Close()
		return err
	}
	return w.file.Close()
}
