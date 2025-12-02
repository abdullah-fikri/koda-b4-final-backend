-- Active: 1764581450873@@localhost@5444@mydb

CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    email varchar(100) UNIQUE NOT NULL,
    password TEXT,
    role TEXT DEFAULT 'user',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE profile (
    id SERIAL PRIMARY KEY,
    users_id BIGINT,
    username varchar(100),
    phone VARCHAR(20),
    address VARCHAR(100),
    profile_picture TEXT,
    FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);