CREATE TABLE bookings (
                          id BIGSERIAL PRIMARY KEY,
                          property_id BIGINT NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
                          user_id BIGINT NOT NULL,
                          check_in_date DATE NOT NULL,
                          check_out_date DATE NOT NULL,
                          total_price INT NOT NULL,
                          status TEXT NOT NULL CHECK (status IN ('confirmed', 'pending', 'cancelled')),
                          created_at DATE NOT NULL
);
