create table user (
                      id varchar(128) not null primary key,
                      name varchar(255) not null,
                      email varchar(255) not null,
                      bio char(255),
                      image char(255),
                      create_at timestamp not null default current_timestamp
);

create table tweet (
                       id char(26) not null primary key,
                       body text not null,
                       posted_by varchar(128) not null,
                       posted_at timestamp not null default current_timestamp,
                       reply_to char(26) default null,
                       like_count int not null default 0
);

create table likes (
                       id char(26) not null primary key,
                       tweet_id char(26) not null,
                       user_id char(128) not null
);

create table tag (
                     id char(26) not null primary key,
                     tag text not null
);

create table follow(
                       id char(26) not null primary key,
                       follower_id char(128) not null,
                       followee_id char(128) not null,
                       created_at timestamp not null default current_timestamp
);

create table tweet_tag (
                           id char(26) not null primary key,
                           tweet_id char(26) not null,
                           tag_id char(26) not null
);