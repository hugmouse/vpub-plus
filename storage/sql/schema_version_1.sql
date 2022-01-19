-- create schema version table
create table schema_version (
    version text not null
);

create table keys (
    id integer primary key autoincrement,
    key text unique check (key <> '' and length(key) <= 20),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    user_id integer unique,
    foreign key (user_id) references users(id)
);

create table users (
    id integer primary key autoincrement,
    name text unique CHECK (name <> '' and length(name) <= 15),
    hash text not null CHECK (hash <> ''),
    about TEXT not null DEFAULT '',
    picture text not null default '',
    is_admin boolean not null default false,
    key_id integer not null unique,
    foreign key (key_id) references keys(id)
);

create table settings (
    name text not null,
    css text not null default ''
);

create table boards (
    id integer primary key autoincrement,
    name text not null check ( name <> '' and length(name) < 120 ),
    position int not null,
    description text,
    is_locked bool not null default false,
    topics_count integer not null default 0,
    posts_count integer not null default 0,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    forum_id not null references forums(id)
);

create table forums (
    id integer primary key autoincrement,
    name text not null check ( name <> '' and length(name) < 120 ),
    position int not null
);

create table posts (
    id integer primary key autoincrement,
    subject text not null check ( length(subject) <= 120 ),
    content text not null check ( length(content) <= 50000 ),
    reply_count integer not null default 0,
    is_sticky boolean not null default false,
    is_locked boolean not null default false,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    topic_id integer references posts(id),
    board_id integer references boards(id),
    user_id integer references users(id)
);

create view topics as
    select
        t.id,
        t.user_id,
        u.name,
        t.board_id,
        t.subject,
        count(r.topic_id) replies,
        t.created_at,
        max(ifnull(r.updated_at, t.updated_at)) updated_at,
        t.is_sticky,
        t.is_locked
    from posts t
             left join posts r on r.topic_id = t.id
             left join users u on t.user_id = u.id
    where t.topic_id is null
    group by t.id order by t.is_sticky desc, updated_at desc;

create view postUsers as
    select
        p.id as post_id,
        p.subject,
        p.content,
        p.created_at,
        p.topic_id,
        p.board_id,
        p.is_sticky,
        p.is_locked,
        u.id as user_id,
        u.name,
        u.picture
    from posts p
         left join users u on p.user_id = u.id;

CREATE TRIGGER check_is_locked_before_insert
    BEFORE INSERT ON posts
BEGIN
    select
        case
            when is_locked is true then
                raise (abort, 'Topic is locked')
            end
    from posts
    where id=new.topic_id;
    select
        case
            when is_locked is true then
                raise (abort, 'Board is locked')
            end
    from boards
    where id=new.board_id;
END;

CREATE TRIGGER count_post_before_insert
    BEFORE INSERT ON posts
BEGIN
    UPDATE
        boards
    SET topics_count = case when new.topic_id is null then topics_count+1 else topics_count end,
        posts_count = posts_count+1
    WHERE id=new.board_id;
    UPDATE
        posts
    SET reply_count = reply_count+1
    WHERE id=new.topic_id;
END;

CREATE TRIGGER count_post_before_delete
    BEFORE DELETE ON posts
BEGIN
    UPDATE
        boards
    SET topics_count = case when old.topic_id is null then topics_count-1 else topics_count end,
        posts_count = posts_count-1
    WHERE id=old.board_id;
    UPDATE
        posts
    SET reply_count = reply_count-1
    WHERE id=old.topic_id;
END;

CREATE TRIGGER count_post_before_update
    BEFORE UPDATE of board_id on posts
    WHEN new.topic_id is null and old.board_id <> new.board_id
BEGIN
    UPDATE
        posts
    SET board_id = new.board_id
    WHERE topic_id=new.id;
    UPDATE
        boards
    SET topics_count=topics_count-1,
        posts_count=posts_count-(old.reply_count+1)
    WHERE id=old.board_id;
    UPDATE
        boards
    SET topics_count=topics_count+1,
        posts_count=posts_count+(new.reply_count+1)
    WHERE id=new.board_id;
end;

insert into settings (name) values ('vpub');
insert into keys (key) values ('admin');