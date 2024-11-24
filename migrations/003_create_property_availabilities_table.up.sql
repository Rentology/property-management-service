CREATE TABLE property_availability (
                                       id BIGSERIAL PRIMARY KEY,
                                       property_id BIGINT NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
                                       date DATE NOT NULL,
                                       is_available BOOLEAN NOT NULL DEFAULT TRUE,
                                       UNIQUE (property_id, date)
);
