-- Insert sample users
INSERT INTO users (username, national_number, email, phone, is_active) VALUES
('jdoe', 'NAT1001', 'jdoe@example.com', '0791111111', true),
('asmith', 'NAT1002', 'asmith@example.com', '0792222222', true),
('mjohnson', 'NAT1003', 'mjohnson@example.com', '0793333333', false), -- inactive user
('kwilson', 'NAT1004', 'kwilson@example.com', '0794444444', true),
('lbrown', 'NAT1005', 'lbrown@example.com', '0795555555', true),
('tjones', 'NAT1006', 'tjones@example.com', '0796666666', true),
('dgarcia', 'NAT1007', 'dgarcia@example.com', '0797777777', true),
('rmiller', 'NAT1008', 'rmiller@example.com', '0798888888', true),
('sdavis', 'NAT1009', 'sdavis@example.com', '0799999999', true),
('mrodriguez', 'NAT1010', 'mrodriguez@example.com', '0790000000', true),
('wmoore', 'NAT1011', 'wmoore@example.com', '0791234567', true); -- user with insufficient data (only 2 salary records)

-- Insert salary data for NAT1001 (jdoe) - Should result in RED status
-- Jan 2025: 1200, Feb 2025: 1300, Mar 2025: 1400, May 2025: 1500, Jun 2025: 1600 (summer)
INSERT INTO salaries (year, month, salary, user_id) VALUES
(2025, 1, 1200.00, 1),  -- January
(2025, 2, 1300.00, 1),  -- February  
(2025, 3, 1400.00, 1),  -- March
(2025, 5, 1500.00, 1),  -- May
(2025, 6, 1600.00, 1);  -- June (summer month, will be adjusted)

-- Insert salary data for NAT1002 (asmith) - Should result in ORANGE status (average = 2000)
INSERT INTO salaries (year, month, salary, user_id) VALUES
(2025, 1, 1900.00, 2),
(2025, 2, 2000.00, 2),
(2025, 3, 2100.00, 2);

-- Insert salary data for NAT1004 (kwilson) - Should result in GREEN status (average > 2000)
INSERT INTO salaries (year, month, salary, user_id) VALUES
(2025, 1, 2500.00, 4),
(2025, 2, 2600.00, 4),
(2025, 3, 2700.00, 4),
(2025, 4, 2800.00, 4);

-- Insert salary data for NAT1005 (lbrown) - High salary with tax deduction scenario
INSERT INTO salaries (year, month, salary, user_id) VALUES
(2025, 1, 3000.00, 5),
(2025, 2, 3100.00, 5),
(2025, 3, 3200.00, 5),
(2025, 4, 3300.00, 5); -- Total > 10000, should trigger tax deduction

-- Insert salary data for NAT1006 (tjones) - December bonus scenario
INSERT INTO salaries (year, month, salary, user_id) VALUES
(2024, 12, 2000.00, 6), -- December with +10% bonus
(2025, 1, 1800.00, 6),
(2025, 2, 1900.00, 6);

-- Insert salary data for NAT1007 (dgarcia) - Summer months scenario
INSERT INTO salaries (year, month, salary, user_id) VALUES
(2025, 6, 2000.00, 7),  -- June (summer, -5%)
(2025, 7, 2100.00, 7),  -- July (summer, -5%)
(2025, 8, 2200.00, 7);  -- August (summer, -5%)

-- Insert salary data for NAT1008 (rmiller) - Mixed seasonal adjustments
INSERT INTO salaries (year, month, salary, user_id) VALUES
(2024, 12, 1800.00, 8), -- December (+10%)
(2025, 6, 1900.00, 8),  -- June (-5%)
(2025, 7, 2000.00, 8),  -- July (-5%)
(2025, 8, 2100.00, 8);  -- August (-5%)

-- Insert salary data for NAT1009 (sdavis) - Regular months, no adjustments
INSERT INTO salaries (year, month, salary, user_id) VALUES
(2025, 1, 1500.00, 9),
(2025, 2, 1600.00, 9),
(2025, 3, 1700.00, 9);

-- Insert salary data for NAT1010 (mrodriguez) - Border case around 2000 average
INSERT INTO salaries (year, month, salary, user_id) VALUES
(2025, 1, 1950.00, 10),
(2025, 2, 2000.00, 10),
(2025, 3, 2050.00, 10);

-- Insert insufficient salary data for NAT1011 (wmoore) - Only 2 records (should trigger INSUFFICIENT_DATA error)
INSERT INTO salaries (year, month, salary, user_id) VALUES
(2025, 1, 2000.00, 11),
(2025, 2, 2100.00, 11);