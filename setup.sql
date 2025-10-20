DROP TABLE IF EXISTS accounts CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS account_roles CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS regions CASCADE;
DROP TABLE IF EXISTS clients CASCADE;
DROP TABLE IF EXISTS comments CASCADE;
DROP TABLE IF EXISTS contacts CASCADE;
DROP TABLE IF EXISTS research_type CASCADE;
DROP TABLE IF EXISTS manufacturers CASCADE;
DROP TABLE IF EXISTS devices CASCADE;
DROP TABLE IF EXISTS classificators CASCADE;
DROP TABLE IF EXISTS ticket_statuses CASCADE;
DROP TABLE IF EXISTS ticket_types CASCADE;
DROP TABLE IF EXISTS ticket_reasons CASCADE;
DROP TABLE IF EXISTS tickets CASCADE;
DROP TABLE IF EXISTS attachments CASCADE;
DROP TABLE IF EXISTS departments CASCADE;
DROP TABLE IF EXISTS agreements CASCADE;


-- Enable citext extension
CREATE EXTENSION IF NOT EXISTS citext;

-- Appends comment ID to corresponding tables
-- CREATE OR REPLACE FUNCTION append_comment_to_reference()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     -- Handle INSERT operations
--     IF TG_OP = 'INSERT' THEN
--         -- Check if reference_id exists in clients table
--         IF EXISTS (SELECT 1 FROM clients WHERE id = NEW.reference_id) THEN
--             UPDATE clients
--             SET comments = array_append(COALESCE(comments, '{}'), NEW.id)
--             WHERE id = NEW.reference_id;
--         END IF;

--     -- Handle DELETE operations
--     ELSIF TG_OP = 'DELETE' THEN
--         -- Check if reference_id exists in clients table
--         IF EXISTS (SELECT 1 FROM clients WHERE id = OLD.reference_id) THEN
--             UPDATE clients
--             SET comments = array_remove(COALESCE(comments, '{}'), OLD.id)
--             WHERE id = OLD.reference_id;
--         END IF;
--     END IF;

--     RETURN COALESCE(NEW, OLD);
-- END;
-- $$ LANGUAGE plpgsql;

BEGIN;

CREATE TABLE IF NOT EXISTS accounts (
    user_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    username CITEXT NOT NULL UNIQUE,
    database TEXT NOT NULL,
    disabled BOOLEAN DEFAULT false,
    password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS departments (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    title TEXT NOT NULL
);

INSERT INTO departments (id, title) VALUES
    ('0d7cf0f9-bbea-11ed-8100-40b0765b1e01', 'Отдел Mindray'),
    ('1f62a255-ef3a-11e5-8d88-001a64d22812', 'Сервис ВББ'),
    ('1f62a256-ef3a-11e5-8d88-001a64d22812', 'Отдел биохимии'),
    ('1f62a257-ef3a-11e5-8d88-001a64d22812', 'Отдел ПЦР'),
    ('1f62a259-ef3a-11e5-8d88-001a64d22812', 'Сервис А78'),
    ('91d50980-4c47-11ec-80ec-40b0765b1e01', 'Методисты СОЛТ'),
    ('9d798b43-9bdd-11ed-80fe-40b0765b1e01', 'Отдел исследования мочи'),
    ('add49497-f8f5-11e6-a1f7-001a64d22812', 'Сервис ЛИС'),
    ('c733447e-2567-11ea-80cd-40b0765b1e01', 'Отдел общих клинических исследований');

CREATE TABLE IF NOT EXISTS users (
    user_id UUID NOT NULL REFERENCES accounts(user_id) ON DELETE CASCADE,
    first_name TEXT,
    last_name TEXT,
    department UUID REFERENCES departments(id) ON DELETE SET NULL,
    email TEXT,
    phone TEXT UNIQUE,
    user_pic UUID
);

-- test123
INSERT INTO accounts (user_id, username, database, disabled, password_hash) VALUES
    ('ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e', 'admin', 'foxygendb', false, '$2a$10$TLTo5KFUlITFAWC.cDk9m.LtlUy22omjg3btZ7AuPi1lqmJRVwKLm'),
    ('84d512de-df6a-4a0b-be28-a8e184bd1d6a', 'coordinator', 'foxygendb', false, '$2a$10$TLTo5KFUlITFAWC.cDk9m.LtlUy22omjg3btZ7AuPi1lqmJRVwKLm'),
    ('73c97b16-09b1-416e-94ad-f8952be14a19', 'user1', 'foxygendb', false, '$2a$10$TLTo5KFUlITFAWC.cDk9m.LtlUy22omjg3btZ7AuPi1lqmJRVwKLm'),
    ('ccb5418b-ac05-4f2c-8bab-6e76a51f86d9', 'user2', 'foxygendb', false, '$2a$10$TLTo5KFUlITFAWC.cDk9m.LtlUy22omjg3btZ7AuPi1lqmJRVwKLm');

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
    ('ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e', 1),
    ('84d512de-df6a-4a0b-be28-a8e184bd1d6a', 2),
    ('73c97b16-09b1-416e-94ad-f8952be14a19', 3),
    ('ccb5418b-ac05-4f2c-8bab-6e76a51f86d9', 3);


INSERT INTO users (user_id, first_name, last_name, department, email, phone, user_pic) VALUES
    ('ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e', 'Админ', '', '1f62a255-ef3a-11e5-8d88-001a64d22812', 'test1@gmail.com', 79992141831, 'ad1fa321-cad1-7bc5-b3e5-f4a3f23cf90e'),
    ('84d512de-df6a-4a0b-be28-a8e184bd1d6a', 'Координатор', '', '1f62a259-ef3a-11e5-8d88-001a64d22812', 'test2@gmail.com', 79992141832, 'ad1fa321-cad1-7bc5-b3e5-f4a3f23cf90e'),
    ('73c97b16-09b1-416e-94ad-f8952be14a19', 'Пользователь', '1', '1f62a256-ef3a-11e5-8d88-001a64d22812', 'test3@gmail.com', 79992146832, 'ad1fa321-cad1-7bc5-b3e5-f4a3f23cf90e'),
    ('ccb5418b-ac05-4f2c-8bab-6e76a51f86d9', 'Пользователь', '2', 'add49497-f8f5-11e6-a1f7-001a64d22812', 'test4@gmail.com', 79992142732, 'ad1fa321-cad1-7bc5-b3e5-f4a3f23cf90e');

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
    author_id UUID NOT NULL REFERENCES accounts(user_id) ON DELETE CASCADE,
    reference_id UUID NOT NULL,
    text TEXT,
    created_at timestamp
);

