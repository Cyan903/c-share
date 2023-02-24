create table users (
	id int auto_increment primary key,
	nickname varchar(32) charset utf8 not null,
	email varchar(254) not null unique,
	pw_bcrypt char(60) not null,
	created_at datetime not null
);

-- 0 = public
-- 1 = private
-- 2 = unlisted
create table files (
	id char(10) not null unique,
	user int not null,
	file_size int not null,
	file_type char(32) not null,
	file_pass char(60) default "",
	permissions int default 0 not null,
	created_at datetime not null
);
