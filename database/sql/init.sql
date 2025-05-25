CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    name TEXT,
    status TEXT,
    tool_status TEXT,
    assigned_tool_id UUID,
    created_at TIMESTAMP
);

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

CREATE TABLE tools (
    id UUID PRIMARY KEY,
    name TEXT,
    version TEXT,
    description TEXT,
    provider_interface TEXT,
    created_at TIMESTAMP
);

CREATE TABLE tool_messages (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL,
    tool_id UUID,
    role TEXT,
    data TEXT,
    created_at TIMESTAMP,
    CONSTRAINT fk_session_tool FOREIGN KEY (session_id) REFERENCES sessions(id)
);