-- Create the trigger
-- CREATE OR REPLACE TRIGGER trigger_append_comment
-- AFTER INSERT OR DELETE ON comments
-- FOR EACH ROW
-- EXECUTE FUNCTION append_comment_to_reference();

-- Clients
CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    region INT REFERENCES regions(id) ON DELETE SET NULL,
    address TEXT,
    location JSONB DEFAULT '{"lat": 0, "lng": 0}',
    laboratory_system UUID,
    manager UUID[] DEFAULT '{}'
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
    'Центральная Больница г. Беломорск',
    '77', -- Reference to Moscow region (must exist in regions table)
    'Беломорск, Меретсковая ул., 6',
    '{"lat": 55.7539, "lng": 37.6208}', -- JSONB location with coordinates
    'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', -- TODO: ADD Laboraroty reference
    '{ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e}'
);

-- Insert a comment for a client
INSERT INTO comments (author_id, reference_id, text, created_at)
VALUES ('73c97b16-09b1-416e-94ad-f8952be14a19', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Комментарий о проделанной работе', (NOW() AT TIME ZONE 'UTC'));
INSERT INTO comments (author_id, reference_id, text, created_at)
VALUES ('ccb5418b-ac05-4f2c-8bab-6e76a51f86d9', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Еще один комментарий', (NOW() AT TIME ZONE 'UTC'));


-- Contacts
CREATE TABLE IF NOT EXISTS contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    position TEXT,
    phone TEXT UNIQUE,
    email TEXT UNIQUE,
    client_id UUID REFERENCES clients(id) ON DELETE CASCADE
);

INSERT INTO contacts (id, name, position, phone, email, client_id)
VALUES ('27b1c3f2-f196-4885-8d56-9169e9f71e52', 'Вероника Васильевна', 'Заведующая лабораторией', '79992191217', 'someemail@gmail.com','a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11');
INSERT INTO contacts (name, phone, email, client_id)
VALUES ('Иван Иванович', '79992161721', 'grdandhedne@gmail.com','a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11');

-- Research Type
CREATE TABLE IF NOT EXISTS research_type (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT UNIQUE
);

INSERT INTO research_type (title) VALUES
('Электрофорез белков'),
('Автоматизация'),
('Скрытая кровь'),
('Пластик'),
('Аллергология'),
('Бактериология'),
('Биохимия'),
('Газы крови'),
('Гематология'),
('Группы крови'),
('Иммунохимия'),
('ИФА'),
('Коагуалогия'),
('Моча'),
('ПЦР'),
('СОЭ'),
('Водоподготовка'),
('Цитология');

-- Manufacturers
CREATE TABLE IF NOT EXISTS manufacturers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT UNIQUE
);

