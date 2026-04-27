-- my-reminders-backend/setup.sql

-- Usamos UUIDs para claves primarias (más seguro y modular)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tabla de Usuarios Actualizada (¡Con password_hash!)
CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                     first_name TEXT NOT NULL,
                                     last_name TEXT NOT NULL,
                                     country_code CHAR(2), -- Ej. 'ES', 'MX'
                                     email TEXT UNIQUE NOT NULL,
                                     password_hash TEXT NOT NULL, -- Aquí guardaremos el hash de la contraseña, no el texto plano
                                     biometrics_enabled BOOLEAN DEFAULT FALSE, -- Para el flujo que diseñamos
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de Grupos
CREATE TABLE IF NOT EXISTS groups (
                                      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                      name TEXT UNIQUE NOT NULL, -- Ej. 'Hogar', 'Oficina'
                                      admin_id UUID REFERENCES users(id), -- El creador/admin inicial
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de Membresía de Grupo (Maneja los 3 estados: PENDING, APPROVED, ADMIN)
CREATE TABLE IF NOT EXISTS group_members (
                                             user_id UUID REFERENCES users(id),
                                             group_id UUID REFERENCES groups(id),
                                             status TEXT DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'ADMIN')),
                                             PRIMARY KEY (user_id, group_id)
);

-- Tabla de Recordatorios (La que ya diseñamos)
CREATE TABLE IF NOT EXISTS reminders (
                                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                         title TEXT NOT NULL,
                                         description TEXT,
                                         due_date TIMESTAMP,
                                         priority TEXT DEFAULT 'MEDIUM' CHECK (priority IN ('HIGH', 'MEDIUM', 'LOW')),
                                         color_accent CHAR(7), -- Ej. '#10B981'
                                         category TEXT, -- Ej. 'Trabajo', 'Casa'

                                         user_id UUID REFERENCES users(id), -- Quién lo creó
                                         group_id UUID REFERENCES groups(id), -- A qué grupo pertenece (NULL si es individual)

                                         is_completed BOOLEAN DEFAULT FALSE,
                                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla de Notificaciones
CREATE TABLE IF NOT EXISTS notifications (
                                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                             user_id UUID REFERENCES users(id), -- El receptor (el ADMIN)
                                             requester_id UUID REFERENCES users(id), -- Quién solicita entrar
                                             group_id UUID REFERENCES groups(id), -- A qué grupo quiere entrar
                                             type TEXT NOT NULL, -- 'GROUP_REQUEST', 'URGENT_TASK'
                                             title TEXT NOT NULL,
                                             message TEXT NOT NULL,
                                             is_read BOOLEAN DEFAULT FALSE,
                                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);