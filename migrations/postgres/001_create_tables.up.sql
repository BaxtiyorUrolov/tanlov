
CREATE TABLE IF NOT EXISTS partners (
    id UUID PRIMARY KEY,
    full_name VARCHAR(50),
    phone VARCHAR(15) UNIQUE,
    email VARCHAR(25) UNIQUE,
    video_link VARCHAR(100) UNIQUE,
    score INT DEFAULT 0,
    video_verify BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS admins (
    id BIGINT PRIMARY KEY
)
