-- idp stores generic identity provider.
CREATE TABLE idp (
  id serial PRIMARY KEY,
  resource_id text NOT NULL,
  name text NOT NULL,
  domain text NOT NULL,
  type text NOT NULL CONSTRAINT idp_type_check CHECK (type IN ('OAUTH2', 'OIDC', 'LDAP')),
  -- config stores the corresponding configuration of the IdP, which may vary depending on the type of the IdP.
  config jsonb NOT NULL DEFAULT '{}'
);

CREATE UNIQUE INDEX idx_idp_unique_resource_id ON idp(resource_id);

ALTER SEQUENCE idp_id_seq RESTART WITH 101;

-- principal
CREATE TABLE principal (
    id serial PRIMARY KEY,
    deleted boolean NOT NULL DEFAULT FALSE,
    created_at timestamptz NOT NULL DEFAULT now(),
    type text NOT NULL CHECK (type IN ('END_USER', 'SYSTEM_BOT', 'SERVICE_ACCOUNT')),
    name text NOT NULL,
    email text NOT NULL,
    password_hash text NOT NULL,
    phone text NOT NULL DEFAULT '',
    mfa_config jsonb NOT NULL DEFAULT '{}',
    profile jsonb NOT NULL DEFAULT '{}'
);

-- Setting
CREATE TABLE setting (
    id serial PRIMARY KEY,
    name text NOT NULL,
    value text NOT NULL
);

CREATE UNIQUE INDEX idx_setting_unique_name ON setting(name);

ALTER SEQUENCE setting_id_seq RESTART WITH 101;

-- Role
CREATE TABLE role (
    id bigserial PRIMARY KEY,
    resource_id text NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    permissions jsonb NOT NULL DEFAULT '{}',
    payload jsonb NOT NULL DEFAULT '{}' -- saved for future use
);

CREATE UNIQUE INDEX idx_role_unique_resource_id on role (resource_id);

ALTER SEQUENCE role_id_seq RESTART WITH 101;

-- Policy
-- policy stores the policies for each resources.
CREATE TABLE policy (
    id serial PRIMARY KEY,
    enforce boolean NOT NULL DEFAULT TRUE,
    updated_at timestamptz NOT NULL DEFAULT now(),
    resource_type text NOT NULL,
    resource TEXT NOT NULL,
    type text NOT NULL,
    payload jsonb NOT NULL DEFAULT '{}',
    inherit_from_parent boolean NOT NULL DEFAULT TRUE
);

CREATE UNIQUE INDEX idx_policy_unique_resource_type_resource_type ON policy(resource_type, resource, type);

ALTER SEQUENCE policy_id_seq RESTART WITH 101;

-- Project
CREATE TABLE project (
    id serial PRIMARY KEY,
    deleted boolean NOT NULL DEFAULT FALSE,
    name text NOT NULL,
    resource_id text NOT NULL,
    data_classification_config_id text NOT NULL DEFAULT '',
    setting jsonb NOT NULL DEFAULT '{}'
);

CREATE UNIQUE INDEX idx_project_unique_resource_id ON project(resource_id);

