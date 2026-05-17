"""Тесты Rust/PyO3 ядра (Bound API, _impl через Python)."""
from __future__ import annotations


def test_total_goals_py(rust_extension) -> None:
    assert rust_extension.total_goals_py(2, 1) == 3
    assert rust_extension.total_goals_py(None, 1) is None


def test_aggregate_league_py(rust_extension) -> None:
    rows = [
        ("Bundesliga", "finished", 4, 0),
        ("Bundesliga", "finished", 2, 0),
        ("Bundesliga", "scheduled", None, None),
    ]
    agg = rust_extension.aggregate_league_py(rows)
    assert agg["match_count"] == 2
    assert agg["goals_sum"] == 6
    assert agg["goals_min"] == 2
    assert agg["goals_max"] == 4
    assert abs(agg["goals_avg"] - 3.0) < 1e-9
