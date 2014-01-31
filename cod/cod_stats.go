package cod

import (
	"database/sql"
	"github.com/adabei/goldenbot/events"
	"github.com/adabei/goldenbot/events/cod"
	"log"
	"time"
)

type Stats struct {
	events chan interface{}
	db     *sql.DB
}

func NewStats(ea events.Aggregator, db *sql.DB) *Stats {
	s := new(Stats)
	s.events = ea.Subscribe(s)
	s.db = db
	return s
}

const gamesSchema = `
create table games (
  started_at text primary key,
  ended_at text,
  mapname text
);`

const statsSchema = `
create table stats (
  games_started_at text,
  players_id text,
  kills integer,
  deaths integer,
  assists integer,
  primary key(games_started_at, players_id)
);`

func (s *Stats) Setup() error {
	_, err := s.db.Exec(gamesSchema)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = s.db.Exec(statsSchema)
	if err != nil {
		log.Println(err)
	}
	return err
}

type playerStats struct {
	Kills   int
	Deaths  int
	Assists int
}

func (s *Stats) Start() {
	currentStats := make(map[string]*playerStats)
	currentStartedAt := time.Now().Unix()
	var currentMap string
	for {
		ev := <-s.events
		switch ev := ev.(type) {
		case cod.InitGame:
			currentStats = make(map[string]*playerStats)
			currentStartedAt = ev.Unix
			currentMap = "mp_backlot" // TODO extract from initgame
		case cod.ExitLevel:
			if len(currentStats) > 0 {
				// write to db
				for k, v := range currentStats {
					log.Println("stats: inserting game", currentStartedAt, "into database")
					_, err := s.db.Exec("insert into games(started_at, ended_at, mapname) values (?, ?, ?);", currentStartedAt, ev.Unix, currentMap)
					if err != nil {
						log.Fatal("stats: failed to insert games", err)
					}

					log.Println("stats: inserting stats for player", k, "into database")
					_, err = s.db.Exec("insert into stats(games_started_at, players_id, kills, deaths, assists) values(?, ?, ?, ?, ?);", currentStartedAt, k, v.Kills, v.Deaths, v.Assists)

					if err != nil {
						log.Fatal("stats: failed to insert stats for player", k, err)
					}
				}
			}
		case cod.ShutdownGame:
			// shutdowngame vs exitlevel?
		case cod.Kill:
			// TODO suicide
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
