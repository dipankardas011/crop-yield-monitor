-- Create the database
CREATE DATABASE IF NOT EXISTS auth;

-- Switch to the newly created database
USE auth;

-- Create the "users" table
CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(20) PRIMARY KEY,
    name VARCHAR(50),
    email VARCHAR(30)
);

-- Create the "password" table
CREATE TABLE IF NOT EXISTS passwords (
    username VARCHAR(20) PRIMARY KEY,
    password VARCHAR(300),
    salt VARCHAR(10)
);

create user 'auth'@'%' identified by '12345';
grant all on *.* to 'auth';
