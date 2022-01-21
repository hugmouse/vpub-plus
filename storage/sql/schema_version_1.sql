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
    key_id integer not null unique references keys(id)
);

create table settings (
    name text not null,
    css text not null default ''
);

create table boards (
    id integer primary key autoincrement,
    name text unique not null check ( name <> '' and length(name) < 120 ),
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
    name text unique not null check ( name <> '' and length(name) < 120 ),
    position int not null
);


create table posts (
    id integer primary key autoincrement,
    subject text not null check ( length(subject) <= 120 ),
    content text not null check ( length(content) <= 50000 ),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    topic_id integer not null references posts(id),
    user_id integer not null references users(id)
);

create table topics (
    id integer primary key autoincrement,
    posts_count integer not null default 0,
    is_sticky boolean not null default false,
    is_locked boolean not null default false,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    board_id integer not null references boards(id),
    post_id integer not null references posts(id) on delete cascade deferrable initially deferred
);

create view postUsers as
    select
        p.topic_id as topic_id,
        p.id as post_id,
        p.subject,
        p.content,
        p.created_at,
        p.updated_at,
        u.id as user_id,
        u.name,
       u.picture,
        p.topic_id
    from posts p
    left join users u on p.user_id = u.id;

create view boardTopics as
    select
        t.Id as topic_id,
        p.subject,
        p.content,
        t.posts_count,
        t.updated_at,
        u.id as user_id,
        u.name,
        t.board_id
    from topics t
        left join posts p on t.post_id = p.id
        left join users u on p.user_id = u.id
    order by t.updated_at desc;

create view forumBoards as
    select
           b.id as board_id,
           f.name as forum_name,
           b.name as board_name,
           b.description,
           b.topics_count,
           b.posts_count,
           b.created_at
    from boards b
    left join forums f on f.id = forum_id
    order by f.position, b.position, f.id;

CREATE TRIGGER increase_topic_count_on_board
    BEFORE INSERT ON topics
BEGIN
    UPDATE boards
    SET topics_count = topics_count+1
    WHERE id=new.board_id;
END;

CREATE TRIGGER increase_post_count_on_topics
    BEFORE INSERT ON posts
BEGIN
    UPDATE topics
    SET posts_count = topics.posts_count+1
    WHERE id=new.topic_id;
END;

CREATE TRIGGER decrease_topic_count_on_board
    BEFORE DELETE ON topics
BEGIN
    UPDATE boards
    SET topics_count = topics_count-1
    WHERE id=old.board_id;
END;

CREATE TRIGGER decrease_post_count_on_topics
    BEFORE DELETE ON posts
BEGIN
    UPDATE topics
    SET posts_count = topics.posts_count-1
    WHERE id=old.topic_id;
END;

CREATE TRIGGER count_post_on_board
    BEFORE UPDATE of posts_count, board_id ON topics
BEGIN
    UPDATE boards
        SET posts_count = posts_count-(old.posts_count + 1)
        WHERE id=old.board_id;
    UPDATE boards
        SET posts_count = posts_count+(new.posts_count + 1)
        WHERE id=new.board_id;
END;

insert into settings (name) values ('vpub');
insert into keys (key) values ('admin');

