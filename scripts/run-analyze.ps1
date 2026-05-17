Set-Location $PSScriptRoot\..\python
python -m sports_pipeline.analyze --input "../data/raw/*.jsonl"
