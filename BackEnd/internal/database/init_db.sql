SELECT 'CREATE DATABASE house' WHERE NOT EXISTS (SELECT * FROM pg_database WHERE datname = 'house');
-- CREATE DATABASE IF NOT EXISTS house;
\c house;

CREATE TABLE IF NOT EXISTS temperatures (
    id SERIAL PRIMARY KEY,
    data JSONB
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    "username" TEXT NOT NULL ,
    "password" TEXT NOT NULL,
    CONSTRAINT unique_user UNIQUE ("username")
);

\c postgres;