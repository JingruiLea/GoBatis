create table user
(
    id         bigint unsigned auto_increment comment '主键id'
        primary key,
    user_id    bigint unsigned         not null comment '用户id',
    username   varchar(150) default '' not null comment '用户名',
    phone      varchar(20)             not null comment '电话',
    age        int          default 0  not null comment '年龄',
    created_at datetime(3)             not null comment '创建时间',
    updated_at datetime(3)             not null comment '修改时间',
    is_deleted tinyint unsigned        not null comment '是否已删除'
);