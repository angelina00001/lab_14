# PROMPT_LOG — журнал промптов

## Студент

| Поле | Значение |
|------|----------|
| **ФИО** | Авдошина Ангелина Евгеньевна |
| **Группа** | 221331 |
| **Вариант** | 18. Сбор и анализ спортивной статистики (HTTP API футбольных лиг) |

---

## Промпт 1 — создание конвейера

**Инструмент:** Cursor Agent (Composer)

**Текст промпта (суть):**
> Предметная область: сбор и анализ спортивной статистики. Источник: API футбольных лиг.
> Реализовать Go-сборщик (горутины, JSONL, буфер, graceful shutdown), Python Polars/DuckDB/Plotly, README.

**Результат:**
- Go-сборщик OpenLigaDB, пакетная запись, SIGINT/SIGTERM
- `python/analyze.py` с Polars → Parquet → DuckDB → графики
- `README.md`, `docker-compose.yml`

**Проблемы:**
- Shell в среде агента не всегда выполнял `go run` (пустой вывод) — проверка перенесена на локальный `go test` / `pytest`
- Для сдачи без сети добавлен `data/sample/matches_sample.jsonl`

---

## Промпт 2 — критерии лабораторной

**Инструмент:** Cursor Agent

**Текст промпта (суть):**
> Адаптировать под чеклист: go.mod + Cargo.toml + pyproject.toml, тесты Go/Rust/Python,
> table-driven subtests, `_impl` в Rust, PyO3 0.23 Bound<>, pytest fixture автосборки Go,
> cProfile, PROMPT_LOG, README с ФИО, репозиторий < 300 КБ.

**Результат:**
- `internal/api/openliga_test.go`, `internal/writer/batch_test.go` — table-driven
- `rust/sports-core/` — `total_goals_impl`, `aggregate_league_impl`, PyO3 0.23, `crate-type = ["cdylib", "rlib"]`
- `rust/sports-core/tests/integration_agg.rs`
- `python/pyproject.toml` (maturin), `python/tests/conftest.py` (session: go build + maturin)
- `sports_pipeline/profile_run.py` — cProfile
- Обновлён `.gitignore` (bin, target, venv, артефакты)

**Проблемы:**
- PyO3-модуль требует `maturin develop` перед pytest (автоматизировано в `conftest.py`)
- Интеграционный тест Go-сборщика с API может `skip` при отсутствии сети
- Python 3.14: PyO3 0.23 не поддерживает 3.14 → обновлено до PyO3 0.25

---

## Промпт 3 — Python 3.14, ошибки тестов, PyO3 API

**Инструмент:** Cursor Agent

**Текст промпта (суть):**
> Работаем с Python 3.14? Ошибка `PyDict::new_bound` not found при `pip install -e ".[dev]"`.
> `pytest` не найден. Обновить PROMPT_LOG и README.

**Результат:**
- PyO3 **0.25** (поддержка Python 3.14)
- В `rust/sports-core/src/python.rs`: `PyDict::new(py)` вместо устаревшего `new_bound` (в 0.25 API переименован)
- Возврат в Python: `dict.into_any().unbind()`
- `Bound<'_, PyList>` / `Bound<'_, PyModule>` — требование лабораторной по Bound API сохранено

**Проблемы:**
- `pip install -e ".[dev]"` падал → Rust не собирался → **pytest не устанавливался** (это следствие, не отдельная ошибка)
- `run-tests.ps1` на Windows: кириллица в строках ломала парсер → скрипт переведён на ASCII
- `maturin develop` нужно запускать из `python\` с активным `.venv` (`VIRTUAL_ENV` или `.\.venv\Scripts\python.exe`)
- `Activate.ps1` может быть заблокирован ExecutionPolicy → использовать прямой путь к `python.exe`

---

## Команды проверки перед push

```powershell
cd sports-stats-pipeline
go test ./...
cd rust\sports-core
cargo test --no-default-features
cd ..\python
py -3.14 -m venv .venv
$env:VIRTUAL_ENV = "$PWD\.venv"
$env:PATH = "$PWD\.venv\Scripts;" + $env:PATH
.\.venv\Scripts\python.exe -m pip install -U pip
.\.venv\Scripts\python.exe -m pip install -e ".[dev]"
.\.venv\Scripts\python.exe -m maturin develop --release --features extension-module
.\.venv\Scripts\python.exe -m pytest -v
```
