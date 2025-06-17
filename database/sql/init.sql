-- Create schema
CREATE SCHEMA IF NOT EXISTS ks_admin;

-- Set search path
SET search_path TO ks_admin;

-- Create sessions table
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    name TEXT,
    status TEXT,
    tool_status TEXT,
    assigned_tool_id UUID,
    created_at TIMESTAMP
);

-- Create chat_messages table
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL,
    role TEXT,
    message TEXT,
    created_at TIMESTAMP,
    message_type INTEGER,
    linked_tool_ids UUID[],
    CONSTRAINT fk_session FOREIGN KEY (session_id) REFERENCES sessions(id)
);

-- Create tools table
CREATE TABLE tools (
    id UUID PRIMARY KEY,
    name TEXT,
    version TEXT,
    description TEXT,
    provider_interface TEXT,
    created_at TIMESTAMP
);

-- Create tool_messages table
CREATE TABLE tool_messages (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL,
    tool_id UUID,
    role TEXT,
    data TEXT,
    created_at TIMESTAMP,
    CONSTRAINT fk_session_tool FOREIGN KEY (session_id) REFERENCES sessions(id)
);

-- Create indexes for better query performance
CREATE INDEX idx_chat_messages_session_id ON chat_messages(session_id);
CREATE INDEX idx_tool_messages_session_id ON tool_messages(session_id);
CREATE INDEX idx_tool_messages_tool_id ON tool_messages(tool_id);

-- Create ordered indexes to replace CLUSTERING ORDER BY
CREATE INDEX idx_chat_messages_created_at_asc ON chat_messages(session_id, created_at ASC);
CREATE INDEX idx_tool_messages_created_at_asc ON tool_messages(session_id, created_at ASC);