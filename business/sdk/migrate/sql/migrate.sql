-- Version: 1.01
-- Description: Create table users
CREATE TABLE users (
    user_id UUID NOT NULL,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    roles TEXT [] NOT NULL,
    password_hash TEXT NOT NULL,
    department TEXT NULL,
    enabled BOOLEAN NOT NULL,
    date_created TIMESTAMP NOT NULL,
    date_updated TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id)
);

CREATE TABLE conversations (
    conversation_id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(user_id)
);

CREATE TABLE messages (
    message_id UUID NOT NULL,
    conversation_id UUID NOT NULL REFERENCES conversations(conversation_id),
    user_id UUID NOT NULL REFERENCES users(user_id),
    "role" TEXT,
    content TEXT,
    "order" INT,
    PRIMARY KEY (message_id, conversation_id)
);