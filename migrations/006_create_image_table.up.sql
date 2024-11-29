CREATE TABLE properties_images (
                                  id SERIAL PRIMARY KEY,
                                  property_id BIGINT NOT NULL,
                                  image_url TEXT NOT NULL,
                                  FOREIGN KEY (property_id) REFERENCES properties(id)
                            -- URL или путь к файлу изображения

);