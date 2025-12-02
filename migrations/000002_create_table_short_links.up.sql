-- Active: 1764581450873@@localhost@5444@mydb
CREATE TABLE short_links (
    id SERIAL PRIMARY KEY,
    user_id INT,
    slug VARCHAR(30) UNIQUE NOT NULL,
    url TEXT NOT NULL,
    clicks INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
