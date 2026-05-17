# Конвейер сбора и анализа спортивной статистики

## Студент

| | |
|---|---|
| **ФИО** | _Иванов Иван Иванович_ |
| **Группа** | _ИВТ-401_ |
| **Вариант** | Сбор и анализ спортивной статистики · HTTP API футбольных лиг (OpenLigaDB) |

> Замените ФИО и группу на свои перед сдачей.

### Версия Python

В задании указано **Python 3.10+**, не обязательно 3.14. **Python 3.14 подходит** с **PyO3 0.25+** (в `rust/sports-core/Cargo.toml`).

PyO3 API (лабораторная, `Bound<>`): `Bound<'_, PyList>`, `Bound<'_, PyModule>`, `PyDict::new(py)` — в 0.25 метод `new_bound` **удалён**, используется `new`.

```powershell
python --version
cd python
py -3.14 -m venv .venv
.\.venv\Scripts\python.exe -m pip install -U pip -e ".[dev]"
.\.venv\Scripts\python.exe -m maturin develop --release --features extension-module
.\.venv\Scripts\python.exe -m pytest -v
```

---

## Предметная область

| Параметр | Значение |
|----------|----------|
| Область | Статистика футбольных матчей (Bundesliga, 2. Bundesliga) |
| Источник | HTTP REST [OpenLigaDB](https://api.openligadb.com) |
| Сбор | Go: горутины, канал, пакетная JSONL, graceful shutdown |
| Анализ | Python Polars; Rust PyO3 (`_impl`); DuckDB; Plotly |
| Профилирование | `cProfile` (`sports_pipeline/profile_run.py`) |

## Архитектура

```
OpenLigaDB → Go collector → JSONL → Polars (Python)
                              ↘ Rust core (PyO3 _impl) — агрегация
                    → Parquet → DuckDB SQL → Plotly
```

## Файлы сборки (чеклист лабораторной)

| Файл | Назначение |
|------|------------|
| `go.mod` | Go-модуль сборщика |
| `rust/sports-core/Cargo.toml` | Rust + PyO3 (`edition = "2021"`, `cdylib` + `rlib`) |
| `python/pyproject.toml` | Python + maturin |

## Тесты (3 языка)

**Windows:** не открывайте `run-tests.ps1` двойным кликом (откроется блокнот). Используйте:

```cmd
scripts\run-tests.cmd
```

Подробно: [docs/TESTING.md](docs/TESTING.md)

```bash
go test ./...
cd rust/sports-core && cargo test --no-default-features
cd python && pip install -e ".[dev]" && maturin develop --release --features extension-module && pytest -v
```

| Язык | Файлы |
|------|-------|
| Go | `openliga_test.go`, `openliga_fetch_test.go` (httptest), `batch_test.go`, `collector/run_test.go` (shutdown) |
| Rust | `lib.rs`, `tests/integration_*.rs` |
| Python | `test_analyze`, `test_parquet_duckdb`, `test_visualize`, `test_rust_*`, `test_profile`, `test_go_collector` |

## Запуск конвейера

```bash
# 1. Сбор
go run ./cmd/collector -out data/raw

# 2. Анализ
cd python && sports-analyze --input "../data/raw/*.jsonl"

# Демо без API
sports-analyze --input "../data/sample/matches_sample.jsonl"

# 3. cProfile
python -m sports_pipeline.profile_run
```

PowerShell: `.\scripts\run-collector.ps1`, `.\scripts\run-analyze.ps1`

## Rust `_impl` + PyO3

- `total_goals_impl` / `aggregate_league_impl` — чистая логика (тесты в Rust)
- `total_goals_py`, `aggregate_league_py` — обёртки PyO3 0.25 (`Bound<'_, PyList>`, `PyDict::new`)

## Чеклист перед `git push`

- [ ] `git clone` → `go test ./...` зелёный
- [ ] `cargo test` в `rust/sports-core` зелёный
- [ ] `pytest` зелёный (после `maturin develop`)
- [ ] Нет в git: `.exe`, `target/`, `.venv`, `__pycache__`, `data/raw/*.jsonl`
- [ ] `PROMPT_LOG.md` заполнен
- [ ] README: ваши ФИО и группа
- [ ] Размер репозитория < 1 МБ (без артефактов)

## Docker (опционально)

`docker compose up -d` — NATS, PostgreSQL для расширений.

## Документация

- [PROMPT_LOG.md](PROMPT_LOG.md) — промпты, инструменты, проблемы
