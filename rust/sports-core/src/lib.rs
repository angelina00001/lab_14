//! Ядро агрегации: `_impl` для Rust-тестов; PyO3 — с feature `extension-module`.

/// Строка матча для агрегации (чистая Rust-логика).
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct MatchRow {
    pub league: String,
    pub status: String,
    pub home_goals: Option<i32>,
    pub away_goals: Option<i32>,
}

/// Результат агрегации по одной лиге.
#[derive(Debug, Clone, PartialEq)]
pub struct LeagueAgg {
    pub match_count: usize,
    pub goals_sum: i64,
    pub goals_min: i32,
    pub goals_max: i32,
    pub goals_avg: f64,
}

/// Считает total_goals для завершённого матча (impl для тестов).
pub fn total_goals_impl(home: Option<i32>, away: Option<i32>) -> Option<i32> {
    match (home, away) {
        (Some(h), Some(a)) => Some(h + a),
        _ => None,
    }
}

/// Агрегация по лиге: COUNT, SUM, MIN, MAX, AVG (impl).
pub fn aggregate_league_impl(rows: &[MatchRow]) -> LeagueAgg {
    let mut finished_totals: Vec<i32> = Vec::new();
    for row in rows {
        if row.status != "finished" {
            continue;
        }
        if let Some(t) = total_goals_impl(row.home_goals, row.away_goals) {
            finished_totals.push(t);
        }
    }

    let match_count = finished_totals.len();
    if match_count == 0 {
        return LeagueAgg {
            match_count: 0,
            goals_sum: 0,
            goals_min: 0,
            goals_max: 0,
            goals_avg: 0.0,
        };
    }

    let goals_sum: i64 = finished_totals.iter().map(|&g| i64::from(g)).sum();
    let goals_min = *finished_totals.iter().min().unwrap();
    let goals_max = *finished_totals.iter().max().unwrap();
    let goals_avg = goals_sum as f64 / match_count as f64;

    LeagueAgg {
        match_count,
        goals_sum,
        goals_min,
        goals_max,
        goals_avg,
    }
}

#[cfg(feature = "extension-module")]
mod python;

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn total_goals_impl_cases() {
        assert_eq!(total_goals_impl(Some(2), Some(1)), Some(3));
        assert_eq!(total_goals_impl(None, Some(1)), None);
    }

    #[test]
    fn aggregate_league_impl_empty() {
        let agg = aggregate_league_impl(&[]);
        assert_eq!(agg.match_count, 0);
    }

    #[test]
    fn aggregate_league_impl_skips_scheduled() {
        let rows = vec![MatchRow {
            league: "BL".into(),
            status: "scheduled".into(),
            home_goals: None,
            away_goals: None,
        }];
        let agg = aggregate_league_impl(&rows);
        assert_eq!(agg.match_count, 0);
    }

    #[test]
    fn aggregate_league_impl_draw_zero_zero() {
        let rows = vec![MatchRow {
            league: "BL".into(),
            status: "finished".into(),
            home_goals: Some(0),
            away_goals: Some(0),
        }];
        let agg = aggregate_league_impl(&rows);
        assert_eq!(agg.match_count, 1);
        assert_eq!(agg.goals_sum, 0);
    }
}
