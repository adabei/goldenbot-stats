package cod

import (
	"database/sql"
	"github.com/adabei/goldenbot/events"
	"github.com/adabei/goldenbot/events/cod"
	"github.com/adabei/goldenbot/rcon/"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type Stats struct {
	events rcon.RCONQuery
	db     sql.DB
}

func NewStats(events events.Aggregator, db sql.DB) *Stats {
	s := new(Stats)
	s.events = events
	s.db = db
	return s
}

const schema = `
create table games (
  started_at text primary key,
  ended_at text,
  mapname text
);

create table stats (
  games_started_at text,
  players_id text,
  kills integer,
  deaths integer,
  assists integer,
  primary key(games_started_at, players_id)
);`

func (s *Stats) Setup() error {
	_, err := s.db.Exec(schema)
}

type playerStats struct {
	Kills   int
	Deaths  int
	Assists int
}

func (s *Stats) Start() {
	currentStats := make(map[string]*playerStats)
	currentStartedAt := 0
	for {
		ev := <-s.events
		switch ev := in.(type) {
		case cod.InitGame:
			currentStats := make(map[string]*playerStats)
			currentStartedAt = time.Now().Unix()
		case cod.ExitLevel:
			if len(currentStats) > 0 {
				// write to db
				for k, v := range currentStats {
					s.db.Exec("insert into stats(games_started_at, players_id, kills, deaths, assists) values((?), (?), (?), (?), (?))",
						currentStartedAt, k, v.Kills, v.Deaths, v.Assists)
				}
			}
		case cod.ShutdownGame:
			// shutdowngame vs exitlevel?
		case cod.Kill:
			if s, ok := currentStats[ev.GUIDA]; ok {
				s.Kills = s.Kills + 1
			} else {
				s = &playerStats{Kills: 1, Deaths: 0, Assists: 0}
				currentStats[ev.GUIDA] = s
			}

			if r, ok := currentStats[ev.GUIDB]; ok {
				r.Deaths = r.Deaths + 1
			} else {
				r := &playerStats{Kills: 0, Deaths: 1, Assists: 0}
				currentStats[ev.GUIDB] = r
			}

		case cod.Damage:
			// not yet implemented (used for assists)
		}
	}
}
