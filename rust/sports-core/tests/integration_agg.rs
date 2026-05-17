use sports_stats_core::{aggregate_league_impl, MatchRow};

#[test]
fn integration_aggregate_bundesliga_sample() {
    let rows = vec![
        MatchRow {
            league: "Bundesliga".into(),
            status: "finished".into(),
            home_goals: Some(4),
            away_goals: Some(0),
        },
        MatchRow {
            league: "Bundesliga".into(),
            status: "finished".into(),
            home_goals: Some(2),
            away_goals: Some(0),
        },
        MatchRow {
            league: "Bundesliga".into(),
            status: "scheduled".into(),
            home_goals: None,
            away_goals: None,
        },
    ];

    let agg = aggregate_league_impl(&rows);
    assert_eq!(agg.match_count, 2);
    assert_eq!(agg.goals_sum, 6);
    assert_eq!(agg.goals_min, 2);
    assert_eq!(agg.goals_max, 4);
    assert!((agg.goals_avg - 3.0).abs() < f64::EPSILON);
}
