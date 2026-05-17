"""cProfile: профилирование загрузки и очистки sample JSONL."""
from __future__ import annotations

import cProfile
import pstats
from pathlib import Path

from sports_pipeline import REPO_ROOT
from sports_pipeline.analyze import clean_matches, load_jsonl


def run_pipeline() -> int:
    sample = REPO_ROOT / "data" / "sample" / "matches_sample.jsonl"
    df = load_jsonl(str(sample))
    clean = clean_matches(df)
    return clean.height


def main() -> None:
    out_dir = REPO_ROOT / "data" / "profile"
    out_dir.mkdir(parents=True, exist_ok=True)
    prof_path = out_dir / "analyze.prof"
    stats_path = out_dir / "analyze_stats.txt"

    profiler = cProfile.Profile()
    profiler.enable()
    rows = run_pipeline()
    profiler.disable()
    profiler.dump_stats(str(prof_path))

    with stats_path.open("w", encoding="utf-8") as fh:
        stats = pstats.Stats(profiler, stream=fh)
        stats.sort_stats("cumulative")
        stats.print_stats(20)

    print(f"cProfile: {prof_path} ({prof_path.stat().st_size} bytes), rows={rows}")
    print(f"Stats: {stats_path}")


if __name__ == "__main__":
    main()
