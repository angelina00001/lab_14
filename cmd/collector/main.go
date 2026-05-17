package main

import (
	"context"
	"flag"
	"log"
	"time"

	"sports-stats-pipeline/internal/collector"
)

func main() {
	outDir := flag.String("out", "data/raw", "каталог для JSONL")
	batchSize := flag.Int("batch", 20, "размер пакета перед записью")
	flushSec := flag.Float64("flush-sec", 5, "макс. интервал (сек) перед сбросом буфера")
	channelBuf := flag.Int("chan-buf", 64, "ёмкость канала записей")
	timeoutSec := flag.Float64("timeout", 30, "HTTP timeout (сек)")
	flag.Parse()

	outPath, err := collector.Run(context.Background(), collector.Config{
		OutDir:        *outDir,
		BatchSize:     *batchSize,
		FlushInterval: time.Duration(*flushSec * float64(time.Second)),
		ChannelBuf:    *channelBuf,
		HTTPTimeout:   time.Duration(*timeoutSec * float64(time.Second)),
	})
	if err != nil {
		log.Fatalf("collector: %v", err)
	}
	log.Printf("saved to %s", outPath)
}
