create table users (
                      id varchar(128) NOT NULL primary key,
                      email varchar(255) NOT NULL
);

create table tweets (
                        id char(26) NOT NULL primary key,
                        posted_by varchar(128) NOT NULL,
                        content varchar(140) NOT NULL
);

insert into tweets values ('00000000000000000000000002', '00000000000000000000000001', 'hello');