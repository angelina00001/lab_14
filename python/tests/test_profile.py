"""cProfile: проверка создания .prof файла."""
from __future__ import annotations

from pathlib import Path

from sports_pipeline import REPO_ROOT
from sports_pipeline.profile_run import run_pipeline


def test_cprofile_run_creates_prof() -> None:
    assert run_pipeline() == 4
    from sports_pipeline.profile_run import main

    main()
    prof = REPO_ROOT / "data" / "profile" / "analyze.prof"
    assert prof.is_file() and prof.stat().st_size > 100
