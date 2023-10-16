/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
CREATE TABLE test (
	id serial PRIMARY KEY,
	name VARCHAR ( 50 ) UNIQUE NOT NULL,
);

INSERT INTO test (name) VALUES ('test1');
INSERT INTO test (name) VALUES ('test2');

-- public.users definition

-- Drop table

-- DROP TABLE users;

CREATE TABLE users (
	id serial4 NOT NULL,
	phone_number varchar NOT NULL,
	full_name varchar(60) NOT NULL,
	"password" varchar(72) NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	last_login timestamp NULL,
	login_count int4 NULL,
	CONSTRAINT users_phone_number_key UNIQUE (phone_number),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);