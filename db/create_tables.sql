
BEGIN;

CREATE EXTENSION IF NOT EXISTS citext;

-- DONE
CREATE TABLE accounts (
    user_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    username CITEXT UNIQUE,
    disabled BOOLEAN DEFAULT false,
    password_hash TEXT NOT NULL
);

-- DONE
CREATE TABLE departments (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL
);

-- DONE
CREATE TABLE users (
    user_id UUID PRIMARY KEY REFERENCES accounts(user_id) ON DELETE CASCADE,
    first_name TEXT DEFAULT '',
    last_name TEXT DEFAULT '',
    department UUID REFERENCES departments(id) ON DELETE SET NULL,
    email TEXT DEFAULT '',
    phone TEXT DEFAULT '',
    logo TEXT DEFAULT '',
    latest_ticket UUID DEFAULT NULL
);

-- DONE
CREATE TABLE roles (
    id INT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255)
);

-- DONE
CREATE TABLE account_roles (
    user_id UUID PRIMARY KEY REFERENCES accounts(user_id) ON DELETE CASCADE,
    role_id INT NOT NULL REFERENCES roles(id)
);

-- DONE
CREATE TABLE regions (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL
);

-- DONE
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    author_id UUID NOT NULL REFERENCES accounts(user_id) ON DELETE CASCADE,
    reference_id UUID NOT NULL,
    text TEXT,
    created_at timestamp DEFAULT (NOW() AT TIME ZONE 'UTC')
);

-- DONE
CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    region UUID REFERENCES regions(id) ON DELETE SET NULL,
    address TEXT,
    location JSONB DEFAULT NULL,
    laboratory_system UUID DEFAULT NULL, -- ссылка на конкретный лис (пока не реализововывать)
    manager UUID[] DEFAULT '{}'
);

-- DONE
CREATE TABLE contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT DEFAULT '',
    position TEXT DEFAULT '',
    phone TEXT DEFAULT '',
    email TEXT DEFAULT '',
    client_id UUID REFERENCES clients(id) ON DELETE CASCADE
);

-- DONE
CREATE TABLE research_type (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT UNIQUE
);

-- DONE
CREATE TABLE manufacturers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT
);

-- DONE
CREATE TABLE classificators (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT,
    manufacturer UUID REFERENCES manufacturers(id) ON DELETE SET NULL,
    research_type UUID REFERENCES research_type(id) ON DELETE SET NULL,
    registration_certificate JSONB DEFAULT '{}',
    maintenance_regulations JSONB DEFAULT '{}',
    attachments TEXT[] DEFAULT '{}',
    images TEXT[] DEFAULT '{}'
);

-- DONE
CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    classificator UUID REFERENCES classificators(id) ON DELETE SET NULL,
    serial_number TEXT,
    properties JSONB DEFAULT '{}',
    connected_to_lis BOOLEAN DEFAULT FALSE,
    is_used BOOLEAN DEFAULT FALSE
);

-- DONE
CREATE TABLE ticket_statuses (
    type VARCHAR(128) PRIMARY KEY,
    title TEXT
);

-- DONE
CREATE TABLE ticket_types (
    type VARCHAR(128) PRIMARY KEY,
    title TEXT
);

-- DONE
CREATE TABLE ticket_reasons (
    id VARCHAR(128) PRIMARY KEY,
    title TEXT,
    past TEXT,
    present TEXT,
    future TEXT
);

CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number INT GENERATED ALWAYS AS IDENTITY,
    created_at timestamp DEFAULT (NOW() AT TIME ZONE 'UTC'),
    assigned_at timestamp DEFAULT NULL,
    workstarted_at timestamp DEFAULT NULL,
    workfinished_at timestamp DEFAULT NULL,
    planned_start timestamp DEFAULT NULL, -- manager field
    planned_end timestamp DEFAULT NULL,  -- manager field
    assigned_start timestamp DEFAULT NULL, -- coordinator field
    assigned_end timestamp DEFAULT NULL, -- coordinator field
    urgent BOOLEAN DEFAULT false,
    closed_at timestamp DEFAULT NULL,
    client UUID REFERENCES clients(id) ON DELETE SET NULL,
    device UUID REFERENCES devices(id) ON DELETE SET NULL,
    ticket_type VARCHAR(128) REFERENCES ticket_types(type) ON DELETE SET NULL,
    author UUID REFERENCES accounts(user_id) ON DELETE SET NULL,
    department UUID REFERENCES departments(id) ON DELETE SET NULL,
    assigned_by UUID REFERENCES accounts(user_id) ON DELETE SET NULL DEFAULT NULL,
    reason VARCHAR(128) REFERENCES ticket_reasons(id) ON DELETE SET NULL,
    description TEXT,
    contact_person UUID REFERENCES contacts(id) ON DELETE SET NULL,
    executor UUID REFERENCES accounts(user_id) ON DELETE SET NULL,
    status VARCHAR(128) REFERENCES ticket_statuses(type) ON DELETE SET NULL,
    result TEXT DEFAULT '',
    used_materials UUID[] DEFAULT '{}',
    reference_ticket UUID REFERENCES tickets(id) DEFAULT NULL,
    double_signed BOOLEAN DEFAULT FALSE
);

-- {
--   "id": "9d737222-cf42-4139-a393-b85ed61152ff.jpg",
--   "name": "IMG_20241129_140002.jpg",
--   "mediaType": "image/jpeg",
--   "imageType": "image",
--   "ext": "jpg"
-- },

CREATE TABLE attachments (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    media_type TEXT NOT NULL,
    ext TEXT NOT NULL,
    ref_id UUID NOT NULL
);

CREATE TABLE  agreements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number INT GENERATED ALWAYS AS IDENTITY,
    actual_client UUID REFERENCES clients(id) ON DELETE SET NULL,
    distributor UUID REFERENCES clients(id) ON DELETE SET NULL,
    device UUID REFERENCES devices(id) ON DELETE SET NULL,
    assigned_at timestamp,
    finished_at timestamp,
    is_active BOOLEAN DEFAULT true,
    on_warranty BOOLEAN DEFAULT true,
    type VARCHAR(128)
);

CREATE TABLE ra_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT DEFAULT ''
);

CREATE TABLE remote_access (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID REFERENCES devices(id) ON DELETE CASCADE,
    parameter_id UUID REFERENCES ra_options(id) ON DELETE SET NULL,
    value JSONB DEFAULT '{}'
);

COMMIT;