--
-- create table schema_version (
--     version text not null
-- );
--
-- create table keys (
--                       id integer primary key autoincrement,
--                       key text unique check (key <> '' and length(key) <= 20),
--                       created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
--                       user_id integer unique,
--                       foreign key (user_id) references users(id)
-- );
--
-- create table users (
--                        id integer primary key autoincrement,
--                        name text unique CHECK (name <> '' and length(name) <= 15),
--                        hash text not null CHECK (hash <> ''),
--                        about TEXT not null DEFAULT '',
--                        picture text not null default '',
--                        is_admin boolean not null default false,
--                        key_id integer not null unique,
--                        foreign key (key_id) references keys(id)
-- );
--
-- create table settings (
--                           name text not null,
--                           css text not null default ''
-- );
--
-- create table boards (
--                         id integer primary key autoincrement,
--                         name text unique not null check ( name <> '' and length(name) < 120 ),
--                         position int not null,
--                         description text,
--                         is_locked bool not null default false,
--                         topics_count integer not null default 0,
--                         posts_count integer not null default 0,
--                         created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
--                         forum_id not null references forums(id)
-- );
--
-- create table forums (
--                         id integer primary key autoincrement,
--                         name text unique not null check ( name <> '' and length(name) < 120 ),
--                         position int not null
-- );
--
--
-- create table posts (
--                        id integer primary key autoincrement,
--                        subject text not null check ( length(subject) <= 120 ),
--                        content text not null check ( length(content) <= 50000 ),
--                        created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
--                        updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
--                        topic_id integer not null references topics(id),
--                        user_id integer not null references users(id)
-- );
--
-- create table topics (
--                         id integer primary key autoincrement,
--                         reply_count integer not null default 0,
--                         is_sticky boolean not null default false,
--                         is_locked boolean not null default false,
--                         updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
--                         board_id integer not null references boards(id)
-- );
--
-- create table topicPosts (
--                             topic_id integer not null references topics(id),
--                             post_id integer not null references posts(id)
-- );
--
-- create view postUsers as
-- select
--     p.id as post_id,
--     p.subject,
--     p.content,
--     p.created_at,
--     u.id as user_id,
--     u.name
-- from posts p
--          left join users u on p.user_id = u.id;
--
-- create view boardTopics as
-- select
--     t.Id,
--     p.subject,
--     p.content,
--     p.created_at,
--     p.user_id,
--     p.name
-- from topicPosts tP
--          left join postUsers p on p.post_id = tP.post_id
--          left join topics t on t.id = tP.topic_id;
--
-- --
-- -- CREATE TRIGGER check_is_locked_before_insert
-- --     BEFORE INSERT ON posts
-- -- BEGIN
-- --     select
-- --         case
-- --             when is_locked is true then
-- --                 raise (abort, 'Topic is locked')
-- --             end
-- --     from topics
-- --     where id=new.topic_id;
-- --     select
-- --         case
-- --             when is_locked is true then
-- --                 raise (abort, 'Board is locked')
-- --             end
-- --     from topics t
-- --     left join boards b on b.id = t.board_id
-- --     where b.id=t.board_id;
-- -- END;
-- --
--
-- CREATE TRIGGER count_topic_before_insert
--     BEFORE INSERT ON topics
-- BEGIN
--     UPDATE boards
--     SET topics_count = topics_count+1
--     WHERE id=new.board_id;
-- END;
--
-- -- CREATE TRIGGER count_post_before_insert
-- --     BEFORE INSERT ON posts
-- -- BEGIN
-- --     UPDATE
-- --         boards
-- --     SET topics_count = case when new.topic_id is null then topics_count+1 else topics_count end,
-- --         posts_count = posts_count+1
-- --     WHERE id=new.board_id;
-- --     UPDATE
-- --         posts
-- --     SET reply_count = reply_count+1
-- --     WHERE id=new.topic_id;
-- -- END;
-- --
-- -- CREATE TRIGGER count_post_before_delete
-- --     BEFORE DELETE ON posts
-- -- BEGIN
-- -- UPDATE
-- --     boards
-- -- SET topics_count = case when old.topic_id is null then topics_count-1 else topics_count end,
-- --     posts_count = posts_count-1
-- -- WHERE id=old.board_id;
-- -- UPDATE
-- --     posts
-- -- SET reply_count = reply_count-1
-- -- WHERE id=old.topic_id;
-- -- END;
-- --
-- -- CREATE TRIGGER count_post_before_update
-- --     BEFORE UPDATE of board_id on posts
-- --     WHEN new.topic_id is null and old.board_id <> new.board_id
-- -- BEGIN
-- -- UPDATE
-- --     posts
-- -- SET board_id = new.board_id
-- -- WHERE topic_id=new.id;
-- -- UPDATE
-- --     boards
-- -- SET topics_count=topics_count-1,
-- --     posts_count=posts_count-(old.reply_count+1)
-- -- WHERE id=old.board_id;
-- -- UPDATE
-- --     boards
-- -- SET topics_count=topics_count+1,
-- --     posts_count=posts_count+(new.reply_count+1)
-- -- WHERE id=new.board_id;
-- -- end;
--
-- insert into settings (name) values ('vpub');
-- insert into keys (key) values ('admin');