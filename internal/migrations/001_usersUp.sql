CREATE TABLE IF NOT EXISTS users (
                       ID SERIAL PRIMARY KEY ,
                       Name TEXT,
                       Phone TEXT,
                       Email TEXT,
                       Password_hash TEXT
);

CREATE INDEX IF NOT EXISTS phone_idx ON users (phone);


