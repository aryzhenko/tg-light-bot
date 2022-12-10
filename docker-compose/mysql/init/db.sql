create database if not exists tg_light_bot;
create table tg_light_bot.users
(
    id                          bigint unsigned not null,
    first_name                  varchar(255),
    last_name                   varchar(255),
    user_name                   varchar(255),
    is_bot                      bool,
    language_code               varchar(2),
    can_join_groups             bool,
    can_read_all_group_messages bool,
    supports_inline_queries     bool,
    created_at                  datetime,
    last_activity               datetime,
    constraint users_pk
        primary key (id)
);

create table tg_light_bot.light_status
(
    id int unsigned not null auto_increment,
    is_on bool,
    changed_at datetime,
        constraint light_status_pk
        primary key (id)
);
