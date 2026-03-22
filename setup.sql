-- Создание базы данных
CREATE DATABASE cyberdefense_db;

-- Выбор базы данных для дальнейшей работы
USE cyberdefense_db;

-- Создание таблицы users
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    phone VARCHAR(15)
);