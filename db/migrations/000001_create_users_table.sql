CREATE TABLE users (
    usr_id SERIAL PRIMARY KEY,
    usr_username VARCHAR(255) NOT NULL UNIQUE,
    usr_email VARCHAR(255) NOT NULL UNIQUE,
    usr_password TEXT NOT NULL,
    usr_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    usr_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    usr_deleted_at TIMESTAMP DEFAULT NULL
);

CREATE OR REPLACE FUNCTION update_timestamp_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.usr_updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_timestamp_column();
