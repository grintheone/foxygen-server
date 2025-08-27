-- Appends comment ID to corresponding tables
CREATE OR REPLACE FUNCTION append_comment_to_reference()
RETURNS TRIGGER AS $$
BEGIN
    -- Handle INSERT operations
    IF TG_OP = 'INSERT' THEN
        -- Check if reference_id exists in clients table
        IF EXISTS (SELECT 1 FROM clients WHERE id = NEW.reference_id) THEN
            UPDATE clients
            SET comments = array_append(COALESCE(comments, '{}'), NEW.id)
            WHERE id = NEW.reference_id;
        END IF;

    -- Handle DELETE operations
    ELSIF TG_OP = 'DELETE' THEN
        -- Check if reference_id exists in clients table
        IF EXISTS (SELECT 1 FROM clients WHERE id = OLD.reference_id) THEN
            UPDATE clients
            SET comments = array_remove(COALESCE(comments, '{}'), OLD.id)
            WHERE id = OLD.reference_id;
        END IF;
    END IF;

    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

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

-- Regions
CREATE TABLE IF NOT EXISTS regions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    district VARCHAR(50) NOT NULL
);

-- Insert data into the regions table
INSERT INTO regions (name, district) VALUES
('Республика Адыгея (Адыгея)', 'Южный'),
('Республика Башкортостан', 'Приволжский'),
('Республика Бурятия', 'Сибирский'),
('Республика Алтай', 'Сибирский'),
('Республика Дагестан', 'Северо-Кавказский'),
('Республика Ингушетия', 'Северо-Кавказский'),
('Кабардино-Балкарская Республика', 'Северо-Кавказский'),
('Республика Калмыкия', 'Южный'),
('Карачаево-Черкесская Республика', 'Северо-Кавказский'),
('Республика Карелия', 'Северо-Западный'),
('Республика Коми', 'Северо-Западный'),
('Республика Марий Эл', 'Приволжский'),
('Республика Мордовия', 'Приволжский'),
('Республика Саха (Якутия)', 'Дальневосточный'),
('Республика Северная Осетия - Алания', 'Северо-Кавказский'),
('Республика Татарстан', 'Приволжский'),
('Республика Тыва', 'Сибирский'),
('Удмуртская Республика', 'Приволжский'),
('Республика Хакасия', 'Сибирский'),
('Чеченская Республика', 'Северо-Кавказский'),
('Чувашская Республика - Чувашия', 'Приволжский'),
('Алтайский край', 'Сибирский'),
('Краснодарский край', 'Южный'),
('Красноярский край', 'Сибирский'),
('Приморский край', 'Дальневосточный'),
('Ставропольский край', 'Северо-Кавказский'),
('Хабаровский край', 'Дальневосточный'),
('Амурская область', 'Дальневосточный'),
('Архангельская область', 'Северо-Западный'),
('Астраханская область', 'Южный'),
('Белгородская область', 'Центральный'),
('Брянская область', 'Центральный'),
('Владимирская область', 'Центральный'),
('Волгоградская область', 'Южный'),
('Вологодская область', 'Северо-Западный'),
('Воронежская область', 'Центральный'),
('Ивановская область', 'Центральный'),
('Иркутская область', 'Сибирский'),
('Калининградская область', 'Северо-Западный'),
('Калужская область', 'Центральный'),
('Камчатский край', 'Дальневосточный'),
('Кемеровская область', 'Сибирский'),
('Кировская область', 'Приволжский'),
('Костромская область', 'Центральный'),
('Курганская область', 'Уральский'),
('Курская область', 'Центральный'),
('Ленинградская область', 'Северо-Западный'),
('Липецкая область', 'Центральный'),
('Магаданская область', 'Дальневосточный'),
('Московская область', 'Центральный'),
('Мурманская область', 'Северо-Западный'),
('Нижегородская область', 'Приволжский'),
('Новгородская область', 'Северо-Западный'),
('Новосибирская область', 'Сибирский'),
('Омская область', 'Сибирский'),
('Оренбургская область', 'Приволжский'),
('Орловская область', 'Центральный'),
('Пензенская область', 'Приволжский'),
('Пермский край', 'Приволжский'),
('Псковская область', 'Северо-Западный'),
('Ростовская область', 'Южный'),
('Рязанская область', 'Центральный'),
('Самарская область', 'Приволжский'),
('Саратовская область', 'Приволжский'),
('Сахалинская область', 'Дальневосточный'),
('Свердловская область', 'Уральский'),
('Смоленская область', 'Центральный'),
('Тамбовская область', 'Центральный'),
('Тверская область', 'Центральный'),
('Томская область', 'Сибирский'),
('Тульская область', 'Центральный'),
('Тюменская область', 'Уральский'),
('Ульяновская область', 'Приволжский'),
('Челябинская область', 'Уральский'),
('Забайкальский край', 'Сибирский'),
('Ярославская область', 'Центральный'),
('Москва', 'Центральный'),
('Санкт-Петербург', 'Северо-Западный'),
('Еврейская автономная область', 'Дальневосточный'),
('Ненецкий автономный округ', 'Северо-Западный'),
('Ханты-Мансийский автономный округ - Югра', 'Уральский'),
('Чукотский автономный округ', 'Дальневосточный'),
('Ямало-Ненецкий автономный округ', 'Уральский'),
('Республика Крым', 'Южный'),
('Севастополь', 'Южный'),
('Иные территории', 'Южный');

-- Comments
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    author_id UUID NOT NULL,
    reference_id UUID NOT NULL,
    text TEXT,
    created_at timestamp
);

-- Create the trigger
CREATE OR REPLACE TRIGGER trigger_append_comment
AFTER INSERT OR DELETE ON comments
FOR EACH ROW
EXECUTE FUNCTION append_comment_to_reference();

-- Clients
CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    region INT REFERENCES regions(id) ON DELETE SET NULL,
    address TEXT,
    location JSONB DEFAULT '{"lat": 0, "lng": 0}',
    comments INT[] DEFAULT '{}',
    laboratory_system UUID,
    manager UUID REFERENCES accounts(user_id) ON DELETE SET NULL
);

INSERT INTO clients (
    id,
    title,
    region,
    address,
    location,
    laboratory_system,
    manager
)
VALUES (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', -- example UUID for client
    'Client Name LLC',
    '77', -- Reference to Moscow region (must exist in regions table)
    'Moscow, Red Square, 1',
    '{"lat": 55.7539, "lng": 37.6208}', -- JSONB location with coordinates
    'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', -- TODO: ADD Laboraroty reference
     'ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e'
);

-- Insert a comment for a client
INSERT INTO comments (author_id, reference_id, text, created_at)
VALUES (gen_random_uuid(), 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Client comment', NOW());
INSERT INTO comments (author_id, reference_id, text, created_at)
VALUES (gen_random_uuid(), 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'New comment', NOW());

COMMIT;
