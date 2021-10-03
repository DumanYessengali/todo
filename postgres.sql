create table todo (
                      id serial primary key,
                      title varchar(255),
                      description varchar(255),
                      created timestamp not null,
                      updated timestamp not null
);

create table users (
                       id serial primary key,
                       email varchar(255) not null,
                       password varchar(255) not null,
                       name varchar(255) not null,
                       todo_id int
);

ALTER TABLE users DROP COLUMN todo_id;

alter table todo add column user_id int;

alter table users add column active BOOLEAN NOT NULL DEFAULT TRUE;