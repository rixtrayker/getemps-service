-- ============================================
-- GetEmpStatus Database Schema
-- PostgreSQL 15+
-- ============================================

-- Drop existing tables if they exist
DROP TABLE IF EXISTS logs CASCADE;
DROP TABLE IF EXISTS salaries CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- ============================================
-- Users Table
-- ============================================
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    national_number VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for users table
CREATE INDEX idx_users_national_number ON users(national_number);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_email ON users(email);

-- Comments for users table
COMMENT ON TABLE users IS 'Stores employee/user information';
COMMENT ON COLUMN users.id IS 'Primary key - auto incrementing';
COMMENT ON COLUMN users.national_number IS 'Unique national identification number';
COMMENT ON COLUMN users.is_active IS 'User active status - only active users can get status';

-- ============================================
-- Salaries Table
-- ============================================
CREATE TABLE salaries (
    id SERIAL PRIMARY KEY,
    year INT NOT NULL CHECK (year >= 2000 AND year <= 2100),
    month INT NOT NULL CHECK (month BETWEEN 1 AND 12),
    salary DECIMAL(10, 2) NOT NULL CHECK (salary >= 0),
    user_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_year_month UNIQUE (user_id, year, month)
);

-- Indexes for salaries table
CREATE INDEX idx_salaries_user_id ON salaries(user_id);
CREATE INDEX idx_salaries_year_month ON salaries(year, month);
CREATE INDEX idx_salaries_user_year ON salaries(user_id, year);

-- Comments for salaries table
COMMENT ON TABLE salaries IS 'Stores monthly salary records for users';
COMMENT ON COLUMN salaries.year IS 'Salary year (2000-2100)';
COMMENT ON COLUMN salaries.month IS 'Salary month (1-12)';
COMMENT ON COLUMN salaries.salary IS 'Salary amount - must be non-negative';
COMMENT ON CONSTRAINT unique_user_year_month ON salaries IS 'Ensures one salary record per user per month';

-- ============================================
-- Logs Table (Bonus Feature)
-- ============================================
CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    level VARCHAR(20) NOT NULL CHECK (level IN ('DEBUG', 'INFO', 'WARN', 'ERROR')),
    message TEXT NOT NULL,
    context JSONB,
    endpoint VARCHAR(255),
    request_id VARCHAR(100),
    user_id INT,
    duration_ms INT,
    status_code INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Indexes for logs table
CREATE INDEX idx_logs_level ON logs(level);
CREATE INDEX idx_logs_created_at ON logs(created_at DESC);
CREATE INDEX idx_logs_endpoint ON logs(endpoint);
CREATE INDEX idx_logs_request_id ON logs(request_id);
CREATE INDEX idx_logs_user_id ON logs(user_id);

-- Comments for logs table
COMMENT ON TABLE logs IS 'Stores application logs for monitoring and debugging';
COMMENT ON COLUMN logs.level IS 'Log level: DEBUG, INFO, WARN, ERROR';
COMMENT ON COLUMN logs.context IS 'Additional context data in JSON format';
COMMENT ON COLUMN logs.duration_ms IS 'Request duration in milliseconds';

-- ============================================
-- Useful Views (Optional)
-- ============================================

-- View for active users with salary count
CREATE OR REPLACE VIEW v_users_salary_summary AS
SELECT 
    u.id,
    u.username,
    u.national_number,
    u.email,
    u.is_active,
    COUNT(s.id) as salary_records_count,
    COALESCE(MIN(s.year), 0) as first_salary_year,
    COALESCE(MAX(s.year), 0) as last_salary_year,
    COALESCE(AVG(s.salary), 0) as raw_average_salary
FROM users u
LEFT JOIN salaries s ON u.id = s.user_id
GROUP BY u.id, u.username, u.national_number, u.email, u.is_active;

COMMENT ON VIEW v_users_salary_summary IS 'Summary view of users with their salary record counts';

-- ============================================
-- Useful Functions (Optional)
-- ============================================

-- Function to get user salary count
CREATE OR REPLACE FUNCTION get_user_salary_count(p_user_id INT)
RETURNS INT AS $$
    SELECT COUNT(*)::INT FROM salaries WHERE user_id = p_user_id;
$$ LANGUAGE SQL STABLE;

COMMENT ON FUNCTION get_user_salary_count(INT) IS 'Returns the number of salary records for a given user';

-- Function to check if user has sufficient data
CREATE OR REPLACE FUNCTION has_sufficient_salary_data(p_user_id INT)
RETURNS BOOLEAN AS $$
    SELECT COUNT(*) >= 3 FROM salaries WHERE user_id = p_user_id;
$$ LANGUAGE SQL STABLE;

COMMENT ON FUNCTION has_sufficient_salary_data(INT) IS 'Returns true if user has at least 3 salary records';

-- ============================================
-- Triggers for updated_at
-- ============================================

-- Function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for users table
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- Grant Permissions (if needed)
-- ============================================
-- Uncomment and modify if you need specific user permissions
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO getemps_app;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO getemps_app;

-- ============================================
-- Verification Queries
-- ============================================
-- Run these to verify schema creation:
-- SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';
-- SELECT * FROM v_users_salary_summary;
-- SELECT indexname FROM pg_indexes WHERE schemaname = 'public';
