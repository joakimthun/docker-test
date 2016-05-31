 create table users (
 	id serial unique primary key, 
 	name varchar(512) not null,
 	email varchar(256) unique not null
)