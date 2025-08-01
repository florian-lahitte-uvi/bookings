-- SQL in section 'Down' is executed when this migration is rolled back
-- Remove the auto-increment from room_restrictions id column
ALTER TABLE room_restrictions ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE IF EXISTS room_restrictions_id_seq;