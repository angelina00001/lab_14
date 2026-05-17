"""Тесты Polars-конвейера."""
from __future__ import annotations

from pathlib import Path

import polars as pl

from sports_pipeline.analyze import aggregate_by_league, clean_matches, load_jsonl


def test_load_and_clean(sample_jsonl: Path) -> None:
    raw = load_jsonl(str(sample_jsonl))
    assert raw.height == 5
    clean = clean_matches(raw)
    assert clean.height == 4  # один дубликат match_id=1
    assert "total_goals" in clean.columns


def test_aggregate_by_league(sample_jsonl: Path) -> None:
    clean = clean_matches(load_jsonl(str(sample_jsonl)))
    agg = aggregate_by_league(clean)
    assert isinstance(agg, pl.DataFrame)
    assert "goals_sum" in agg.columns
    assert agg.height >= 1