INSERT INTO manufacturers (title) VALUES
('СОЛТ'),
('Урит'),
('АО "Вектор-Бест-Балтика"'),
('West Medica Produktions');

-- Classificator
CREATE TABLE IF NOT EXISTS classificators (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT,
    manufacturer UUID REFERENCES manufacturers(id) ON DELETE SET NULL,
    research_type UUID REFERENCES research_type(id) ON DELETE SET NULL,
    registration_certificate JSONB DEFAULT '{}',
    maintenance_regulations JSONB DEFAULT '{}',
    attachments TEXT[] DEFAULT '{}',
    images TEXT[] DEFAULT '{}'
);

INSERT INTO classificators (id, title)
VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Экспресс-анализатор Triage MeterPro');

-- Devices
CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    classificator UUID REFERENCES classificators(id) ON DELETE SET NULL,
    serial_number TEXT UNIQUE,
    properties JSONB DEFAULT '{}',
    connected_to_lis BOOLEAN DEFAULT FALSE,
    is_used BOOLEAN DEFAULT FALSE
);

-- Insert a new device
INSERT INTO devices (id, classificator, serial_number)
VALUES (
    '2ecc4df8-cd7a-412d-9362-09b047a67c30',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'SN123456'
);
INSERT INTO devices (classificator, serial_number, properties)
VALUES (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'SN123453',
    '{"manufacturer": "Company XYZ", "model": "Device Pro", "firmware": "v2.1"}'
);

-- Ticket statuses
CREATE TABLE IF NOT EXISTS ticket_statuses (
    type VARCHAR(128) PRIMARY KEY,
    title TEXT
);

INSERT INTO ticket_statuses (type, title) VALUES
('created', 'создан'),
('assigned', 'назначен'),
('inWork', 'в работе'),
('worksDone', 'работы завершены'),
('closed', 'закрыт'),
('cancelled', 'отменен');

-- Ticket types
CREATE TABLE IF NOT EXISTS ticket_types (
    type VARCHAR(128) PRIMARY KEY,
    title TEXT
);

INSERT INTO ticket_types (type, title) VALUES
('internal', 'внутренний'),
('external', 'внешний');

-- Ticket reasons
CREATE TABLE IF NOT EXISTS ticket_reasons (
    id VARCHAR(128) PRIMARY KEY,
    title TEXT,
    past TEXT,
    present TEXT,
    future TEXT
);

INSERT INTO ticket_reasons (id, title, past, present, future) VALUES
('commissioning', 'Ввод в эксплуатацию', 'Введен в эксплуатацию', 'Ввод в эксплуатацию', 'Ввести в эксплуатацию'),
('consultation', 'Консультация', 'Проведена консультация', 'Консультация', 'Провести консультацию'),
('deinstallation', 'Деинсталляция', 'Деинсталлирован', 'Деинсталляция', 'Деинсталлировать'),
('diagnostic', 'Диагностика', 'Проведена диагностика', 'Проведение диагностики', 'Провести диагностику'),
('installation', 'Инсталляция', 'Инсталлирован', 'Инсталляция', 'Инсталлировать'),
('maintanence', 'Техническое обслуживание', 'Проведено ТО', 'Проведение ТО', 'Провести ТО'),
('methodInput', 'Ввод методик', 'Введены методики', 'Ввод методик', 'Ввести методики'),
('other', 'Прочее', 'Прочее', 'Прочее', 'Прочее'),
('repair', 'Ремонт', 'Проведен ремонт', 'Ремонт', 'Отремонтировать'),
('service', 'Сервисный центр', 'Сервисный центр', 'Сервисный центр', 'Сервисный центр'),
('staffTraining', 'Обучение', 'Провести обучение', 'Обучение', 'Проведено обучение');

