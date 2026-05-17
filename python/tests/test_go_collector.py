"""Тесты Go-сборщика (бинарник собирается fixture-ом)."""
from __future__ import annotations

import json
import subprocess
from pathlib import Path

import pytest


def test_go_collector_builds(go_collector_exe: Path) -> None:
    assert go_collector_exe.stat().st_size > 0


def test_go_collector_writes_jsonl(go_collector_exe: Path, repo_root: Path, tmp_path: Path) -> None:
    out = tmp_path / "raw"
    proc = subprocess.run(
        [str(go_collector_exe), "-out", str(out), "-batch", "5", "-flush-sec", "1"],
        cwd=repo_root,
        capture_output=True,
        text=True,
        timeout=120,
    )
    if proc.returncode != 0 and "fetch" in (proc.stderr + proc.stdout).lower():
        pytest_skip_network(proc)
    assert proc.returncode == 0, proc.stderr

    files = list(out.glob("*.jsonl"))
    assert files, "ожидался JSONL"
    lines = files[0].read_text(encoding="utf-8").strip().splitlines()
    assert lines
    obj = json.loads(lines[0])
    assert "match_id" in obj
    assert "league" in obj


def pytest_skip_network(proc: subprocess.CompletedProcess) -> None:
    pytest.skip(f"API недоступен: {proc.stderr[:200]}")
