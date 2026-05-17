"""
Конвейер анализа: JSONL → Polars → Parquet → DuckDB → визуализация.
"""
from __future__ import annotations

import argparse
import glob
import time
from pathlib import Path

import duckdb
import polars as pl
import plotly.express as px

from sports_pipeline import REPO_ROOT


def load_jsonl(pattern: str) -> pl.DataFrame:
    paths = sorted(glob.glob(pattern))
    if not paths:
        raise FileNotFoundError(f"Нет файлов по шаблону: {pattern}")
    return pl.concat([pl.read_ndjson(p) for p in paths], how="diagonal_relaxed")


def print_overview(df: pl.DataFrame, title: str) -> None:
    print(f"\n{'=' * 60}\n{title}\n{'=' * 60}")
    print(df.head(5))
    print(f"\nСтрок: {df.height}, столбцов: {df.width}")
    print("\nСхема:", df.schema)


def clean_matches(df: pl.DataFrame) -> pl.DataFrame:
    df = df.unique(subset=["match_id"], keep="first")
    df = df.with_columns(
        pl.col("match_id").cast(pl.Int64),
        pl.col("matchday").cast(pl.Int32),
        pl.col("match_date").str.to_datetime(time_zone="UTC", strict=False),
        pl.col("home_goals").cast(pl.Int32),
        pl.col("away_goals").cast(pl.Int32),
        pl.col("league").cast(pl.Utf8),
        pl.col("status").cast(pl.Utf8),
    )
    df = df.with_columns(
        pl.when(pl.col("status") == "finished")
        .then(pl.col("home_goals").fill_null(0))
        .otherwise(pl.col("home_goals"))
        .alias("home_goals"),
        pl.when(pl.col("status") == "finished")
        .then(pl.col("away_goals").fill_null(0))
        .otherwise(pl.col("away_goals"))
        .alias("away_goals"),
    )
    df = df.filter(
        pl.col("home_team").is_not_null() & pl.col("away_team").is_not_null()
    )
    return df.with_columns(
        (pl.col("home_goals") + pl.col("away_goals")).alias("total_goals"),
        pl.col("match_date").dt.date().alias("match_day"),
    )


def aggregate_by_league(df: pl.DataFrame) -> pl.DataFrame:
    return df.group_by("league").agg(
        pl.len().alias("match_count"),
        pl.col("total_goals").sum().alias("goals_sum"),
        pl.col("total_goals").mean().alias("goals_avg"),
        pl.col("total_goals").min().alias("goals_min"),
        pl.col("total_goals").max().alias("goals_max"),
    ).sort("league")


def polars_query_timing(df: pl.DataFrame) -> pl.DataFrame:
    return (
        df.filter(pl.col("status") == "finished")
        .group_by("league", "matchday")
        .agg(pl.len().alias("matches"), pl.col("total_goals").sum().alias("goals"))
        .sort(["league", "matchday"])
    )


def duckdb_query(parquet_path: str) -> pl.DataFrame:
    path = parquet_path.replace("\\", "/")
    con = duckdb.connect()
    sql = f"""
    SELECT league, matchday, COUNT(*) AS matches, SUM(total_goals) AS goals
    FROM read_parquet('{path}')
    WHERE status = 'finished'
    GROUP BY league, matchday
    ORDER BY league, matchday
    """
    return con.execute(sql).pl()


def visualize(df: pl.DataFrame, out_dir: Path) -> None:
    out_dir.mkdir(parents=True, exist_ok=True)
    finished = df.filter(pl.col("status") == "finished")
    daily = (
        finished.group_by(["match_day", "league"])
        .agg(pl.col("total_goals").mean().alias("avg_goals"))
        .sort("match_day")
    )
    fig1 = px.line(
        daily.to_pandas(),
        x="match_day",
        y="avg_goals",
        color="league",
        title="Среднее число голов по дням",
    )
    fig1.write_html(str(out_dir / "timeseries_avg_goals.html"))
    fig2 = px.histogram(
        finished.select("total_goals", "league").to_pandas(),
        x="total_goals",
        color="league",
        title="Распределение total_goals",
    )
    fig2.write_html(str(out_dir / "histogram_total_goals.html"))


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--input", default=str(REPO_ROOT / "data/raw/*.jsonl"))
    parser.add_argument("--parquet", default=str(REPO_ROOT / "data/processed/matches.parquet"))
    parser.add_argument("--charts", default=str(REPO_ROOT / "data/charts"))
    args = parser.parse_args()

    raw = load_jsonl(args.input)
    print_overview(raw, "Исходные данные")
    clean = clean_matches(raw)
    aggregate_by_league(clean)

    Path(args.parquet).parent.mkdir(parents=True, exist_ok=True)
    clean.write_parquet(args.parquet)

    t0 = time.perf_counter()
    polars_query_timing(clean)
    t_polars = time.perf_counter() - t0
    t1 = time.perf_counter()
    duckdb_query(args.parquet)
    t_duck = time.perf_counter() - t1
    print(f"Polars: {t_polars*1000:.2f} ms | DuckDB: {t_duck*1000:.2f} ms")

    visualize(clean, Path(args.charts))


if __name__ == "__main__":
    main()