-- Project Hook
CREATE TABLE project_webhook (
    id serial PRIMARY KEY,
    project text NOT NULL REFERENCES project(resource_id),
    type text NOT NULL CHECK (type LIKE 'bb.plugin.webhook.%'),
    name text NOT NULL,
    url text NOT NULL,
    event_list text ARRAY NOT NULL,
    payload jsonb NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_project_webhook_project ON project_webhook(project);

ALTER SEQUENCE project_webhook_id_seq RESTART WITH 101;

-- Instance
CREATE TABLE instance (
    id serial PRIMARY KEY,
    deleted boolean NOT NULL DEFAULT FALSE,
    environment text,
    resource_id text NOT NULL,
    metadata jsonb NOT NULL DEFAULT '{}'
);

CREATE UNIQUE INDEX idx_instance_unique_resource_id ON instance(resource_id);

ALTER SEQUENCE instance_id_seq RESTART WITH 101;

-- db stores the databases for a particular instance
-- data is synced periodically from the instance
CREATE TABLE db (
    id serial PRIMARY KEY,
    deleted boolean NOT NULL DEFAULT FALSE,
    project text NOT NULL REFERENCES project(resource_id),
    instance text NOT NULL REFERENCES instance(resource_id),
    name text NOT NULL,
    environment text,
    metadata jsonb NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_db_project ON db(project);

CREATE UNIQUE INDEX idx_db_unique_instance_name ON db(instance, name);

ALTER SEQUENCE db_id_seq RESTART WITH 101;

-- db_schema stores the database schema metadata for a particular database.
CREATE TABLE db_schema (
    id serial PRIMARY KEY,
    instance text NOT NULL,
    db_name text NOT NULL,
    metadata json NOT NULL DEFAULT '{}',
    raw_dump text NOT NULL DEFAULT '',
    config jsonb NOT NULL DEFAULT '{}',
    todo boolean NOT NULL DEFAULT TRUE,
    CONSTRAINT db_schema_instance_db_name_fkey FOREIGN KEY(instance, db_name) REFERENCES db(instance, name)
);

CREATE UNIQUE INDEX idx_db_schema_unique_instance_db_name ON db_schema(instance, db_name);

ALTER SEQUENCE db_schema_id_seq RESTART WITH 101;

-- Deprecated. To be deleted later.
CREATE TABLE data_source (
    id serial PRIMARY KEY,
    instance text NOT NULL REFERENCES instance(resource_id),
    options jsonb NOT NULL DEFAULT '{}'
);

ALTER SEQUENCE data_source_id_seq RESTART WITH 101;

CREATE TABLE sheet_blob (
	sha256 bytea NOT NULL PRIMARY KEY,
	content text NOT NULL
);

-- sheet table stores general statements.
CREATE TABLE sheet (
    id serial PRIMARY KEY,
    creator_id integer NOT NULL REFERENCES principal(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    project text NOT NULL REFERENCES project(resource_id),
    name text NOT NULL,
    sha256 bytea NOT NULL,
    payload jsonb NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_sheet_project ON sheet(project);

ALTER SEQUENCE sheet_id_seq RESTART WITH 101;

-----------------------
-- Pipeline related BEGIN
-- pipeline table
CREATE TABLE pipeline (
    id serial PRIMARY KEY,
    creator_id integer NOT NULL REFERENCES principal(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    project text NOT NULL REFERENCES project(resource_id)
);

ALTER SEQUENCE pipeline_id_seq RESTART WITH 101;

-- task table stores the task for the pipeline
CREATE TABLE task (
    id serial PRIMARY KEY,
    pipeline_id integer NOT NULL REFERENCES pipeline(id),
    instance text NOT NULL REFERENCES instance(resource_id),
    environment text,
    db_name text,
    type text NOT NULL,
    payload jsonb NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_task_pipeline_id_environment ON task(pipeline_id, environment);

ALTER SEQUENCE task_id_seq RESTART WITH 101;

-- task run table stores the task run
CREATE TABLE task_run (
    id serial PRIMARY KEY,
    creator_id integer NOT NULL REFERENCES principal(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    task_id integer NOT NULL REFERENCES task(id),
    sheet_id integer REFERENCES sheet(id),
    attempt integer NOT NULL,
    status text NOT NULL CHECK (status IN ('PENDING', 'RUNNING', 'DONE', 'FAILED', 'CANCELED')),
    started_at timestamptz NULL,
    run_at timestamptz,
    code integer NOT NULL DEFAULT 0,
    -- result saves the task run result in json format
    result jsonb NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_task_run_task_id ON task_run(task_id);

CREATE UNIQUE INDEX uk_task_run_task_id_attempt ON task_run (task_id, attempt);

ALTER SEQUENCE task_run_id_seq RESTART WITH 101;

CREATE TABLE task_run_log (
    id bigserial PRIMARY KEY,
    task_run_id integer NOT NULL REFERENCES task_run(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    payload jsonb NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_task_run_log_task_run_id ON task_run_log(task_run_id);

ALTER SEQUENCE task_run_log_id_seq RESTART WITH 101;

-- Pipeline related END
-----------------------
-- Plan related BEGIN
CREATE TABLE plan (
    id bigserial PRIMARY KEY,
    deleted boolean NOT NULL DEFAULT FALSE,
    creator_id integer NOT NULL REFERENCES principal(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    project text NOT NULL REFERENCES project(resource_id),
    pipeline_id integer REFERENCES pipeline(id),
    name text NOT NULL,
    description text NOT NULL,
    config jsonb NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_plan_project ON plan(project);

CREATE INDEX idx_plan_pipeline_id ON plan(pipeline_id);

ALTER SEQUENCE plan_id_seq RESTART WITH 101;

CREATE TABLE plan_check_run (
    id serial PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    plan_id bigint NOT NULL REFERENCES plan(id),
    status text NOT NULL CHECK (status IN ('RUNNING', 'DONE', 'FAILED', 'CANCELED')),
    type text NOT NULL CHECK (type LIKE 'bb.plan-check.%'),
    config jsonb NOT NULL DEFAULT '{}',
    result jsonb NOT NULL DEFAULT '{}',
    payload jsonb NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_plan_check_run_plan_id ON plan_check_run (plan_id);

ALTER SEQUENCE plan_check_run_id_seq RESTART WITH 101;

-- Plan related END
-----------------------
-- issue
CREATE TABLE issue (
    id serial PRIMARY KEY,
    creator_id integer NOT NULL REFERENCES principal(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    project text NOT NULL REFERENCES project(resource_id),
    plan_id bigint REFERENCES plan(id),
    pipeline_id integer REFERENCES pipeline(id),
    name text NOT NULL,
    status text NOT NULL CHECK (status IN ('OPEN', 'DONE', 'CANCELED')),
    type text NOT NULL,
    description text NOT NULL DEFAULT '',
    payload jsonb NOT NULL DEFAULT '{}',
    ts_vector tsvector
);

CREATE INDEX idx_issue_project ON issue(project);

CREATE INDEX idx_issue_plan_id ON issue(plan_id);

CREATE INDEX idx_issue_pipeline_id ON issue(pipeline_id);

CREATE INDEX idx_issue_creator_id ON issue(creator_id);

CREATE INDEX idx_issue_ts_vector ON issue USING GIN(ts_vector);

ALTER SEQUENCE issue_id_seq RESTART WITH 101;

-- instance change history records the changes an instance and its databases.
CREATE TABLE instance_change_history (
    id bigserial PRIMARY KEY,
    version text NOT NULL
);

CREATE UNIQUE INDEX idx_instance_change_history_unique_version ON instance_change_history (version);

ALTER SEQUENCE instance_change_history_id_seq RESTART WITH 101;

CREATE TABLE audit_log (
    id bigserial PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    payload jsonb NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_audit_log_created_at ON audit_log(created_at);

CREATE INDEX idx_audit_log_payload_parent ON audit_log((payload->>'parent'));

CREATE INDEX idx_audit_log_payload_method ON audit_log((payload->>'method'));

CREATE INDEX idx_audit_log_payload_resource ON audit_log((payload->>'resource'));

CREATE INDEX idx_audit_log_payload_user ON audit_log((payload->>'user'));

ALTER SEQUENCE audit_log_id_seq RESTART WITH 101;

CREATE TABLE issue_comment (
    id bigserial PRIMARY KEY,
    creator_id integer NOT NULL REFERENCES principal(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    issue_id integer NOT NULL REFERENCES issue(id),
    payload jsonb NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_issue_comment_issue_id ON issue_comment(issue_id);

ALTER SEQUENCE issue_comment_id_seq RESTART WITH 101;

CREATE TABLE query_history (
    id bigserial PRIMARY KEY,
    creator_id integer NOT NULL REFERENCES principal(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    project_id text NOT NULL, -- the project resource id
    database text NOT NULL, -- the database resource name, for example, instances/{instance}/databases/{database}
    statement text NOT NULL,
    type text NOT NULL, -- the history type, support QUERY and EXPORT.
    payload jsonb NOT NULL DEFAULT '{}' -- saved for details, like error, duration, etc.
);

CREATE INDEX idx_query_history_creator_id_created_at_project_id ON query_history(creator_id, created_at, project_id DESC);

ALTER SEQUENCE query_history_id_seq RESTART WITH 101;

-- worksheet table stores worksheets in SQL Editor.
CREATE TABLE worksheet (
    id serial PRIMARY KEY,
    creator_id integer NOT NULL REFERENCES principal(id),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    project text NOT NULL REFERENCES project(resource_id),
    instance text,
    db_name text,
    name text NOT NULL,
    statement text NOT NULL,
    visibility text NOT NULL,
    payload jsonb NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_worksheet_creator_id_project ON worksheet(creator_id, project);

ALTER SEQUENCE worksheet_id_seq RESTART WITH 101;

-- worksheet_organizer table stores the sheet status for a principal.
CREATE TABLE worksheet_organizer (
    id serial PRIMARY KEY,
    worksheet_id integer NOT NULL REFERENCES worksheet(id) ON DELETE CASCADE,
    principal_id integer NOT NULL REFERENCES principal(id),
    starred boolean NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX idx_worksheet_organizer_unique_sheet_id_principal_id ON worksheet_organizer(worksheet_id, principal_id);

CREATE INDEX idx_worksheet_organizer_principal_id ON worksheet_organizer(principal_id);

-- risk stores the definition of a risk.
CREATE TABLE risk (
    id bigserial PRIMARY KEY,
    source text NOT NULL CHECK (source LIKE 'bb.risk.%'),
    -- how risky is the risk, the higher the riskier
    level bigint NOT NULL,
    name text NOT NULL,
    active boolean NOT NULL,
    expression jsonb NOT NULL
);

ALTER SEQUENCE risk_id_seq RESTART WITH 101;

CREATE TABLE db_group (
    id bigserial PRIMARY KEY,
    project text NOT NULL REFERENCES project(resource_id),
    resource_id text NOT NULL,
    placeholder text NOT NULL DEFAULT '',
    expression jsonb NOT NULL DEFAULT '{}',
    payload jsonb NOT NULL DEFAULT '{}'
);

CREATE UNIQUE INDEX idx_db_group_unique_project_resource_id ON db_group(project, resource_id);

CREATE UNIQUE INDEX idx_db_group_unique_project_placeholder ON db_group(project, placeholder);

ALTER SEQUENCE db_group_id_seq RESTART WITH 101;

-- changelist table stores project changelists.
CREATE TABLE changelist (
    id serial PRIMARY KEY,
    creator_id integer NOT NULL REFERENCES principal (id),
    updated_at timestamptz NOT NULL DEFAULT now(),
    project text NOT NULL REFERENCES project(resource_id),
    name text NOT NULL,
    payload jsonb NOT NULL DEFAULT '{}'
);

CREATE UNIQUE INDEX idx_changelist_project_name ON changelist(project, name);

ALTER SEQUENCE changelist_id_seq RESTART WITH 101;

CREATE TABLE export_archive (
  id serial PRIMARY KEY,
  created_at timestamptz NOT NULL DEFAULT now(),
  bytes bytea,
  payload jsonb NOT NULL DEFAULT '{}'
);

CREATE TABLE user_group (
  email text PRIMARY KEY,
  name text NOT NULL,
  description text NOT NULL DEFAULT '',
  payload jsonb NOT NULL DEFAULT '{}'
);

-- review config table.
CREATE TABLE review_config (
    id text NOT NULL PRIMARY KEY,
    enabled boolean NOT NULL DEFAULT TRUE,
    name text NOT NULL,
    payload jsonb NOT NULL DEFAULT '{}'
);

CREATE TABLE revision (
    id bigserial PRIMARY KEY,
    instance text NOT NULL,
    db_name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    deleter_id integer REFERENCES principal(id),
    deleted_at timestamptz,
    version text NOT NULL,
    payload jsonb NOT NULL DEFAULT '{}',
    CONSTRAINT revision_instance_db_name_fkey FOREIGN KEY(instance, db_name) REFERENCES db(instance, name)
);

ALTER SEQUENCE revision_id_seq RESTART WITH 101;

CREATE UNIQUE INDEX idx_revision_unique_instance_db_name_version_deleted_at_null ON revision(instance, db_name, version) WHERE deleted_at IS NULL;

CREATE INDEX idx_revision_instance_db_name_version ON revision(instance, db_name, version);

CREATE TABLE sync_history (
    id bigserial PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    instance text NOT NULL,
    db_name text NOT NULL,
    metadata json NOT NULL DEFAULT '{}',
    raw_dump text NOT NULL DEFAULT '',
    CONSTRAINT sync_history_instance_db_name_fkey FOREIGN KEY(instance, db_name) REFERENCES db(instance, name)
);

ALTER SEQUENCE sync_history_id_seq RESTART WITH 101;

CREATE INDEX idx_sync_history_instance_db_name_created_at ON sync_history (instance, db_name, created_at);

CREATE TABLE changelog (
    id bigserial PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    instance text NOT NULL,
    db_name text NOT NULL,
    status text NOT NULL CONSTRAINT changelog_status_check CHECK (status IN ('PENDING', 'DONE', 'FAILED')),
    prev_sync_history_id bigint REFERENCES sync_history(id),
    sync_history_id bigint REFERENCES sync_history(id),
    payload jsonb NOT NULL DEFAULT '{}',
    CONSTRAINT changelog_instance_db_name_fkey FOREIGN KEY(instance, db_name) REFERENCES db(instance, name)
);

ALTER SEQUENCE changelog_id_seq RESTART WITH 101;

CREATE INDEX idx_changelog_instance_db_name ON changelog (instance, db_name);

CREATE TABLE release (
    id bigserial PRIMARY KEY,
    deleted boolean NOT NULL DEFAULT FALSE,
    project text NOT NULL REFERENCES project(resource_id),
    creator_id integer NOT NULL REFERENCES principal (id),
    created_at timestamptz NOT NULL DEFAULT now(),
    digest text NOT NULL DEFAULT '',
    payload jsonb NOT NULL DEFAULT '{}'
);

ALTER SEQUENCE release_id_seq RESTART WITH 101;

CREATE INDEX idx_release_project ON release(project);

CREATE UNIQUE INDEX idx_release_unique_project_digest ON release(project, digest) WHERE digest != '';


-- Default bytebase system account id is 1.
INSERT INTO principal (id, type, name, email, password_hash) VALUES (1, 'SYSTEM_BOT', 'Bytebase', 'support@bytebase.com', '');

ALTER SEQUENCE principal_id_seq RESTART WITH 101;

-- Default project.
INSERT INTO project (id, name, resource_id) VALUES (1, 'Default', 'default');

ALTER SEQUENCE project_id_seq RESTART WITH 101;
