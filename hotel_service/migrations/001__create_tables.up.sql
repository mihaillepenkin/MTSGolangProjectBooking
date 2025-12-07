CREATE TABLE IF NOT EXISTS hotels {
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(30) NOT NULL,
    description VARCHAR(500) NOT NULL,
    location VARCHAR(50) NOT NULL,
    owner_id BIGINT NOT NULL
}

CREATE TABLE IF NOT EXISTS rooms {
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number INT NOT NULL,
    price INT NOT NULL,
    hotel_id BIGINT NOT NULL,
    CONSTRAINT fk_hotel FOREIGN KEY (hotel_id) REFERENCES hotels(id)
}
