CREATE TABLE reviews (
                         id BIGSERIAL PRIMARY KEY,
                         property_id BIGINT NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
                         user_id BIGINT NOT NULL,
                         rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
                         comment TEXT,
                         created_at DATE NOT NULL
);
