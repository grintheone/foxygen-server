BEGIN;

-- INSERT INTO departments (id, title) VALUES
--     ('0d7cf0f9-bbea-11ed-8100-40b0765b1e01', 'Отдел ПЦР'),
--     ('1f62a255-ef3a-11e5-8d88-001a64d22812', 'Отдел ИФА'),
--     ('1f62a256-ef3a-11e5-8d88-001a64d22812', 'Отдел Биохимии'),
--     ('1f62a257-ef3a-11e5-8d88-001a64d22812', 'Инженерная служба');

-- test123
-- INSERT INTO accounts (user_id, username, disabled, password_hash) VALUES
--     ('ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e', 'admin', false, '$2a$10$TLTo5KFUlITFAWC.cDk9m.LtlUy22omjg3btZ7AuPi1lqmJRVwKLm'),
--     ('84d512de-df6a-4a0b-be28-a8e184bd1d6a', 'coordinator', false, '$2a$10$TLTo5KFUlITFAWC.cDk9m.LtlUy22omjg3btZ7AuPi1lqmJRVwKLm'),
--     ('73c97b16-09b1-416e-94ad-f8952be14a19', 'user1', false, '$2a$10$TLTo5KFUlITFAWC.cDk9m.LtlUy22omjg3btZ7AuPi1lqmJRVwKLm'),
--     ('ccb5418b-ac05-4f2c-8bab-6e76a51f86d9', 'user2', false, '$2a$10$TLTo5KFUlITFAWC.cDk9m.LtlUy22omjg3btZ7AuPi1lqmJRVwKLm');

INSERT INTO roles (id, name, description) VALUES
    (1, 'admin', 'System administrator with full access'),
    (2, 'coordinator', 'Can manage content and users but not system settings'),
    (3, 'user', 'Regular user with basic access');

-- INSERT INTO account_roles (user_id, role_id) VALUES
--     ('ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e', 1),
--     ('84d512de-df6a-4a0b-be28-a8e184bd1d6a', 2),
--     ('73c97b16-09b1-416e-94ad-f8952be14a19', 3),
--     ('ccb5418b-ac05-4f2c-8bab-6e76a51f86d9', 3);

-- INSERT INTO users (user_id, first_name, last_name, department, email) VALUES
--     ('ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e', 'Админ', '', '1f62a255-ef3a-11e5-8d88-001a64d22812', 'test1@gmail.com'),
--     ('84d512de-df6a-4a0b-be28-a8e184bd1d6a', 'Координатор', '', '1f62a256-ef3a-11e5-8d88-001a64d22812', 'test2@gmail.com'),
--     ('73c97b16-09b1-416e-94ad-f8952be14a19', 'Владимир', 'Инженер', '1f62a256-ef3a-11e5-8d88-001a64d22812', 'test3@gmail.com'),
--     ('ccb5418b-ac05-4f2c-8bab-6e76a51f86d9', 'Михаил', 'Инженер', '1f62a257-ef3a-11e5-8d88-001a64d22812', 'test4@gmail.com');

