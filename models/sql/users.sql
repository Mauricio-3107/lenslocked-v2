CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

SELECT
    users.id,
    users.email,
    sessions.id as session_id
FROM
    users
    JOIN sessions ON users.id = sessions.user_id;

SELECT
    users.id,
    users.email,
    users.password_hash
FROM
    sessions
    JOIN users ON users.id = sessions.user_id
WHERE
    sessions.token_hash = $ 1;