CREATE TABLE IF NOT EXISTS tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number TEXT,
    created_at timestamp DEFAULT (NOW() AT TIME ZONE 'UTC'),
    assigned_at timestamp DEFAULT (NOW() AT TIME ZONE 'UTC'), -- Change to NULL later,
    workstarted_at timestamp DEFAULT NULL,
    workfinished_at timestamp DEFAULT NULL,
    deadline timestamp DEFAULT '2025-10-15T09:19:34.169Z',
    urgent BOOLEAN DEFAULT false,
    closed_at timestamp DEFAULT NULL,
    client UUID REFERENCES clients(id) ON DELETE SET NULL,
    device UUID REFERENCES devices(id) ON DELETE SET NULL,
    ticket_type VARCHAR(128) REFERENCES ticket_types(type) ON DELETE SET NULL,
    author UUID REFERENCES accounts(user_id) ON DELETE SET NULL,
    department UUID DEFAULT NULL,
    assigned_by UUID REFERENCES accounts(user_id) ON DELETE SET NULL,
    reason VARCHAR(128) REFERENCES ticket_reasons(id) ON DELETE SET NULL,
    description TEXT,
    contact_person UUID REFERENCES contacts(id) ON DELETE SET NULL,
    executor UUID REFERENCES accounts(user_id) ON DELETE SET NULL,
    status VARCHAR(128) REFERENCES ticket_statuses(type) ON DELETE SET NULL,
    result TEXT DEFAULT '',
    used_materials UUID[] DEFAULT '{}',
    recommendation TEXT
);

INSERT INTO tickets (number, client, device, ticket_type, author, assigned_by, reason, contact_person, executor, status, description, deadline, urgent, department) VALUES
('0002314', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2ecc4df8-cd7a-412d-9362-09b047a67c30', 'internal', 'ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e', '84d512de-df6a-4a0b-be28-a8e184bd1d6a', 'installation', '27b1c3f2-f196-4885-8d56-9169e9f71e52', '73c97b16-09b1-416e-94ad-f8952be14a19', 'assigned', 'Контроль прохождения 9004 ', '2025-10-17T09:19:34.169Z', true, '1f62a256-ef3a-11e5-8d88-001a64d22812');

INSERT INTO tickets (number, client, device, ticket_type, author, assigned_by, reason, contact_person, executor, status, description, deadline, urgent, department) VALUES
('04144', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2ecc4df8-cd7a-412d-9362-09b047a67c30', 'internal', 'ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e', '84d512de-df6a-4a0b-be28-a8e184bd1d6a', 'diagnostic', '27b1c3f2-f196-4885-8d56-9169e9f71e52', '73c97b16-09b1-416e-94ad-f8952be14a19', 'assigned', 'Выдаёт ошибку холостой пробы, превышение предела RBC. При выполнении анализов не считает эритроциты.', '2025-10-13T09:19:34.169Z', true, '1f62a256-ef3a-11e5-8d88-001a64d22812');

INSERT INTO tickets (number, created_at, client, device, ticket_type, author, assigned_by, reason, contact_person, executor, status, description, department) VALUES
('0002311', '2025-09-18T11:24:42.072Z', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2ecc4df8-cd7a-412d-9362-09b047a67c30', 'internal', 'ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e', '84d512de-df6a-4a0b-be28-a8e184bd1d6a', 'installation', '27b1c3f2-f196-4885-8d56-9169e9f71e52', '73c97b16-09b1-416e-94ad-f8952be14a19', 'assigned', 'Описание тикета', '1f62a256-ef3a-11e5-8d88-001a64d22812');

CREATE TABLE IF NOT EXISTS attachments (
    id SERIAL PRIMARY KEY,
    file_name TEXT NOT NULL,
    original_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    ref_id UUID NOT NULL
);

-- Agreements
CREATE TABLE IF NOT EXISTS agreements (
    id UUID PRIMARY KEY,
    number TEXT,
    actual_client UUID REFERENCES clients(id) ON DELETE SET NULL,
    distributor UUID REFERENCES clients(id) ON DELETE SET NULL,
    device UUID REFERENCES devices(id) ON DELETE SET NULL,
    assigned_at timestamp DEFAULT NULL,
    finished_at timestamp DEFAULT NULL,
    is_active BOOLEAN DEFAULT true,
    on_warranty BOOLEAN DEFAULT true,
    type VARCHAR(128)
);

INSERT INTO agreements (id, actual_client, distributor, device, assigned_at, type) VALUES
    (gen_random_uuid(), 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2ecc4df8-cd7a-412d-9362-09b047a67c30', NOW(), 'rent');
INSERT INTO agreements (id, actual_client, distributor, assigned_at, type) VALUES
    (gen_random_uuid(), 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NOW(), 'bought');

COMMIT;
