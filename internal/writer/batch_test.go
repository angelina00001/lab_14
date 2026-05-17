package writer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"sports-stats-pipeline/internal/model"
)

func sampleMatch(id int) model.Match {
	g1, g2 := 1, 2
	return model.Match{
		MatchID:   id,
		League:    "Bundesliga",
		HomeTeam:  "A",
		AwayTeam:  "B",
		HomeGoals: &g1,
		AwayGoals: &g2,
		Status:    "finished",
	}
}

func TestBatchWriter_WriteFlush(t *testing.T) {
	tests := []struct {
		name      string
		batchSize int
		writes    int
	}{
		{"single record", 5, 1},
		{"exact batch", 2, 2},
		{"two batches", 3, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "out.jsonl")
			w, err := NewBatchWriter(path, tt.batchSize, 0)
			if err != nil {
				t.Fatalf("NewBatchWriter: %v", err)
			}

			for i := 1; i <= tt.writes; i++ {
				if err := w.Write(sampleMatch(i)); err != nil {
					t.Fatalf("Write: %v", err)
				}
			}

			if err := w.Close(); err != nil {
				t.Fatalf("Close: %v", err)
			}

			lines, err := countJSONLLines(path)
			if err != nil {
				t.Fatalf("count: %v", err)
			}
			if lines != tt.writes {
				t.Errorf("lines=%d want %d", lines, tt.writes)
			}
		})
	}
}

func TestBatchWriter_WriteAfterClose(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.jsonl")
	w, err := NewBatchWriter(path, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	if err := w.Write(sampleMatch(1)); err == nil {
		t.Fatal("expected error after close")
	}
}

func TestBatchWriter_TimerFlush(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.jsonl")
	w, err := NewBatchWriter(path, 100, 50*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.Write(sampleMatch(1)); err != nil {
		t.Fatal(err)
	}
	time.Sleep(120 * time.Millisecond)
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	lines, err := countJSONLLines(path)
	if err != nil {
		t.Fatal(err)
	}
	if lines != 1 {
		t.Errorf("timer flush: lines=%d", lines)
	}
}

func countJSONLLines(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	if len(data) == 0 {
		return 0, nil
	}
	n := 0
	for _, b := range data {
		if b == '\n' {
			n++
		}
	}
	return n, nil
}

func TestBatchWriter_JSONValid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.jsonl")
	w, err := NewBatchWriter(path, 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.Write(sampleMatch(42)); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(data) {
		t.Fatal("invalid json line")
	}
}
