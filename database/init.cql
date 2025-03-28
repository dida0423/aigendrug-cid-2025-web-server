CREATE KEYSPACE ks_admin WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1 };

USE ks_admin;

CREATE TABLE sessions (
    id uuid PRIMARY KEY,
    name text,
    status text,
    tool_status text,
    assigned_tool_id uuid,
    created_at timestamp
);

CREATE TABLE chat_messages (
    session_id uuid,
    id uuid,
    role text,
    message text,
    created_at timestamp,
    message_type int,
    linked_tool_ids list<uuid>,

    PRIMARY KEY (session_id, created_at)
) WITH CLUSTERING ORDER BY (created_at ASC);

CREATE TABLE tools (
    id uuid PRIMARY KEY,
    name text,
    version text,
    description text,
    provider_interface text,
    created_at timestamp
);

CREATE TABLE tool_messages (
    id uuid,
    session_id uuid,
    tool_id uuid,
    role text,
    data text,
    created_at timestamp,

    PRIMARY KEY (session_id, created_at)
) WITH CLUSTERING ORDER BY (created_at ASC);