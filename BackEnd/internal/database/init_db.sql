SELECT 'CREATE DATABASE house' WHERE NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'house');
\c house;

CREATE TABLE IF NOT EXISTS temperatures (
    id SERIAL PRIMARY KEY,
    data JSONB
);
