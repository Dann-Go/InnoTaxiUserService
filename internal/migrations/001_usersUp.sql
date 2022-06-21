CREATE TABLE IF NOT EXISTS users (
                       ID SERIAL PRIMARY KEY ,
                       Name TEXT,
                       Phone TEXT,
                       Email TEXT,
                       Password_hash TEXT,
                       Rating FLOAT DEFAULT 0.0
);

CREATE INDEX IF NOT EXISTS phone_idx ON users (phone);
CREATE INDEX IF NOT EXISTS email_idx ON users (email);


