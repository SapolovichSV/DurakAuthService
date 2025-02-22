CREATE TYPE role AS ENUM ('admin','user');
CREATE TYPE status AS ENUM ('offline','online');
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) NOT NULL,
    passwordHASH VARCHAR(255) NOT NULL,
    status status default 'offline',
    current_role role
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  
);