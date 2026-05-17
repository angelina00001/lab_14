# Zapusk testov (Windows)

## Vse komandy po poryadku

Otkroyte **PowerShell** i skopiruyte po odnoy stroke (ili vse srazu).

```powershell
cd C:\Users\Ангелина\.cursor\projects\empty-window\sports-stats-pipeline
```

### Variant A: odin skript (rekomenduetsya)

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\run-tests.ps1
```

Ili:

```cmd
scripts\run-tests.cmd
```

### Variant B: vruchnuyu (3 yazyka)

**1. Go** (iz kornya proekta):

```powershell
cd C:\Users\Ангелина\.cursor\projects\empty-window\sports-stats-pipeline
go test ./... -v -count=1
```

**2. Rust** (bez Python):

```powershell
cd C:\Users\Ангелина\.cursor\projects\empty-window\sports-stats-pipeline\rust\sports-core
cargo test --no-default-features -v
```

**3. Python** (venv + maturin + pytest):

```powershell
cd C:\Users\Ангелина\.cursor\projects\empty-window\sports-stats-pipeline\python

python -m venv .venv

$env:VIRTUAL_ENV = "$PWD\.venv"
$env:PATH = "$PWD\.venv\Scripts;" + $env:PATH

.\.venv\Scripts\python.exe -m pip install -U pip
.\.venv\Scripts\python.exe -m pip install -e ".[dev]"
.\.venv\Scripts\python.exe -m pip install maturin
.\.venv\Scripts\python.exe -m maturin develop --release --features extension-module
.\.venv\Scripts\python.exe -m pytest -v
```

> **Ne nuzhen** `Activate.ps1` — ispolzuem polny put k `python.exe` v `.venv`.

---

## Chastye oshibki

| Oshibka | Prichina | Reshenie |
|---------|----------|----------|
| `.ps1` otkryvaetsya v Bloknote | dvoynoy klik po `.ps1` | `run-tests.cmd` ili `powershell -ExecutionPolicy Bypass -File ...` |
| `TerminatorExpectedAtEndOfString` | kodirovka skripta | obnovlen `run-tests.ps1` (tolko ASCII) |
| `cd rust\sports-core` ne naiden | uzhe v etoy papke | snachala `cd` v **koren** proekta |
| maturin: need virtualenv | zapusk ne iz `python\` | `cd ...\python`, sozdat `.venv` |
| Activate.ps1 zapreschen | ExecutionPolicy | ne aktivirovat venv, sm. komandy vyshe |

---

## Ofline (bez API)

```powershell
cd C:\Users\Ангелина\.cursor\projects\empty-window\sports-stats-pipeline\python
.\.venv\Scripts\python.exe -m pytest -v -k "not writes_jsonl"
```
