CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    nickname      VARCHAR(32)  NOT NULL,
    email         VARCHAR(64)  NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    dob           VARCHAR(32)  NOT NULL
);
CREATE TABLE subscriptions
(
    id               SERIAL PRIMARY KEY,
    user_id          INT    NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    subscribed_to_id INT    NOT NULL,
    FOREIGN KEY (subscribed_to_id) REFERENCES users(id) ON DELETE CASCADE
);
