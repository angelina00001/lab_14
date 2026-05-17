"""Тесты генерации графиков Plotly."""
from __future__ import annotations

from pathlib import Path

from sports_pipeline.analyze import clean_matches, load_jsonl, visualize


def test_visualize_creates_html_files(sample_jsonl: Path, tmp_path: Path) -> None:
    clean = clean_matches(load_jsonl(str(sample_jsonl)))
    out = tmp_path / "charts"
    visualize(clean, out)

    ts = out / "timeseries_avg_goals.html"
    hist = out / "histogram_total_goals.html"
    assert ts.is_file() and ts.stat().st_size > 200
    assert hist.is_file() and hist.stat().st_size > 200
    body = ts.read_bytes().lower()
    assert b"plotly" in body or b"<html" in body
