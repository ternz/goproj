drop database if exists db_authority;
create database db_authority character set utf8;
use db_authority;

drop table if exists tb_user;
create table tb_user(
	`id`		bigint unsigned		auto_increment,
	`code`	varchar(256) 		not null unique,
	`name`	varchar(64) 		not null,
	primary key (id)
);

drop table if exists tb_authority;
create table tb_authority(
	`id`		bigint unsigned		auto_increment,
	`code`	varchar(256) 		not null unique,
	`name`	varchar(64) 		not null,
	`group_id`	bigint unsigned	 default 0,
	primary key (id)
);

drop table if exists tb_user_authority;
create table tb_user_authority(
	`id`	bigint unsigned		auto_increment,
	`user_id`	bigint unsigned,
	`authority_id`	bigint unsigned,
	primary key (id)
);