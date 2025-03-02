-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER  NOT NULL,
    task_name   TEXT     NOT NULL,
    description TEXT,
    status      TEXT DEFAULT 'Pending',
    created_at  datetime NOT NULL,
    updated_at  datetime NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_tasks_userid ON tasks (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if exists tasks;
-- +goose StatementEnd
