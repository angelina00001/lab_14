Set-Location $PSScriptRoot\..
go run ./cmd/collector -out data/raw -batch 20 -flush-sec 5 -chan-buf 64
