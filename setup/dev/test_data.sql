
SET SEARCH_PATH TO community, auth;

-- AUTHENTICATION --
-- due to the ability to share categories or plan between collaborators, we have test cases for different users as well. --
-- new test users are automatically added at the end of the script so that they can test sharable groups and containers --
-- two users are admin@gmail.com or user@gmail.com --
INSERT INTO auth.users(email, username, role, encrypted_password, birthdate, is_verified, email_confirmed_at) values ('admin@gmail.com', 'batman', 'TESTER', '$2a$08$DtBpJRoDduzqR/ERz/JvFe2toYO9UhifP/.kmxdeamz0VmWr7kQuW', '1984-01-12', 'true', now());

-- UPDATE PROFILE TABLE WITH USER DETAILS --
UPDATE community.profiles SET 
username = 'john',
full_name = 'John Doe',
phone_number = '1234567890',
about_me = 'I like to climb trees and hike with my friends'
WHERE email_address = 'admin@gmail.com';

-- ADD NOTES SQL TEST DATA --
INSERT INTO community.notes (title, description, status, color, location, created_by, updated_by, sharable_groups)
VALUES (
    'Buy kitty litter for four of my kittens',
    'Do not buy the brand from walmart, buy from a generic well known place',
    (SELECT id FROM community.statuses s LIMIT 1),
    '#2a6dbc',
     '(-72.625481, 42.203217)',
    (SELECT id FROM community.profiles p LIMIT 1),
    (SELECT id FROM community.profiles p LIMIT 1),
    ARRAY[(SELECT id FROM community.profiles p LIMIT 1)::UUID]
);

-- ADD CATGORIES SQL TEST DATA ---
INSERT INTO community.category (name, description, status, color, created_by, updated_by, sharable_groups)
VALUES ('Groceries', 'used for grocery related items', (SELECT id FROM community.statuses s LIMIT 1), '#d20a0a',  (SELECT id FROM community.profiles p LIMIT 1), (SELECT id FROM community.profiles p LIMIT 1), ARRAY [(SELECT id FROM community.profiles p LIMIT 1)::UUID]),
       ('Utilities', 'store boxes and utility related stuffs', (SELECT id FROM community.statuses s LIMIT 1), '#e7d3da', (SELECT id FROM community.profiles p LIMIT 1), (SELECT id FROM community.profiles p LIMIT 1), ARRAY [(SELECT id FROM community.profiles p LIMIT 1)::UUID]),
       ('Entertainment', 'store gaming consoles, tvs, and / or any audio video equipment', (SELECT id FROM community.statuses s LIMIT 1), '#963256', (SELECT id FROM community.profiles p LIMIT 1), (SELECT id FROM community.profiles p LIMIT 1), ARRAY [(SELECT id FROM community.profiles p LIMIT 1)::UUID]);

-- ADD MAINTENANCE PLAN SQL TEST DATA --
INSERT INTO community.maintenance_plan 
(name, description, status, color, plan_type, plan_due, location, created_by, updated_by, sharable_groups)
VALUES 
('Daily maintenance plan', 'used to validate daily items', 
 (SELECT id FROM community.statuses mps LIMIT 1), 
 '#d20a0a', 'annual', 
 now() + interval '1 year', 
 '(-119.170898, 34.196411)', 
 (SELECT id FROM community.profiles p LIMIT 1), 
 (SELECT id FROM community.profiles p LIMIT 1), 
 ARRAY [(SELECT id FROM community.profiles p LIMIT 1)::UUID]),

('Weekly maintenance plan', 'used for weekly maintenance', 
 (SELECT id FROM community.statuses mps LIMIT 1), 
 '#28b391', 'weekly', 
 now() + interval '7 days', 
 '(-117.182541, 34.055569)', 
 (SELECT id FROM community.profiles p LIMIT 1), 
 (SELECT id FROM community.profiles p LIMIT 1), 
 ARRAY [(SELECT id FROM community.profiles p LIMIT 1)::UUID]);

-- ADD INVENTORIES SQL TEST DATA --
INSERT INTO community.inventory (name, description, price, status, barcode, sku, quantity, bought_at, location,
                                 storage_location_id, is_returnable, return_location, return_datetime, max_weight, min_weight,
                                 max_height, min_height, created_by, updated_by, sharable_groups)
VALUES ('4 pounds of kitty litter',
        'Bought from tractor supply in fm969',
        12.00,
        'DRAFT',
        'barcode#1123928',
        'sku#123456734',
        12,
        'Walmart',
        'Kitchen Pantry',
        (SELECT id from community.storage_locations WHERE location = 'Kitchen Pantry'),
        false,
        'amazon return',
        'now()',
        '12',
        '4',
        '20',
        '12',
        (SELECT id FROM community.profiles p LIMIT 1),
        (SELECT id FROM community.profiles p LIMIT 1),
        ARRAY [(SELECT id FROM community.profiles p LIMIT 1)::UUID]);

-- insert a null datetime --
INSERT INTO community.inventory (name, description, price, status, barcode, sku, quantity, bought_at, location,
                                 storage_location_id, is_returnable, return_location, max_weight, min_weight,
                                 max_height, min_height, created_by, updated_by, sharable_groups)
VALUES ('Dog food',
        '6 pounds of food bought from tractor supply',
        96.00,
        'HIDDEN',
        'barcode#1123928',
        'sku#123456734',
        1,
        'Walmart',
        'Utility Closet',
        (SELECT id from community.storage_locations WHERE location = 'Utility Closet'),
        false,
        'amazon return',
        '12',
        '4',
        '20',
        '12',
        (SELECT id FROM community.profiles p LIMIT 1),
        (SELECT id FROM community.profiles p LIMIT 1),
        ARRAY [(SELECT id FROM community.profiles p LIMIT 1)::UUID]);


-- ADD FAVOURITE ITEMS --
INSERT INTO community.favourite_items (category_id, created_by, updated_by, sharable_groups)
VALUES (
    (SELECT id FROM community.category c LIMIT 1),
    (SELECT id FROM community.profiles p LIMIT 1),
    (SELECT id FROM community.profiles p LIMIT 1),
    ARRAY[(SELECT id FROM community.profiles p LIMIT 1)]
);

INSERT INTO community.favourite_items (maintenance_plan_id, created_by, updated_by, sharable_groups)
VALUES (
    (SELECT id FROM community.maintenance_plan mp LIMIT 1),
    (SELECT id FROM community.profiles p LIMIT 1),
    (SELECT id FROM community.profiles p LIMIT 1),
    ARRAY[(SELECT id FROM community.profiles p LIMIT 1)]
);


-- AUTHENTICATION FOR REGULAR USER --
INSERT INTO auth.users(email, username, role, encrypted_password, birthdate, is_verified, email_confirmed_at) values ('user@gmail.com', 'dark knight', 'TESTER', '$2a$08$DtBpJRoDduzqR/ERz/JvFe2toYO9UhifP/.kmxdeamz0VmWr7kQuW', '1900-01-01', 'true', now());