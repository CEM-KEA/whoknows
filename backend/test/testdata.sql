-- Ensure we are operating in the correct schema
SET search_path TO public;

-- Create the "users" table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create the "pages" table
CREATE TABLE IF NOT EXISTS pages (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL UNIQUE,
    language VARCHAR(2) NOT NULL DEFAULT 'en' CHECK (language IN ('en', 'da')),
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create the "jwts" table
CREATE TABLE IF NOT EXISTS jwts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for optimized lookups
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_pages_url ON pages (url);
CREATE INDEX IF NOT EXISTS idx_jwts_user_id ON jwts (user_id);

-- Insert some test data
INSERT INTO users (username, email, password_hash) VALUES
    ('testuser', 'testuser@test.com', '$2a$12$v5MHSoQqaboza8trubEnVu/z6gVKHKGBtrna59OwA1QdBuBYzILa2');
INSERT INTO pages (title, url, language, content) VALUES
    ('Go Programming', '/go-programming', 'en', 'A comprehensive guide to Go programming.'),
    ('Python Programming', '/python-programming', 'en', 'A comprehensive guide to Python programming.'),
    ('Danish Guide', '/danish-guide', 'da', 'Guide to Danish culture and language.');


-- Ensure sequences are in line with existing data (if any)
SELECT setval('users_id_seq', (SELECT COALESCE(MAX(id), 1) FROM users));
SELECT setval('pages_id_seq', (SELECT COALESCE(MAX(id), 1) FROM pages));
SELECT setval('jwts_id_seq', (SELECT COALESCE(MAX(id), 1) FROM jwts));

