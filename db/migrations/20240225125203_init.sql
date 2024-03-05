-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
    id BIGINT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS "chat" (
    id BIGSERIAL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS chat_user (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    chat_id BIGINT,
    CONSTRAINT fk_chat_id FOREIGN KEY (chat_id) REFERENCES chat (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS message (
    id BIGSERIAL PRIMARY KEY,
    from_user_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,
    text TEXT NOT NULL,
    sent_at timestamp with time zone NOT NULL,
    CONSTRAINT fk_from_user_id FOREIGN KEY (from_user_id) REFERENCES "user" (id) ON DELETE CASCADE,
    CONSTRAINT fk_chat_id FOREIGN KEY (chat_id) REFERENCES chat (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chat_user;
DROP TABLE IF EXISTS message;
DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS chat;

-- +goose StatementEnd
