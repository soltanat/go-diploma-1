CREATE SCHEMA IF NOT EXISTS service_diploma_1;

CREATE TABLE service_diploma_1.users
(
    login      VARCHAR(255) PRIMARY KEY,
    password   VARCHAR(255) NOT NULL,
    whole      BIGINT       NOT NULL,
    decimal    BIGINT       NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE service_diploma_1.orders
(
    number      BIGINT PRIMARY KEY,
    status      VARCHAR(255) NOT NULL,
    whole       BIGINT       NOT NULL,
    decimal     BIGINT       NOT NULL,
    uploaded_at TIMESTAMP WITH TIME ZONE,
    user_id     VARCHAR(255) REFERENCES service_diploma_1.users (login)
);

CREATE TABLE service_diploma_1.withdrawals
(
    number       BIGINT PRIMARY KEY,
    whole        BIGINT NOT NULL,
    decimal      BIGINT NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE,
    user_id      VARCHAR(255) REFERENCES service_diploma_1.users (login)
);
