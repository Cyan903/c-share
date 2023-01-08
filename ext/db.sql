create table users (
	id int auto_increment primary key,
	nickname varchar(32) charset utf8 not null,
	email varchar(254) not null unique,
	pw_bcrypt char(60) not null,
	created_at datetime not null
);