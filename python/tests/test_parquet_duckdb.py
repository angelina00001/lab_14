"""Parquet round-trip, DuckDB SQL, сравнение Polars vs DuckDB."""
from __future__ import annotations

import time
from pathlib import Path

import polars as pl

from sports_pipeline.analyze import (
    clean_matches,
    duckdb_query,
    load_jsonl,
    polars_query_timing,
)


def test_parquet_round_trip(sample_jsonl: Path, tmp_path: Path) -> None:
    clean = clean_matches(load_jsonl(str(sample_jsonl)))
    path = tmp_path / "matches.parquet"
    clean.write_parquet(path)

    loaded = pl.read_parquet(path)
    assert loaded.height == clean.height
    assert set(loaded.columns) >= {"match_id", "league", "total_goals", "status"}


def test_duckdb_query_on_parquet(sample_jsonl: Path, tmp_path: Path) -> None:
    clean = clean_matches(load_jsonl(str(sample_jsonl)))
    path = tmp_path / "matches.parquet"
    clean.write_parquet(path)

    result = duckdb_query(str(path))
    assert result.height >= 1
    assert "matches" in result.columns
    assert "goals" in result.columns
    finished = clean.filter(pl.col("status") == "finished")
    assert result["matches"].sum() <= finished.height


def test_polars_and_duckdb_same_shape(sample_jsonl: Path, tmp_path: Path) -> None:
    clean = clean_matches(load_jsonl(str(sample_jsonl)))
    path = tmp_path / "matches.parquet"
    clean.write_parquet(path)

    polars_df = polars_query_timing(clean)
    duck_df = duckdb_query(str(path))

    assert polars_df.columns == duck_df.columns
    assert polars_df.height == duck_df.height
    p = polars_df.sort(["league", "matchday"])
    d = duck_df.sort(["league", "matchday"])
    assert p["matches"].to_list() == d["matches"].to_list()
    assert p["goals"].to_list() == d["goals"].to_list()


def test_query_timing_returns_positive(sample_jsonl: Path, tmp_path: Path) -> None:
    clean = clean_matches(load_jsonl(str(sample_jsonl)))
    path = tmp_path / "matches.parquet"
    clean.write_parquet(path)

    t0 = time.perf_counter()
    polars_query_timing(clean)
    t_polars = time.perf_counter() - t0

    t1 = time.perf_counter()
    duckdb_query(str(path))
    t_duck = time.perf_counter() - t1

    assert t_polars >= 0
    assert t_duck >= 0
