
-- schema.sql

drop database if exists gold;

create database gold;

use gold;

grant select, insert, update, delete on gold.* to 'root'@'127.0.0.1:3306' identified by 'root';

CREATE TABLE GoldEveryDay
(
`getDate` date NOT NULL,
`price` decimal(10,2),
`open` decimal(10,2),
`high` decimal(10,2),
`low` decimal(10,2),
PRIMARY KEY (`getDate`)
) engine=innodb default charset=utf8;

CREATE TABLE GoldEveryMin
(
`GetTime` datetime NOT NULL,
`price` decimal(10,2),
`buy` decimal(10,2),
`sell` decimal(10,2),
`highMid` decimal(10,2),
`lowMid` decimal(10,2),
PRIMARY KEY (`getTime`)
) engine=innodb default charset=utf8;

CREATE TABLE actionLog
(
`ActionTime` datetime NOT NULL,
`action` int,
`price` decimal(10,2),
`amount` decimal(10,2),
`actionType` int,
`earn` decimal(10,2),
`watermark` decimal(10,2),
PRIMARY KEY (`ActionTime`)
) engine=innodb default charset=utf8;