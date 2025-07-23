CREATE TABLE IF NOT EXISTS city_phones (
    id              SERIAL PRIMARY KEY,
    city            TEXT NOT NULL,
    diversion_phone TEXT NOT NULL
);