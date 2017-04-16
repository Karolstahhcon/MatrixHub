CREATE TABLE IF NOT EXISTS registries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    type TEXT NOT NULL,
    url TEXT NOT NULL,
    credential_type TEXT DEFAULT NULL,
    auth_info TEXT,
    insecure INTEGER NOT NULL DEFAULT 0 CHECK (insecure IN (0, 1)),
    status INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_registries_name ON registries (name);

CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT COLLATE BINARY DEFAULT '',
    type INTEGER NOT NULL DEFAULT 0,
    registry_id INTEGER DEFAULT NULL,
    organization TEXT DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_projects_registry_id
        FOREIGN KEY (registry_id) REFERENCES registries (id)
);
CREATE INDEX idx_projects_name ON projects (name);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT COLLATE NOCASE NOT NULL,
    password TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL DEFAULT '',
    session_version INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uniq_users_username UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    permissions TEXT NOT NULL,
    scope TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_roles_name ON roles (name);

CREATE TABLE IF NOT EXISTS members_roles_projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    member_id INTEGER NOT NULL,
    member_type TEXT COLLATE NOCASE NOT NULL,
    role_id INTEGER DEFAULT NULL,
    project_id INTEGER DEFAULT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uniq_members_roles_projects_project_member
        UNIQUE (project_id, member_id, member_type),
    CONSTRAINT fk_members_roles_projects_project_id
        FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE,
    CONSTRAINT fk_members_roles_projects_role_id
        FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE
);
CREATE INDEX idx_members_roles_projects_project_id
    ON members_roles_projects (project_id);
CREATE INDEX idx_members_roles_projects_member_id
    ON members_roles_projects (member_id);

CREATE TABLE IF NOT EXISTS models (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT COLLATE NOCASE NOT NULL,
    project_id INTEGER NOT NULL,
    size INTEGER NOT NULL,
    default_branch TEXT NOT NULL,
    parameter_count INTEGER NOT NULL,
    readme_content TEXT NOT NULL,
    is_popular INTEGER NOT NULL DEFAULT 0 CHECK (is_popular IN (0, 1)),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uniq_models_project_name UNIQUE (project_id, name),
    CONSTRAINT fk_models_project_id
        FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);
CREATE INDEX idx_models_updated_at ON models (updated_at);

CREATE TABLE IF NOT EXISTS labels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT COLLATE NOCASE NOT NULL,
    category TEXT COLLATE NOCASE NOT NULL,
    scope TEXT COLLATE NOCASE NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uniq_labels_name_category_scope UNIQUE (name, category, scope)
);

CREATE TABLE IF NOT EXISTS models_labels (
    model_id INTEGER NOT NULL,
    label_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (model_id, label_id),
    CONSTRAINT fk_models_labels_model_id
        FOREIGN KEY (model_id) REFERENCES models (id) ON DELETE CASCADE,
    CONSTRAINT fk_models_labels_label_id
        FOREIGN KEY (label_id) REFERENCES labels (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS datasets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT COLLATE NOCASE NOT NULL,
    project_id INTEGER NOT NULL,
    default_branch TEXT NOT NULL DEFAULT 'main',
    is_popular INTEGER NOT NULL DEFAULT 0 CHECK (is_popular IN (0, 1)),
    num_rows TEXT NOT NULL DEFAULT '',
    size INTEGER NOT NULL DEFAULT 0,
    readme_content TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uniq_datasets_project_name UNIQUE (project_id, name),
    CONSTRAINT fk_datasets_project_id
        FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);
CREATE INDEX idx_datasets_updated_at ON datasets (updated_at);

CREATE TABLE IF NOT EXISTS datasets_labels (
    dataset_id INTEGER NOT NULL,
    label_id INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (dataset_id, label_id),
    CONSTRAINT fk_datasets_labels_dataset_id
        FOREIGN KEY (dataset_id) REFERENCES datasets (id) ON DELETE CASCADE,
    CONSTRAINT fk_datasets_labels_label_id
        FOREIGN KEY (label_id) REFERENCES labels (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS access_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    token_hash TEXT COLLATE NOCASE NOT NULL DEFAULT '',
    enabled INTEGER NOT NULL DEFAULT 0 CHECK (enabled IN (0, 1)),
    expire_at DATETIME DEFAULT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uniq_access_tokens_token_hash UNIQUE (token_hash)
);
CREATE INDEX idx_access_tokens_user_id ON access_tokens (user_id);

CREATE TABLE IF NOT EXISTS ssh_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    public_key TEXT NOT NULL,
    fingerprint TEXT COLLATE NOCASE NOT NULL,
    expire_at DATETIME DEFAULT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uniq_ssh_keys_fingerprint UNIQUE (fingerprint)
);
CREATE INDEX idx_ssh_keys_user_id ON ssh_keys (user_id);

CREATE TABLE IF NOT EXISTS sync_policies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
