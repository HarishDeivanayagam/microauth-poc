CREATE TABLE organizations (
    id          VARCHAR(36) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    domain      VARCHAR(255) NOT NULL UNIQUE,
    created_at  INTEGER NOT NULL,
    updated_at  INTEGER NOT NULL
);

CREATE TABLE members (
    id               VARCHAR(36) PRIMARY KEY,
    organization_id  VARCHAR(36) NOT NULL,
    user_id          VARCHAR(36) NOT NULL,
    role             VARCHAR(255) NOT NULL,
    app_role         VARCHAR(255) NOT NULL,
    UNIQUE (organization_id, user_id),
    FOREIGN KEY (organization_id) REFERENCES organizations (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);
