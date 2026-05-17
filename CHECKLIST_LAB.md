# Чеклист лабораторной (соответствие проекта)

| Критерий | Статус | Где |
|----------|--------|-----|
| `go.mod` | ✅ | корень |
| `Cargo.toml` | ✅ | `rust/sports-core/Cargo.toml` |
| `pyproject.toml` | ✅ | `python/pyproject.toml` |
| Тесты Go | ✅ | `go test ./...` |
| Тесты Rust | ✅ | `cargo test` (+ `tests/integration_agg.rs`) |
| Тесты Python | ✅ | `pytest` в `python/tests/` |
| Table-driven subtests (Go) | ✅ | `openliga_test.go`, `batch_test.go` |
| `_impl` в Rust | ✅ | `total_goals_impl`, `aggregate_league_impl` |
| `crate-type = ["cdylib", "rlib"]` | ✅ | `Cargo.toml` |
| `edition = "2021"` | ✅ | `Cargo.toml` |
| Автосборка Go в pytest (session) | ✅ | `conftest.py` → `go_collector_exe` |
| PyO3 0.25 + `Bound<>` | ✅ | `python.rs` (`PyDict::new`, `Bound<PyList>`) |
| cProfile | ✅ | `profile_run.py`, `test_profile.py` |
| HTTP httptest (Go) | ✅ | `openliga_fetch_test.go` |
| Graceful shutdown (Go) | ✅ | `internal/collector/run_test.go` |
| Parquet / DuckDB / viz (Python) | ✅ | `test_parquet_duckdb.py`, `test_visualize.py` |
| Rust ↔ Polars | ✅ | `test_rust_polars_consistency.py` |
| Обработка `err` в Go | ✅ | проверки + `slog` в timer flush |
| PROMPT_LOG | ✅ | `PROMPT_LOG.md` (допишите свои промпты) |
| README ФИО/группа/вариант | ⚠️ | `README.md` — **замените шаблон** |
| Conventional commits | ⚠️ | делаете при `git commit` вручную |
| Репозиторий < 300 КБ | ✅ | при соблюдении `.gitignore` |
| Docker Compose | ✅ | опционально, не обязателен для тестов |

## Перед push

```powershell
.\scripts\run-tests.ps1
```

Или по шагам — см. `README.md`.
