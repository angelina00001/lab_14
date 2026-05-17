use sports_stats_core::{aggregate_league_impl, total_goals_impl, MatchRow};

#[test]
fn integration_only_scheduled_matches() {
    let rows = vec![
        MatchRow {
            league: "2. Bundesliga".into(),
            status: "scheduled".into(),
            home_goals: None,
            away_goals: None,
        },
        MatchRow {
            league: "2. Bundesliga".into(),
            status: "scheduled".into(),
            home_goals: None,
            away_goals: None,
        },
    ];
    let agg = aggregate_league_impl(&rows);
    assert_eq!(agg.match_count, 0);
    assert_eq!(agg.goals_sum, 0);
}

#[test]
fn integration_total_goals_impl_partial_scores() {
    assert_eq!(total_goals_impl(Some(1), None), None);
    assert_eq!(total_goals_impl(Some(2), Some(2)), Some(4));
}
