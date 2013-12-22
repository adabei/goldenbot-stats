create table games (
  started_at integer primary key,
  ended_at integer,
  mapname text
);

create table stats (
  games_started_at integer,
  players_id text,
  kills integer,
  deaths integer,
  assists integer,
  primary key(games_started_at, players_id)
);
