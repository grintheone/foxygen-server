BEGIN;

CREATE TABLE accounts (
    user_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    username CITEXT NOT NULL UNIQUE,
    database TEXT NOT NULL,
    disabled BOOLEAN DEFAULT false,
    password_hash TEXT NOT NULL
);

CREATE TABLE roles (
    id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255)
);

INSERT INTO roles (id, name, description) VALUES
    (1, 'admin', 'System administrator with full access'),
    (2, 'coordinator', 'Can manage content and users but not system settings'),
    (3, 'user', 'Regular user with basic access');

CREATE TABLE account_roles (
    user_id UUID REFERENCES accounts(user_id) ON DELETE CASCADE,
    role_id INT REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

COMMIT;