-- Insert data into the regions table
-- INSERT INTO regions (name, district) VALUES
-- ('Республика Адыгея (Адыгея)', 'Южный'),
-- ('Республика Башкортостан', 'Приволжский'),
-- ('Республика Бурятия', 'Сибирский'),
-- ('Республика Алтай', 'Сибирский'),
-- ('Республика Дагестан', 'Северо-Кавказский'),
-- ('Республика Ингушетия', 'Северо-Кавказский'),
-- ('Кабардино-Балкарская Республика', 'Северо-Кавказский'),
-- ('Республика Калмыкия', 'Южный'),
-- ('Карачаево-Черкесская Республика', 'Северо-Кавказский'),
-- ('Республика Карелия', 'Северо-Западный'),
-- ('Республика Коми', 'Северо-Западный'),
-- ('Республика Марий Эл', 'Приволжский'),
-- ('Республика Мордовия', 'Приволжский'),
-- ('Республика Саха (Якутия)', 'Дальневосточный'),
-- ('Республика Северная Осетия - Алания', 'Северо-Кавказский'),
-- ('Республика Татарстан', 'Приволжский'),
-- ('Республика Тыва', 'Сибирский'),
-- ('Удмуртская Республика', 'Приволжский'),
-- ('Республика Хакасия', 'Сибирский'),
-- ('Чеченская Республика', 'Северо-Кавказский'),
-- ('Чувашская Республика - Чувашия', 'Приволжский'),
-- ('Алтайский край', 'Сибирский'),
-- ('Краснодарский край', 'Южный'),
-- ('Красноярский край', 'Сибирский'),
-- ('Приморский край', 'Дальневосточный'),
-- ('Ставропольский край', 'Северо-Кавказский'),
-- ('Хабаровский край', 'Дальневосточный'),
-- ('Амурская область', 'Дальневосточный'),
-- ('Архангельская область', 'Северо-Западный'),
-- ('Астраханская область', 'Южный'),
-- ('Белгородская область', 'Центральный'),
-- ('Брянская область', 'Центральный'),
-- ('Владимирская область', 'Центральный'),
-- ('Волгоградская область', 'Южный'),
-- ('Вологодская область', 'Северо-Западный'),
-- ('Воронежская область', 'Центральный'),
-- ('Ивановская область', 'Центральный'),
-- ('Иркутская область', 'Сибирский'),
-- ('Калининградская область', 'Северо-Западный'),
-- ('Калужская область', 'Центральный'),
-- ('Камчатский край', 'Дальневосточный'),
-- ('Кемеровская область', 'Сибирский'),
-- ('Кировская область', 'Приволжский'),
-- ('Костромская область', 'Центральный'),
-- ('Курганская область', 'Уральский'),
-- ('Курская область', 'Центральный'),
-- ('Ленинградская область', 'Северо-Западный'),
-- ('Липецкая область', 'Центральный'),
-- ('Магаданская область', 'Дальневосточный'),
-- ('Московская область', 'Центральный'),
-- ('Мурманская область', 'Северо-Западный'),
-- ('Нижегородская область', 'Приволжский'),
-- ('Новгородская область', 'Северо-Западный'),
-- ('Новосибирская область', 'Сибирский'),
-- ('Омская область', 'Сибирский'),
-- ('Оренбургская область', 'Приволжский'),
-- ('Орловская область', 'Центральный'),
-- ('Пензенская область', 'Приволжский'),
-- ('Пермский край', 'Приволжский'),
-- ('Псковская область', 'Северо-Западный'),
-- ('Ростовская область', 'Южный'),
-- ('Рязанская область', 'Центральный'),
-- ('Самарская область', 'Приволжский'),
-- ('Саратовская область', 'Приволжский'),
-- ('Сахалинская область', 'Дальневосточный'),
-- ('Свердловская область', 'Уральский'),
-- ('Смоленская область', 'Центральный'),
-- ('Тамбовская область', 'Центральный'),
-- ('Тверская область', 'Центральный'),
-- ('Томская область', 'Сибирский'),
-- ('Тульская область', 'Центральный'),
-- ('Тюменская область', 'Уральский'),
-- ('Ульяновская область', 'Приволжский'),
-- ('Челябинская область', 'Уральский'),
-- ('Забайкальский край', 'Сибирский'),
-- ('Ярославская область', 'Центральный'),
-- ('Москва', 'Центральный'),
-- ('Санкт-Петербург', 'Северо-Западный'),
-- ('Еврейская автономная область', 'Дальневосточный'),
-- ('Ненецкий автономный округ', 'Северо-Западный'),
-- ('Ханты-Мансийский автономный округ - Югра', 'Уральский'),
-- ('Чукотский автономный округ', 'Дальневосточный'),
-- ('Ямало-Ненецкий автономный округ', 'Уральский'),
-- ('Республика Крым', 'Южный'),
-- ('Севастополь', 'Южный'),
-- ('Иные территории', 'Южный');

