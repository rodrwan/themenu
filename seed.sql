-- Enable required extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

INSERT INTO users (id, name, email) VALUES ('00000000-0000-0000-0000-000000000001', 'Usuario de Prueba', 'prueba@correo.com') ON CONFLICT (id) DO NOTHING;

INSERT INTO user_roles (user_id, role_id) VALUES ('00000000-0000-0000-0000-000000000001', (SELECT id FROM roles WHERE name = 'client')) ON CONFLICT DO NOTHING;


-- Seed roles
INSERT INTO roles (id, name) VALUES
  ('416af891-6368-46d1-9129-1de0b57bdd16', 'client'),
  ('054ab879-a1ad-408e-819d-cf6d1c614af8', 'kitchen'),
  ('670202cd-3983-4ed8-9fc3-029bff2bdc56', 'admin');

-- Seed permissions
INSERT INTO permissions (id, name) VALUES
  ('920f02df-2db9-4b49-8a5a-4c75e0ca083b', 'view_menu'),
  ('d8a32f2f-d25f-4335-8a01-01b069aeee5d', 'place_order'),
  ('9536d456-2b3b-4df5-8082-41fe8f006da4', 'update_order_status'),
  ('b20bf9f7-be0e-4bc7-a1c9-9f5d3facdeaa', 'manage_dishes'),
  ('dcba757a-2479-48f1-b97e-4c1c8c3f12ae', 'manage_users');

-- Assign permissions to client role
INSERT INTO role_permissions (role_id, permission_id) VALUES
  ('416af891-6368-46d1-9129-1de0b57bdd16', '920f02df-2db9-4b49-8a5a-4c75e0ca083b'),
  ('416af891-6368-46d1-9129-1de0b57bdd16', 'd8a32f2f-d25f-4335-8a01-01b069aeee5d');

-- Assign permissions to kitchen role
INSERT INTO role_permissions (role_id, permission_id) VALUES
  ('054ab879-a1ad-408e-819d-cf6d1c614af8', '9536d456-2b3b-4df5-8082-41fe8f006da4');

-- Assign permissions to admin role
INSERT INTO role_permissions (role_id, permission_id) VALUES
  ('670202cd-3983-4ed8-9fc3-029bff2bdc56', '920f02df-2db9-4b49-8a5a-4c75e0ca083b'),
  ('670202cd-3983-4ed8-9fc3-029bff2bdc56', 'd8a32f2f-d25f-4335-8a01-01b069aeee5d'),
  ('670202cd-3983-4ed8-9fc3-029bff2bdc56', '9536d456-2b3b-4df5-8082-41fe8f006da4'),
  ('670202cd-3983-4ed8-9fc3-029bff2bdc56', 'b20bf9f7-be0e-4bc7-a1c9-9f5d3facdeaa'),
  ('670202cd-3983-4ed8-9fc3-029bff2bdc56', 'dcba757a-2479-48f1-b97e-4c1c8c3f12ae');

-- Insertar roles
INSERT INTO roles (id, name) VALUES
    ('11111111-1111-1111-1111-111111111111', 'client') ON CONFLICT DO NOTHING;

-- Insertar usuario de prueba
INSERT INTO users (id, name, email) VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Usuario Cliente', 'cliente@test.com') ON CONFLICT DO NOTHING;

-- Asignar rol al usuario
INSERT INTO user_roles (user_id, role_id) VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111') ON CONFLICT DO NOTHING;

-- Insertar platos de ejemplo para hoy
INSERT INTO dishes (id, name, description, price, prep_time_minutes, available_on) VALUES
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', 'Pasta Carbonara', 'Pasta con salsa cremosa, panceta y queso parmesano', 15.99, 20, CURRENT_DATE),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'Ensalada César', 'Lechuga romana, pollo a la parrilla, crutones y aderezo césar', 12.99, 15, CURRENT_DATE),
    ('ffffffff-ffff-ffff-ffff-ffffffffffff', 'Pizza Margherita', 'Pizza con salsa de tomate, mozzarella y albahaca', 14.99, 25, CURRENT_DATE)
ON CONFLICT DO NOTHING;

-- name: SeedDishes :exec
INSERT INTO dishes (id, name, description, price, prep_time_minutes, available_on)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'Ensalada César', 'Lechuga romana, crutones, parmesano y aderezo césar', 12.99, 15, CURRENT_DATE),
    ('22222222-2222-2222-2222-222222222222', 'Pasta Carbonara', 'Espaguetis con salsa cremosa, panceta y parmesano', 16.99, 20, CURRENT_DATE),
    ('33333333-3333-3333-3333-333333333333', 'Filete de Salmón', 'Salmón a la parrilla con vegetales asados', 24.99, 25, CURRENT_DATE),
    ('44444444-4444-4444-4444-444444444444', 'Tiramisú', 'Postre italiano con café, mascarpone y cacao', 8.99, 10, CURRENT_DATE);