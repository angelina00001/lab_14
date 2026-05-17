//! PyO3-обёртки (feature `extension-module`, собирается через maturin).

use pyo3::prelude::*;
use pyo3::types::{PyDict, PyList, PyTuple};

use crate::{aggregate_league_impl, total_goals_impl, MatchRow};

#[pyfunction]
fn aggregate_league_py(py: Python<'_>, rows: Bound<'_, PyList>) -> PyResult<PyObject> {
    let mut parsed: Vec<MatchRow> = Vec::with_capacity(rows.len());
    for i in 0..rows.len() {
        let item = rows.get_item(i)?;
        let tuple = item.downcast::<PyTuple>()?;
        if tuple.len() != 4 {
            return Err(pyo3::exceptions::PyValueError::new_err(
                "ожидается кортеж (league, status, home_goals, away_goals)",
            ));
        }
        let league: String = tuple.get_item(0)?.extract()?;
        let status: String = tuple.get_item(1)?.extract()?;
        let home: Option<i32> = tuple.get_item(2)?.extract()?;
        let away: Option<i32> = tuple.get_item(3)?.extract()?;
        parsed.push(MatchRow {
            league,
            status,
            home_goals: home,
            away_goals: away,
        });
    }

    let agg = aggregate_league_impl(&parsed);
    let dict = PyDict::new(py);
    dict.set_item("match_count", agg.match_count)?;
    dict.set_item("goals_sum", agg.goals_sum)?;
    dict.set_item("goals_min", agg.goals_min)?;
    dict.set_item("goals_max", agg.goals_max)?;
    dict.set_item("goals_avg", agg.goals_avg)?;
    Ok(dict.into_any().unbind())
}

#[pyfunction]
fn total_goals_py(home: Option<i32>, away: Option<i32>) -> Option<i32> {
    total_goals_impl(home, away)
}

#[pymodule]
fn sports_stats_core(m: &Bound<'_, PyModule>) -> PyResult<()> {
    m.add_function(wrap_pyfunction!(total_goals_py, m)?)?;
    m.add_function(wrap_pyfunction!(aggregate_league_py, m)?)?;
    Ok(())
}
