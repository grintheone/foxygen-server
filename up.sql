BEGIN;

CREATE TABLE IF NOT EXISTS accounts (
    user_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    username CITEXT NOT NULL UNIQUE,
    database TEXT NOT NULL,
    disabled BOOLEAN DEFAULT false,
    password_hash TEXT NOT NULL
);

-- test123
INSERT INTO accounts (user_id, username, database, disabled, password_hash) VALUES
    ('ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e', 'admin', 'foxygendb', false, '$2a$10$TLTo5KFUlITFAWC.cDk9m.LtlUy22omjg3btZ7AuPi1lqmJRVwKLm');

CREATE TABLE IF NOT EXISTS roles (
    id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255)
);

INSERT INTO roles (id, name, description) VALUES
    (1, 'admin', 'System administrator with full access'),
    (2, 'coordinator', 'Can manage content and users but not system settings'),
    (3, 'user', 'Regular user with basic access');

CREATE TABLE IF NOT EXISTS account_roles (
    user_id UUID NOT NULL REFERENCES accounts(user_id) ON DELETE CASCADE,
    role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    UNIQUE (user_id, role_id)
);

INSERT INTO account_roles (user_id, role_id) VALUES
    ('ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e', 1);

COMMIT;
