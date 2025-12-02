CREATE TABLE daily_analytics (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL UNIQUE,
    date DATE NOT NULL UNIQUE,
    total_links INT NOT NULL DEFAULT 0,
    total_visits INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, date)
);