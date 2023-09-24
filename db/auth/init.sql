-- Create the database
CREATE DATABASE IF NOT EXISTS auth;

-- Switch to the newly created database
USE auth;

-- Create the "users" table
CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(255)
);

-- Create the "password" table
CREATE TABLE IF NOT EXISTS password (
    username VARCHAR(50) PRIMARY KEY,
    password VARCHAR(255)
);

create user 'auth'@'%' identified by '12345';
grant all on *.* to 'auth';
