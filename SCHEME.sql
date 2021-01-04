
drop table alarm_target;
drop table todo_alarm;
drop table todos;

create table todos (
    id bigint(20) not null auto_increment comment 'id',
    user_id bigint(20) null comment '사용자 아이디',
    title varchar(500) not null comment '제목',
    priority varchar(20) not null comment '우선순위: T, A, B, C, D',
    status varchar(20) not null comment '상태: REGISTED, DOING, DONE, HOLD, ARCHIVE',
    completion_level int null comment '완성도: 0 ~ 100 %',
    created_at timestamp not null comment '생성일',
    modified_at timestamp null comment '수정일',
    done_at timestamp null comment '완료일',
    primary key (id)
)
comment 'TODO 목록';

create table todo_alarm (
    id bigint(20) not null auto_increment comment 'id',
    todo_id bigint(20) not null comment 'todo id',
    period_type varchar(20) not null comment 'alarm type: ONCE, INTERVAL',
    alarm_type varchar(20) not null comment 'PUSH, MAIL',
    alarm_date date null comment '알람일자, 특정일',
    alarm_time time not null comment '알람시간',
    active_yn char(1) not null default 'Y' comment '활성화여부',
    created_at timestamp not null comment '생성일',
    modified_at timestamp null comment '수정일',
    primary key (id)
)
comment 'ALARM';

alter table todo_alarm add constraint fk_with_todos foreign key(todo_id) references todos (id) ON DELETE CASCADE;

create table alarm_target (
    id bigint(20) not null auto_increment comment 'id',
    todo_alarm_id bigint(20) not null comment '알람 아이디',
    phone varchar(100) null comment '전화번호',
    email varchar(100) null comment '메일주소',
    user_id bigint(20) null comment '사용자 아이디',
    active_yn char(1) not null default 'Y' comment '활성화여부',
    created_at timestamp not null comment '생성일',
    modified_at timestamp null comment '수정일',
    primary key (id)
)
comment '알람대상정보';

alter table alarm_target add constraint fk_with_alarm_target foreign key(todo_alarm_id) references todo_alarm(id) ON DELETE CASCADE;