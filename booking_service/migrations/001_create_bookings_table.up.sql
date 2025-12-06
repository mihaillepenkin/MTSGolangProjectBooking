CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    hotel_id bigint NOT NULL,
    room_number bigint NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'RUB',
    check_in DATE NOT NULL,
    check_out DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'unpaid',
    payment_id text NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, hotel_id, room_number, check_in, check_out)
);

CREATE INDEX IF NOT EXISTS idx_bookings_hotel_room_dates ON bookings(hotel_id, room_number, check_in, check_out);
CREATE INDEX IF NOT EXISTS idx_bookings_hotel_dates ON bookings(hotel_id, check_in, check_out);
CREATE INDEX IF NOT EXISTS idx_bookings_user ON bookings(user_id);