-- INSERT INTO clients (
--     id,
--     title,
--     region,
--     address,
--     location,
--     laboratory_system,
--     manager
-- )
-- VALUES (
--     'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', -- example UUID for client
--     'Центральная Больница г. Беломорск',
--     '77', -- Reference to Moscow region (must exist in regions table)
--     'Беломорск, Меретсковая ул., 6',
--     '{"lat": 55.7539, "lng": 37.6208}', -- JSONB location with coordinates
--     'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', -- TODO: ADD Laboraroty reference
--     '{ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e}'
-- );
--
-- INSERT INTO clients (
--     id,
--     title,
--     region,
--     address,
--     location,
--     manager
-- )
-- VALUES (
--     'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', -- example UUID for client
--     'Отделение СберБанка',
--     '77', -- Reference to Moscow region (must exist in regions table)
--     'Санкт-Петербург, Портовая ул., 89',
--     '{"lat": 55.7539, "lng": 37.6208}', -- JSONB location with coordinates
--     '{ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e}'
-- );
--
-- -- Insert a comment for a client
-- INSERT INTO comments (author_id, reference_id, text, created_at)
-- VALUES ('73c97b16-09b1-416e-94ad-f8952be14a19', '2ecc4df8-cd7a-412d-9362-09b047a67c30', 'Комментарий о проделанной работе', (NOW() AT TIME ZONE 'UTC'));
-- INSERT INTO comments (author_id, reference_id, text, created_at)
-- VALUES ('ccb5418b-ac05-4f2c-8bab-6e76a51f86d9', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Еще один комментарий', (NOW() AT TIME ZONE 'UTC'));
--
-- INSERT INTO contacts (id, name, position, phone, email, client_id)
-- VALUES ('27b1c3f2-f196-4885-8d56-9169e9f71e52', 'Вероника Васильевна', 'Заведующая лабораторией', '79992191217', 'someemail@gmail.com','a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11');
-- INSERT INTO contacts (name, phone, email, client_id)
-- VALUES ('Иван Иванович', '79992161721', 'grdandhedne@gmail.com','a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11');
--
-- INSERT INTO research_type (title) VALUES
-- ('Электрофорез белков'),
-- ('Автоматизация'),
-- ('Скрытая кровь'),
-- ('Пластик'),
-- ('Аллергология'),
-- ('Бактериология'),
-- ('Биохимия'),
-- ('Газы крови'),
-- ('Гематология'),
-- ('Группы крови'),
-- ('Иммунохимия'),
-- ('ИФА'),
-- ('Коагуалогия'),
-- ('Моча'),
-- ('ПЦР'),
-- ('СОЭ'),
-- ('Водоподготовка'),
-- ('Цитология');
--
-- INSERT INTO manufacturers (title) VALUES
-- ('СОЛТ'),
-- ('Урит'),
-- ('АО "Вектор-Бест-Балтика"'),
-- ('West Medica Produktions');
--
-- INSERT INTO classificators (id, title)
-- VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Экспресс-анализатор Triage MeterPro');
-- INSERT INTO classificators (id, title)
-- VALUES ('dff731a3-b77f-4839-a593-0345f1a5b081', 'Ultima 900');
--
-- -- Insert a new device
-- INSERT INTO devices (id, classificator, serial_number)
-- VALUES (
--     '2ecc4df8-cd7a-412d-9362-09b047a67c30',
--     'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
--     'SN123456'
-- );
-- INSERT INTO devices (classificator, serial_number, properties)
-- VALUES (
--     'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
--     'SN123453',
--     '{"manufacturer": "Company XYZ", "model": "Device Pro", "firmware": "v2.1"}'
-- );
--
-- INSERT INTO devices (id, classificator, serial_number)
-- VALUES (
--     'ddf432a3-b37f-4139-a523-2335f1a5b041',
--     'dff731a3-b77f-4839-a593-0345f1a5b081',
--     'KDP-12982'
-- );
--
INSERT INTO ticket_statuses (type, title) VALUES
('created', 'создан'),
('assigned', 'назначен'),
('inWork', 'в работе'),
('worksDone', 'работы завершены'),
('closed', 'закрыт'),
('cancelled', 'отменен');

INSERT INTO ticket_types (type, title) VALUES
('internal', 'внутренний'),
('external', 'внешний');

INSERT INTO ticket_reasons (id, title, past, present, future) VALUES
('commissioning', 'Ввод в эксплуатацию', 'Введен в эксплуатацию', 'Ввод в эксплуатацию', 'Ввести в эксплуатацию'),
('consultation', 'Консультация', 'Проведена консультация', 'Консультация', 'Провести консультацию'),
('deinstallation', 'Деинсталляция', 'Деинсталлирован', 'Деинсталляция', 'Деинсталлировать'),
('diagnostic', 'Диагностика', 'Проведена диагностика', 'Проведение диагностики', 'Провести диагностику'),
('installation', 'Инсталляция', 'Инсталлирован', 'Инсталляция', 'Инсталлировать'),
('maintenance', 'Техническое обслуживание', 'Проведено ТО', 'Проведение ТО', 'Провести ТО'),
('methodInput', 'Ввод методик', 'Введены методики', 'Ввод методик', 'Ввести методики'),
('other', 'Прочее', 'Прочее', 'Прочее', 'Прочее'),
('repair', 'Ремонт', 'Проведен ремонт', 'Ремонт', 'Отремонтировать'),
('service', 'Сервисный центр', 'Сервисный центр', 'Сервисный центр', 'Сервисный центр'),
('staffTraining', 'Обучение', 'Провести обучение', 'Обучение', 'Проведено обучение');

-- INSERT INTO tickets (client, device, ticket_type, author, assigned_by, reason, contact_person, executor, status, description, urgent, department, created_at, assigned_interval, assigned_at) VALUES
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2ecc4df8-cd7a-412d-9362-09b047a67c30', 'internal', 'ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e', '84d512de-df6a-4a0b-be28-a8e184bd1d6a', 'installation', '27b1c3f2-f196-4885-8d56-9169e9f71e52', '73c97b16-09b1-416e-94ad-f8952be14a19', 'assigned', 'Контроль прохождения 9004 ', false, '1f62a256-ef3a-11e5-8d88-001a64d22812', '2025-12-13T09:19:34.169Z', '{"start": "2025-10-15T09:19:34.169Z", "end": "2025-12-02T09:19:34.169Z"}', '2025-11-11T09:19:34.169Z');
--
-- INSERT INTO tickets (created_at, client, device, ticket_type, author, assigned_by, reason, contact_person, executor, status, description, department, assigned_interval, assigned_at) VALUES
-- ('2025-08-18T11:24:42.072Z', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2ecc4df8-cd7a-412d-9362-09b047a67c30', 'internal', 'ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e', '84d512de-df6a-4a0b-be28-a8e184bd1d6a', 'installation', '27b1c3f2-f196-4885-8d56-9169e9f71e52', 'ccb5418b-ac05-4f2c-8bab-6e76a51f86d9', 'assigned', 'Описание тикета', '1f62a257-ef3a-11e5-8d88-001a64d22812', '{"start": "2025-10-15T09:19:34.169Z", "end": "2025-12-09T09:19:34.169Z"}', '2025-11-11T09:19:34.169Z');
--
-- INSERT INTO tickets (client, device, ticket_type, author, assigned_by, reason, contact_person, executor, status, description, urgent, department, created_at, assigned_interval, assigned_at) VALUES
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'ddf432a3-b37f-4139-a523-2335f1a5b041', 'internal', 'ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e', '84d512de-df6a-4a0b-be28-a8e184bd1d6a', 'diagnostic', '27b1c3f2-f196-4885-8d56-9169e9f71e52', '73c97b16-09b1-416e-94ad-f8952be14a19', 'assigned', 'Выдаёт ошибку холостой пробы, превышение предела RBC. При выполнении анализов не считает эритроциты.', true, '1f62a256-ef3a-11e5-8d88-001a64d22812', '2025-11-13T09:19:34.169Z', '{"start": "2025-10-15T09:19:34.169Z", "end": "2025-12-12T09:19:34.169Z"}', '2025-11-11T09:19:34.169Z');
--
-- -- One in created phase 
-- INSERT INTO tickets (client, device, ticket_type, author, reason, contact_person, status, description, urgent, department, created_at, assigned_interval, executor) VALUES
-- ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'ddf432a3-b37f-4139-a523-2335f1a5b041', 'internal', 'ad9fa963-cad8-4bc3-b8e2-f4a4f70cf95e','diagnostic', '27b1c3f2-f196-4885-8d56-9169e9f71e52', 'created', 'Выдаёт ошибку холостой пробы, превышение предела RBC. При выполнении анализов не считает эритроциты.', true, '1f62a256-ef3a-11e5-8d88-001a64d22812', '2025-11-13T09:19:34.169Z', '{"start": "2025-10-15T09:19:34.169Z", "end": "2025-12-10T09:19:34.169Z"}', '84d512de-df6a-4a0b-be28-a8e184bd1d6a');
--
-- INSERT INTO agreements (id, actual_client, distributor, device, assigned_at, type) VALUES
--     (gen_random_uuid(), 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2ecc4df8-cd7a-412d-9362-09b047a67c30', NOW(), 'rent');
-- INSERT INTO agreements (id, actual_client, distributor, device, assigned_at, type) VALUES
--     (gen_random_uuid(), 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2ecc4df8-cd7a-412d-9362-09b047a67c30', NOW(), 'rent');
-- INSERT INTO agreements (id, actual_client, distributor, assigned_at, type) VALUES
--     (gen_random_uuid(), 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', NOW(), 'buy');
--
-- INSERT INTO ra_options (id, title) VALUES
--     ('95e0f8d9-b497-4d38-844d-10e5746f3aa1', 'RDP'),
--     ('51dfb946-b0c9-4503-9fbe-96ea5329095f', 'AmmyAdmin'),
--     ('fd23c3d9-2a55-45db-8f74-a1f7af9769b5', 'TeamViewer');
--
-- INSERT INTO remote_access (device_id, parameter_id) VALUES
--     ('2ecc4df8-cd7a-412d-9362-09b047a67c30', '95e0f8d9-b497-4d38-844d-10e5746f3aa1'),
--     ('2ecc4df8-cd7a-412d-9362-09b047a67c30', '51dfb946-b0c9-4503-9fbe-96ea5329095f'),
--     ('2ecc4df8-cd7a-412d-9362-09b047a67c30', 'fd23c3d9-2a55-45db-8f74-a1f7af9769b5');

COMMIT;
