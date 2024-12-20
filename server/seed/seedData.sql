

-- seed community status --
INSERT INTO community.statuses(name, description) 
VALUES ('draft', 'items under this bucket are in draft state'),
('archived', 'items under this bucket are archived'),
('completed', 'items under this bucket are marked complete'),
('pending', 'items under this bucket are in pending state'),
('urgent', 'items under this bucket require immediate attention'),
('general', 'items under this bucket are generalized items'),
('on_hold', 'items under this bucket are on hold and needs more information'),
('cancelled', 'items under this bucket are cancelled and pending for deletion');

-- seed storage locations --
INSERT INTO community.storage_locations (location)
VALUES ('Kitchen Pantry'),
       ('Master Bedroom Closet'),
       ('Garage'),
       ('Living Room Cabinet'),
       ('Bathroom Closet'),
       ('Dining Room Hutch'),
       ('Home Office Desk'),
       ('Basement Storage'),
       ('Kids'' Playroom'),
       ('Garage Workshop'),
       ('Guest Bedroom Closet'),
       ('Outdoor Shed'),
       ('Utility Closet'),
       ('Attic Storage'),
       ('Guest Bathroom Cabinet'),
       ('Children''s Bedroom Closet'),
       ('Outdoor Storage Box'),
       ('Home Gym Closet'),
       ('Patio Storage Bench'),
       ('Study Room Bookshelf'),
       ('Laundry Room Cabinet'),
       ('Home Theater Shelf'),
       ('Backyard Storage Shed'),
       ('Closet Under the Stairs'),
       ('Home Bar Cabinet'),
       ('Tool Shed'),
       ('Linen Closet'),
       ('Shoe Rack'),
       ('Home Library'),
       ('Tool Bench');