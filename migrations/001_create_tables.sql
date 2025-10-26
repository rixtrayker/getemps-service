-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    national_number VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for users table
CREATE INDEX idx_national_number ON users(national_number);
CREATE INDEX idx_is_active ON users(is_active);

-- Create salaries table
CREATE TABLE salaries (
    id SERIAL PRIMARY KEY,
    year INT NOT NULL,
    month INT NOT NULL CHECK (month BETWEEN 1 AND 12),
    salary DECIMAL(10, 2) NOT NULL CHECK (salary >= 0),
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_year_month UNIQUE (user_id, year, month)
);

-- Create indexes for salaries table
CREATE INDEX idx_user_id ON salaries(user_id);
CREATE INDEX idx_year_month ON salaries(year, month);

-- Create logs table (for bonus logging feature)
CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    level VARCHAR(20) NOT NULL,  -- INFO, ERROR, WARN, DEBUG
    message TEXT NOT NULL,
    context JSONB,               -- Additional context data
    endpoint VARCHAR(255),
    request_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for logs table
CREATE INDEX idx_level ON logs(level);
CREATE INDEX idx_created_at ON logs(created_at);
CREATE INDEX idx_endpoint ON logs(endpoint);