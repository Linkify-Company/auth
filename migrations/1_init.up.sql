CREATE TABLE IF NOT EXISTS role (
    id INTEGER NOT NULL UNIQUE
);
CREATE INDEX IF NOT EXISTS id_idx ON role (id);
INSERT INTO role (id) VALUES (0);
INSERT INTO role (id) VALUES (1);
INSERT INTO role (id) VALUES (2);


CREATE TABLE IF NOT EXISTS "user" (
    id SERIAL PRIMARY KEY ,
    email TEXT NOT NULL UNIQUE,
    pass_hash TEXT NOT NULL,
    role INTEGER NOT NULL REFERENCES role (id),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP
);
CREATE INDEX IF NOT EXISTS id_idx ON "user" (id);
CREATE INDEX IF NOT EXISTS email_idx ON "user" (email);