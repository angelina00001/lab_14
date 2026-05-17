package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchLeagueMatchesAt(t *testing.T) {
	league := LeagueConfig{ShortName: "bl1", Label: "Bundesliga", Season: "2024"}

	tests := []struct {
		name       string
		status     int
		body       string
		wantErr    bool
		errContain string
		wantLen    int
		wantHome   string
	}{
		{
			name:     "HTTP 200 valid JSON",
			status:   http.StatusOK,
			body:     `[{"matchID":1,"matchDateTime":"2024-08-23T18:30:00","matchIsFinished":true,"matchday":1,"team1":{"teamName":"FC A","goals":2},"team2":{"teamName":"FC B","goals":1}}]`,
			wantLen:  1,
			wantHome: "FC A",
		},
		{
			name:       "HTTP 500",
			status:     http.StatusInternalServerError,
			body:       `{"error":"fail"}`,
			wantErr:    true,
			errContain: "HTTP 500",
		},
		{
			name:       "invalid JSON",
			status:     http.StatusOK,
			body:       `{not-json`,
			wantErr:    true,
			errContain: "",
		},
		{
			name:    "empty array",
			status:  http.StatusOK,
			body:    `[]`,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if !strings.Contains(r.URL.Path, "/getmatchdata/bl1/2024") {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			t.Cleanup(srv.Close)

			client := srv.Client()
			got, err := FetchLeagueMatchesAt(context.Background(), client, league, srv.URL)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				if tt.errContain != "" && !strings.Contains(err.Error(), tt.errContain) {
					t.Errorf("err=%v want contain %q", err, tt.errContain)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if len(got) != tt.wantLen {
				t.Fatalf("len=%d want %d", len(got), tt.wantLen)
			}
			if tt.wantLen > 0 && got[0].HomeTeam != tt.wantHome {
				t.Errorf("home=%q want %q", got[0].HomeTeam, tt.wantHome)
			}
			if tt.wantLen > 0 && got[0].Status != "finished" {
				t.Errorf("status=%q", got[0].Status)
			}
		})
	}
}

func TestFetchLeagueMatchesAt_contextCanceled(t *testing.T) {
	league := LeagueConfig{ShortName: "bl1", Label: "BL", Season: "2024"}
	block := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-block
	}))
	t.Cleanup(func() {
		close(block)
		srv.Close()
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := FetchLeagueMatchesAt(ctx, srv.Client(), league, srv.URL)
	if err == nil {
		t.Fatal("expected cancel/connection error")
	}
}
