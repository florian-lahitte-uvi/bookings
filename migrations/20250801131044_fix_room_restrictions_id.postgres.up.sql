-- SQL in section 'Up' is executed when this migration is applied
-- Fix the room_restrictions id column to be auto-incrementing
CREATE SEQUENCE IF NOT EXISTS room_restrictions_id_seq;
ALTER TABLE room_restrictions ALTER COLUMN id SET DEFAULT nextval('room_restrictions_id_seq');
ALTER SEQUENCE room_restrictions_id_seq OWNED BY room_restrictions.id;
SELECT setval('room_restrictions_id_seq', COALESCE(MAX(id), 0) + 1, false) FROM room_restrictions;