CREATE TABLE properties (
                            id BIGSERIAL PRIMARY KEY,
                            owner_id BIGINT NOT NULL,
                            title TEXT NOT NULL,
                            location TEXT NOT NULL,
                            price INT NOT NULL,
                            property_type TEXT NOT NULL CHECK (property_type IN ('house', 'apartment')),
                            rental_type TEXT NOT NULL CHECK (rental_type IN ('shortTerm', 'longTerm')),
                            max_guests INT NOT NULL,
                            created_at DATE NOT NULL
);
