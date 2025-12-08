CREATE EXTENSION IF NOT EXISTS btree_gist;

ALTER TABLE bookings ADD CONSTRAINT no_overlapping_bookings
    EXCLUDE USING GIST (
    hotel_id WITH =,
    room_id WITH =,
    daterange(check_in, check_out, '[)') WITH &&
);