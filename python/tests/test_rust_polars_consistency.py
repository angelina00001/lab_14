"""Согласованность агрегации Rust (PyO3) и Polars."""
from __future__ import annotations

from pathlib import Path

from sports_pipeline.analyze import clean_matches, load_jsonl


def test_rust_aggregate_matches_polars_finished_goals(
    rust_extension, sample_jsonl: Path
) -> None:
    clean = clean_matches(load_jsonl(str(sample_jsonl)))
    finished = clean.filter(
        (clean["status"] == "finished")
        & clean["home_goals"].is_not_null()
        & clean["away_goals"].is_not_null()
    )

    rows = [
        (row["league"], row["status"], row["home_goals"], row["away_goals"])
        for row in finished.iter_rows(named=True)
    ]
    rust_agg = rust_extension.aggregate_league_py(rows)

    polars_sum = int(finished["total_goals"].sum())
    assert rust_agg["goals_sum"] == polars_sum
    assert rust_agg["match_count"] == finished.height
