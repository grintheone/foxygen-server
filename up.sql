BEGIN;

CREATE TABLE IF NOT EXISTS accounts (
    user_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    username CITEXT NOT NULL UNIQUE,
    database TEXT NOT NULL,
    disabled BOOLEAN DEFAULT false,
    password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    user_id UUID NOT NULL REFERENCES accounts(user_id) ON DELETE CASCADE,
    first_name TEXT,
    last_name TEXT,
    department UUID,
    email TEXT,
    phone TEXT UNIQUE,
    user_pic UUID
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


INSERT INTO users (user_id, first_name, last_name, department, email, phone, user_pic) VALUES
    ('ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e', 'test_name', 'test_lastname', 'ad2fa382-cad3-2bc1-b4e7-f4a4f12cf54e', 'test@gmail.com', 79992142832, 'ad1fa321-cad1-7bc5-b3e5-f4a3f23cf90e');

COMMIT;
