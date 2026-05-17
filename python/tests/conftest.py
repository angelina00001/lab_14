"""Pytest fixtures: автосборка Go и maturin (session scope)."""
from __future__ import annotations

import os
import shutil
import subprocess
import sys
from pathlib import Path

import pytest

REPO_ROOT = Path(__file__).resolve().parents[2]
PYTHON_ROOT = Path(__file__).resolve().parents[1]
BIN_DIR = REPO_ROOT / "bin"


@pytest.fixture(scope="session")
def repo_root() -> Path:
    return REPO_ROOT


@pytest.fixture(scope="session")
def go_collector_exe() -> Path:
    """Собирает Go-сборщик один раз на сессию pytest."""
    BIN_DIR.mkdir(parents=True, exist_ok=True)
    name = "collector.exe" if os.name == "nt" else "collector"
    out = BIN_DIR / name
    subprocess.run(
        ["go", "build", "-o", str(out), "./cmd/collector"],
        cwd=REPO_ROOT,
        check=True,
    )
    assert out.is_file(), "go build не создал бинарник"
    return out


@pytest.fixture(scope="session")
def rust_extension():
    """Импорт PyO3-модуля; при отсутствии — maturin develop."""
    try:
        import sports_stats_core

        return sports_stats_core
    except ImportError:
        maturin = shutil.which("maturin")
        if maturin is None:
            pytest.skip("maturin не найден; pip install maturin")
        subprocess.run(
            [
                maturin,
                "develop",
                "--release",
                "--features",
                "extension-module",
            ],
            cwd=PYTHON_ROOT,
            check=True,
        )
        import sports_stats_core

        return sports_stats_core


@pytest.fixture(scope="session")
def sample_jsonl(repo_root: Path) -> Path:
    path = repo_root / "data" / "sample" / "matches_sample.jsonl"
    assert path.is_file(), f"нет sample: {path}"
    return path
