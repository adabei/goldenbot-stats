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
);
