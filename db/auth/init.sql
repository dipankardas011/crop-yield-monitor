create database auth;

create table users(
	ID auto-increment primary;
	Name string;
	password passwords;
);

create table passwords(
	ID primary;// refered by the user table
	hashedpass string;
	salt string;
)
