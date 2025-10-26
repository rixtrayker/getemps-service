-- ============================================
-- GetEmpStatus Sample Data
-- Insert statements for testing and development
-- ============================================

-- ============================================
-- Users Table - Sample Data
-- ============================================
INSERT INTO users (id, username, national_number, email, phone, is_active) VALUES
(1, 'jdoe', 'NAT1001', 'jdoe@example.com', '0791111111', TRUE),
(2, 'asalem', 'NAT1002', 'asalem@example.com', '0792222222', TRUE),
(3, 'rhamdan', 'NAT1003', 'rhamdan@example.com', '0793333333', FALSE),
(4, 'lbarakat', 'NAT1004', 'lbarakat@example.com', '0794444444', TRUE),
(5, 'mfaris', 'NAT1005', 'mfaris@example.com', '0795555555', TRUE),
(6, 'nsaleh', 'NAT1006', 'nsaleh@example.com', '0796666666', FALSE),
(7, 'zobeidat', 'NAT1007', 'zobeidat@example.com', '0797777777', TRUE),
(8, 'ahalaseh', 'NAT1008', 'ahalaseh@example.com', '0798888888', TRUE),
(9, 'tkhalaf', 'NAT1009', 'tkhalaf@example.com', '0799999999', FALSE),
(10, 'sshaheen', 'NAT1010', 'sshaheen@example.com', '0781010101', TRUE),
(11, 'tmart', 'NAT1011', 'tmart@example.com', '0781099101', FALSE),
(12, 'aali', 'NAT1012', 'aali@example.com', '0781088101', TRUE);

-- Reset sequence for users table
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));

-- ============================================
-- Salaries Table - Sample Data
-- ============================================

-- User 1 (jdoe - NAT1001): 5 salary records
-- Expected: RED status (average < 2000 after adjustments)
INSERT INTO salaries (id, year, month, salary, user_id) VALUES
(1, 2025, 1, 1200.00, 1),
(2, 2025, 2, 1300.00, 1),
(3, 2025, 3, 1400.00, 1),
(4, 2025, 5, 1500.00, 1),
(5, 2025, 6, 1600.00, 1);  -- Summer month: -5% deduction

-- User 2 (asalem - NAT1002): 5 salary records
-- Expected: RED status (low salaries)
INSERT INTO salaries (id, year, month, salary, user_id) VALUES
(6, 2025, 1, 900.00, 2),
(7, 2025, 2, 950.00, 2),
(8, 2025, 3, 980.00, 2),
(9, 2025, 4, 1100.00, 2),
(10, 2025, 5, 1150.00, 2);

-- User 3 (rhamdan - NAT1003): 2 salary records
-- Expected: INSUFFICIENT_DATA error (< 3 records)
-- Also: User is not active
INSERT INTO salaries (id, year, month, salary, user_id) VALUES
(11, 2025, 1, 400.00, 3),
(15, 2025, 5, 800.00, 3);

-- User 4 (lbarakat - NAT1004): 5 salary records
-- Expected: GREEN status (average > 2000)
INSERT INTO salaries (id, year, month, salary, user_id) VALUES
(16, 2025, 1, 2000.00, 4),
(17, 2025, 2, 2050.00, 4),
(18, 2025, 3, 2100.00, 4),
(19, 2025, 4, 2200.00, 4),
(20, 2025, 5, 2300.00, 4);

-- User 5 (mfaris - NAT1005): 4 salary records
-- Expected: RED status (low salaries)
INSERT INTO salaries (id, year, month, salary, user_id) VALUES
(21, 2025, 1, 600.00, 5),
(22, 2025, 2, 700.00, 5),
(23, 2025, 3, 750.00, 5),
(25, 2025, 5, 850.00, 5);

-- User 6 (nsaleh - NAT1006): 6 salary records
-- Expected: User is not active error
-- Has December bonus month and good salaries
INSERT INTO salaries (id, year, month, salary, user_id) VALUES
(26, 2025, 11, 1500.00, 6),
(27, 2025, 12, 1550.00, 6),  -- December: +10% bonus
(28, 2025, 1, 1600.00, 6),
(29, 2025, 2, 1650.00, 6),
(30, 2025, 3, 1700.00, 6),
(31, 2025, 4, 2000.00, 6);

-- User 7 (zobeidat - NAT1007): 7 salary records
-- Expected: RED status (average around 1200)
-- Has summer months
INSERT INTO salaries (id, year, month, salary, user_id) VALUES
(32, 2025, 1, 1000.00, 7),
(33, 2025, 2, 1100.00, 7),
(34, 2025, 3, 1150.00, 7),
(35, 2025, 4, 1200.00, 7),
(36, 2025, 5, 1250.00, 7),
(37, 2025, 6, 1350.00, 7),  -- Summer: -5%
(38, 2025, 7, 1500.00, 7);  -- Summer: -5%

-- User 8 (ahalaseh - NAT1008): 6 salary records
-- Expected: GREEN status (high salaries > 2000)
-- Has December bonus, total > 10000 (tax applies)
INSERT INTO salaries (id, year, month, salary, user_id) VALUES
(39, 2025, 10, 2200.00, 8),
(40, 2025, 11, 2300.00, 8),
(41, 2025, 12, 2400.00, 8),  -- December: +10% bonus
(42, 2025, 1, 2500.00, 8),
(43, 2025, 2, 2600.00, 8),
(44, 2025, 3, 2800.00, 8);

