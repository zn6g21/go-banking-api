CREATE DATABASE IF NOT EXISTS api_database;

USE api_database;

CREATE TABLE customers (
    cif_no INT PRIMARY KEY AUTO_INCREMENT,
    name_kana VARCHAR(255) NOT NULL,
    name_kanji VARCHAR(255) NOT NULL,
    birth_date DATE NOT NULL,
    prefecture VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    town VARCHAR(255) NOT NULL,
    street VARCHAR(255) NOT NULL,
    building VARCHAR(255),
    room VARCHAR(255),
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE accounts (
    id INT PRIMARY KEY AUTO_INCREMENT,
    cif_no INT NOT NULL,
    status VARCHAR(20) NOT NULL,
    branch_code VARCHAR(3) NOT NULL,
    account_number VARCHAR(7) NOT NULL,
    account_type VARCHAR(2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_accounts_customers FOREIGN KEY (cif_no) REFERENCES customers(cif_no)
);

CREATE TABLE tokens (
    access_token VARCHAR(255) PRIMARY KEY,
    refresh_token VARCHAR(255) NOT NULL,
    scopes TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    cif_no INT NOT NULL,
    UNIQUE KEY uk_tokens_refresh_token (refresh_token),
    CONSTRAINT fk_tokens_customers FOREIGN KEY (cif_no) REFERENCES customers(cif_no)
);

