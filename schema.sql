-- MySQL schema for users table
-- Note: GORM AutoMigrate will create this table automatically based on the User model
-- This schema is provided as a reference for manual database setup

DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    pass_hash VARCHAR(255) NOT NULL,
    phone_num VARCHAR(20),
    role VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Note: Do not insert users with plaintext passwords
-- Use the /register endpoint to create users with properly hashed passwords
-- Example seed data (with bcrypt hashed passwords) can be added manually if needed:
-- Password 'password123' hashed with bcrypt: $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy