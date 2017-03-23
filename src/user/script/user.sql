/*
tb_user
id			bigint			auto, pk
name		varchar(64)
nickname	varchar(32)
phone		varchar(32)
email		varchar(256)
origin		varchar(256)

tb_role
id			bigint			auto, pk
name		varchar(256)

tb_user_role
user_id		bigint			fk
role_id		bigint			fk
pk(user_id, role_id)

tb_privilege
id			bigint			auto, pk
name		varchar(256)	unique

tb_role_privilege
role_id		bigint			fk
privilege_id	bigint		fk
pk(role_id, privilege_id)
*/

create database db_user character set utf8;
use db_user;

create table tb_user (
	id		bigint unsigned		auto_increment,
	name	varchar(64),
	nickname	varchar(32),
	phone		varchar(32),
	email		varchar(256),
	origin		varchar(256),
	primary key (id)
)

create table tb_role (
	id		bigint unsigned		auto_increment,
	name	varchar(256)	not null,
	primary key(id)
)

create table tb_role_user (
	user_id	bigint unsigned,
	role_id	bigint unsigned,
	primary key(user_id, role_id)
)

create table tb_privilege (
	id		bigint unsigned auto_increment,
	name	varchar(256) not null,
	primary key(id)
)

create table tb_role_privilege (
	role_id		bigint unsigned,
	privilege_id bigint unsigned,
	primary key(role_id, privilege_id)
)