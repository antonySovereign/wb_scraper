CREATE SCHEMA wb_scraper;

CREATE TABLE IF NOT EXISTS wb_scraper.categories (
    id SERIAL PRIMARY KEY,
    wb_id INTEGER UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    parent_id INTEGER,
    url TEXT,
    shard VARCHAR(50),
    query TEXT
);
