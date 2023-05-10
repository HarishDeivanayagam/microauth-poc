CREATE TABLE users (
    id              VARCHAR(36) PRIMARY KEY,
    first_name      VARCHAR(255) NOT NULL,
    last_name       VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    is_email_verified BOOLEAN NOT NULL,
    password        VARCHAR(255) NOT NULL,
    created_at      INTEGER NOT NULL,
    updated_at      INTEGER NOT NULL,
    reset_otp       VARCHAR(255),
    reset_expiry    INTEGER
);
