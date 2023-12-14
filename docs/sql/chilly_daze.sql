CREATE SCHEMA IF NOT EXISTS chilly_daze;

CREATE TABLE IF NOT EXISTS chilly_daze.achievement_categories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  display_name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS chilly_daze.achievements (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  display_name TEXT NOT NULL,
  description TEXT NOT NULL,
  category_id UUID NOT NULL REFERENCES chilly_daze.achievement_categories(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chilly_daze.users (
  id VARCHAR(256) PRIMARY KEY,
  name TEXT NOT NULL,
  avatar UUID REFERENCES chilly_daze.achievements(id) ON UPDATE CASCADE ON DELETE CASCADE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS chilly_daze.user_achievements (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id VARCHAR(256) NOT NULL REFERENCES chilly_daze.users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  achievement_id UUID NOT NULL REFERENCES chilly_daze.achievements(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chilly_daze.chills (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  ended_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS chilly_daze.user_chills (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id VARCHAR(256) NOT NULL REFERENCES chilly_daze.users(id) ON UPDATE CASCADE ON DELETE CASCADE,
  chill_id UUID NOT NULL REFERENCES chilly_daze.chills(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chilly_daze.trace_points (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  chill_id UUID NOT NULL REFERENCES chilly_daze.chills(id) ON UPDATE CASCADE ON DELETE CASCADE,
  latitude DOUBLE PRECISION NOT NULL,
  longitude DOUBLE PRECISION NOT NULL,
  timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS chilly_daze.photos (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  chill_id UUID NOT NULL REFERENCES chilly_daze.chills(id) ON UPDATE CASCADE ON DELETE CASCADE,
  timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  url TEXT NOT NULL
);