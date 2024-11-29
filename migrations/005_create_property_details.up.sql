CREATE TABLE property_details (
                                  property_id BIGINT PRIMARY KEY,
                                  floor INT,
                                  max_floor INT,
                                  area INT,
                                  rooms INT,
                                  house_creation_year INT,
                                  house_type TEXT,
                                  description TEXT,
                                  FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE
);