-- User 9 (tkhalaf - NAT1009): 5 salary records
-- Expected: User is not active error
-- Has multiple summer months
INSERT INTO salaries (id, year, month, salary, user_id) VALUES
(45, 2025, 1, 1700.00, 9),
(46, 2025, 2, 1750.00, 9),
(47, 2025, 6, 1800.00, 9),  -- Summer: -5%
(48, 2025, 7, 1850.00, 9),  -- Summer: -5%
(49, 2025, 8, 1900.00, 9);  -- Summer: -5%

-- User 10 (sshaheen - NAT1010): 6 salary records
-- Expected: RED status (low to medium salaries)
-- Has summer month
INSERT INTO salaries (id, year, month, salary, user_id) VALUES
(50, 2025, 1, 800.00, 10),
(51, 2025, 2, 850.00, 10),
(52, 2025, 3, 900.00, 10),
(53, 2025, 8, 950.00, 10),   -- Summer: -5%
(54, 2025, 9, 1000.00, 10),
(55, 2025, 10, 1200.00, 10);

-- User 11 (tmart - NAT1011): 0 salary records
-- Expected: INSUFFICIENT_DATA error (no records)
-- Also: User is not active

-- User 12 (aali - NAT1012): 3 salary records (minimum)
-- Expected: RED status (exactly at minimum data requirement)
INSERT INTO salaries (id, year, month, salary, user_id) VALUES
(56, 2025, 1, 1500.00, 12),
(57, 2025, 2, 1550.00, 12),
(58, 2025, 3, 1600.00, 12);

-- Reset sequence for salaries table
SELECT setval('salaries_id_seq', (SELECT MAX(id) FROM salaries));

-- ============================================
-- Sample Logs Data (Optional - for testing)
-- ============================================
INSERT INTO logs (level, message, context, endpoint, request_id, status_code) VALUES
('INFO', 'Application started', '{"version": "1.0.0"}', NULL, NULL, NULL),
('INFO', 'Database connection established', '{"database": "getemps_db"}', NULL, NULL, NULL),
('INFO', 'Employee status retrieved successfully', '{"national_number": "NAT1001"}', '/api/GetEmpStatus', 'req-001', 200),
('WARN', 'Cache miss for employee', '{"national_number": "NAT1002"}', '/api/GetEmpStatus', 'req-002', 200),
('ERROR', 'Invalid national number', '{"national_number": "NAT9999"}', '/api/GetEmpStatus', 'req-003', 404);

-- ============================================
-- Verification Queries
-- ============================================

-- Count users by active status
-- SELECT is_active, COUNT(*) as count FROM users GROUP BY is_active;

-- Count salary records per user
-- SELECT u.national_number, u.username, COUNT(s.id) as salary_count
-- FROM users u
-- LEFT JOIN salaries s ON u.id = s.user_id
-- GROUP BY u.id, u.national_number, u.username
-- ORDER BY salary_count DESC;

-- View users with insufficient data
-- SELECT u.national_number, u.username, u.is_active, COUNT(s.id) as records
-- FROM users u
-- LEFT JOIN salaries s ON u.id = s.user_id
-- GROUP BY u.id, u.national_number, u.username, u.is_active
-- HAVING COUNT(s.id) < 3;

-- Salaries with special months (December and Summer)
-- SELECT u.national_number, s.year, s.month, s.salary,
--        CASE 
--          WHEN s.month = 12 THEN 'December Bonus (+10%)'
--          WHEN s.month IN (6,7,8) THEN 'Summer Deduction (-5%)'
--          ELSE 'Normal'
--        END as adjustment_type
-- FROM salaries s
-- JOIN users u ON s.user_id = u.id
-- WHERE s.month IN (6,7,8,12)
-- ORDER BY u.national_number, s.year, s.month;

-- ============================================
-- Test Scenarios Summary
-- ============================================

/*
TEST CASE SUMMARY:

1. NAT1001 (jdoe) - ACTIVE, 5 records → Expected: RED (average < 2000)
   - Has summer month (June)

2. NAT1002 (asalem) - ACTIVE, 5 records → Expected: RED (low salaries)

3. NAT1003 (rhamdan) - INACTIVE, 2 records → Expected: 406 Error (not active)

4. NAT1004 (lbarakat) - ACTIVE, 5 records → Expected: GREEN (average > 2000)

5. NAT1005 (mfaris) - ACTIVE, 4 records → Expected: RED (low salaries)

6. NAT1006 (nsaleh) - INACTIVE, 6 records → Expected: 406 Error (not active)

7. NAT1007 (zobeidat) - ACTIVE, 7 records → Expected: RED
   - Has summer months (June, July)

8. NAT1008 (ahalaseh) - ACTIVE, 6 records → Expected: GREEN (high salaries)
   - Has December bonus
   - Total > 10000, tax applies

9. NAT1009 (tkhalaf) - INACTIVE, 5 records → Expected: 406 Error (not active)
   - Has summer months

10. NAT1010 (sshaheen) - ACTIVE, 6 records → Expected: RED
    - Has summer month (August)

11. NAT1011 (tmart) - INACTIVE, 0 records → Expected: 422 Error (insufficient data)

12. NAT1012 (aali) - ACTIVE, 3 records → Expected: RED (at minimum threshold)

13. NAT9999 (non-existent) → Expected: 404 Error (invalid national number)

*/
