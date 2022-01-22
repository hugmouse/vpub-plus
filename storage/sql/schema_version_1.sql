create table schema_version (
    version text not null
);

create table keys (
    id integer primary key autoincrement,
    key text unique check (key <> '' and length(key) <= 20),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
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
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    forum_id not null references forums(id)
);

create table forums (
    id integer primary key autoincrement,
    name text unique not null check ( name <> '' and length(name) < 120 ),
    position int not null,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

create table posts (
    id integer primary key autoincrement,
    subject text not null check ( length(subject) <= 120 ),
    content text not null check ( length(content) <= 50000 ),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    topic_id integer not null references posts(id),
    user_id integer not null references users(id)
);

create table topics (
    id integer primary key autoincrement,
    posts_count integer not null default 0,
    is_sticky boolean not null default false,
    is_locked boolean not null default false,
    updated_at NOTtimestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
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
        u.picture
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
    AFTER INSERT ON posts
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
    AFTER DELETE ON posts
BEGIN
    UPDATE topics
    SET posts_count = topics.posts_count-1
    WHERE id=old.topic_id;
END;

-- CREATE TRIGGER count_post_on_board
--     BEFORE UPDATE of posts_count, board_id ON topics
-- BEGIN
--     UPDATE boards
--         SET posts_count = posts_count-(old.posts_count + 1)
--         WHERE id=old.board_id;
--     UPDATE boards
--         SET posts_count = posts_count+(new.posts_count + 1)
--         WHERE id=new.board_id;
-- END;
--
-- CREATE TRIGGER get_topic_updated_at
--     AFTER UPDATE of posts_count on topics
-- BEGIN
--     UPDATE topics
--     SET updated_at = (SELECT updated_at from posts where topic_id=old.id order by updated_at desc limit 1)
--     WHERE id=old.id;
-- end;

insert into settings (name) values ('vpub');
insert into keys (key) values ('admin');