// Code generated by go generate; DO NOT EDIT.

package storage

var SqlMap = map[string]string{
	"schema_version_1": `create table schema_version
(
    version text not null
);

create table users
(
    id       int primary key generated always as identity,
    name     text unique CHECK (name <> '' and length(name) <= 15),
    hash     text    not null CHECK (hash <> ''),
    about    TEXT    not null DEFAULT '',
    picture  text    not null default '',
    is_admin boolean not null default false
);

create table keys
(
    id         int primary key generated always as identity,
    key        text unique check (key <> '' and length(key) <= 20),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    user_id    integer unique references users (id)
);

create table settings
(
    name text not null,
    css  text not null default ''
);

create table forums
(
    id         int primary key generated always as identity,
    name       text unique                                        not null check ( name <> '' and length(name) < 120 ),
    position   int                                                not null,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

create table boards
(
    id           int primary key generated always as identity,
    name         text unique not null check ( name <> '' and length(name) < 120 ),
    position     int         not null,
    description  text,
    is_locked    bool        not null     default false,
    topics_count integer     not null     default 0,
    posts_count  integer     not null     default 0,
    created_at   timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at   timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    forum_id     integer     not null references forums (id)
);

create table topics
(
    id          int primary key generated always as identity,
    posts_count integer not null         default 0,
    is_sticky   boolean not null         default false,
    is_locked   boolean not null         default false,
    updated_at  timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    board_id    integer not null references boards (id),
    post_id     integer not null
);

create table posts
(
    id         int primary key generated always as identity,
    subject    text                                               not null check ( length(subject) <= 120 ),
    content    text                                               not null check ( length(content) <= 50000 ),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    topic_id   integer                                            null references topics (id),
    user_id    integer                                            not null references users (id)
);

alter table topics
    add foreign key (post_id) references posts (id) on delete cascade deferrable initially deferred;

create view posts_full as
select p.topic_id as topic_id,
       p.id       as post_id,
       p.subject,
       p.content,
       p.created_at,
       p.updated_at,
       u.id       as user_id,
       u.name,
       u.picture,
       u.about
from posts p
         left join users u on p.user_id = u.id;

create view topics_summary as
select t.Id as topic_id,
       p.subject,
       p.content,
       t.posts_count,
       t.updated_at,
       u.id as user_id,
       u.name,
       t.board_id,
       t.is_sticky,
       t.is_locked,
       t.post_id
from topics t
         left join posts p on t.post_id = p.id
         left join users u on p.user_id = u.id
order by t.is_sticky desc, t.updated_at desc;

create view forums_summary as
select b.id   as board_id,
       f.id   as forum_id,
       f.name as forum_name,
       b.name as board_name,
       b.description,
       b.topics_count,
       b.posts_count,
       b.updated_at
from boards b
         left join forums f on f.id = forum_id
order by f.position, b.position, f.id;

-- Triggers
--
-- Enforce topic and board lock rules
--
CREATE OR REPLACE FUNCTION check_is_locked() RETURNS TRIGGER AS
$$
DECLARE
    _is_locked bool;
    _board_id  int;
    _posts     int;
    _is_admin  bool;
BEGIN
    SELECT is_admin into _is_admin from users where id = NEW.user_id;
    IF (_is_admin) THEN
        RETURN NEW;
    end if;
    SELECT board_id into _board_id from topics where id = NEW.topic_id;
    SELECT is_locked into _is_locked from boards where id = _board_id;
    IF (_is_locked) THEN
        RETURN NULL;
    end if;
    SELECT count(id) into _posts from posts where topic_id = NEW.topic_id LIMIT 1;
    IF (_posts = 0) THEN
        RETURN NEW;
    end if;
    SELECT is_locked into _is_locked from topics where id = NEW.topic_id;
    IF (_is_locked) THEN
        RETURN NULL;
    end if;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER check_is_locked_topic
    BEFORE INSERT
    on posts
    FOR EACH ROW
EXECUTE PROCEDURE check_is_locked();
--
-- Increase topic and post count on boards
--
CREATE OR REPLACE FUNCTION count_on_board() RETURNS TRIGGER AS
$$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        UPDATE boards
        SET topics_count = topics_count + 1
        WHERE id = new.board_id;
        RETURN NEW;
    ELSIF (TG_OP = 'UPDATE') THEN
        UPDATE boards
        SET topics_count = boards.topics_count - 1,
            posts_count  = posts_count - (OLD.posts_count)
        WHERE id = OLD.board_id;
        UPDATE boards
        SET topics_count = boards.topics_count + 1,
            posts_count  = posts_count + (NEW.posts_count)
        WHERE id = NEW.board_id;
        RETURN NEW;
    ELSIF (TG_OP = 'DELETE') THEN
        UPDATE boards
        SET topics_count = topics_count - 1,
            posts_count  = posts_count - 1
        WHERE id = old.board_id;
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER increase_topic_count_on_board
    AFTER INSERT
    ON topics
    FOR EACH ROW
EXECUTE PROCEDURE count_on_board();
CREATE TRIGGER decrease_count_on_board
    AFTER DELETE
    ON topics
    FOR EACH ROW
EXECUTE PROCEDURE count_on_board();
CREATE TRIGGER update_count_on_board
    AFTER UPDATE of posts_count, board_id
    ON topics
    FOR EACH ROW
EXECUTE PROCEDURE count_on_board();
--
-- Count posts on topics
--
CREATE OR REPLACE FUNCTION count_posts_on_topic() RETURNS TRIGGER AS
$$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        UPDATE topics
        SET posts_count = topics.posts_count + 1
        WHERE id = NEW.topic_id;
        RETURN NEW;
    ELSIF (TG_OP = 'DELETE') THEN
        UPDATE topics
        SET posts_count = topics.posts_count - 1
        WHERE id = OLD.topic_id;
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER increase_post_count_on_topics
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE count_posts_on_topic();
CREATE TRIGGER decrease_post_count_on_topics
    AFTER DELETE
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE count_posts_on_topic();
--
-- Get last updated_at time of a topic
--
CREATE OR REPLACE FUNCTION get_topic_updated_at() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE topics
    SET updated_at = (SELECT updated_at from posts where topic_id = old.id order by updated_at desc limit 1)
    WHERE id = old.id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER get_topic_updated_at
    AFTER UPDATE of posts_count
    on topics
    FOR EACH ROW
EXECUTE PROCEDURE get_topic_updated_at();
--
-- Get last updated_at time of a board
--
CREATE OR REPLACE FUNCTION get_board_updated_at() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE boards
    SET updated_at = coalesce((SELECT updated_at from topics where board_id = old.id order by updated_at desc limit 1),
                              boards.created_at)
    WHERE id = old.id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER get_board_updated_at
    AFTER UPDATE of posts_count
    on boards
    FOR EACH ROW
EXECUTE PROCEDURE get_board_updated_at();

insert into settings (name)
values ('vpub');
insert into keys (key)
values ('admin');`,
	"schema_version_10": `ALTER TABLE boards
    ADD COLUMN textsearchable_index_col tsvector
        GENERATED ALWAYS AS ( to_tsvector('english', coalesce(name, '') || ' ' || coalesce(description, ''))) STORED;

CREATE INDEX textsearch_idx_boards ON boards USING GIN (textsearchable_index_col);

ALTER TABLE posts
    ADD COLUMN textsearchable_index_col tsvector
        GENERATED ALWAYS AS ( to_tsvector('english', coalesce(subject, '') || ' ' || coalesce(content, ''))) STORED;

CREATE INDEX textsearch_idx_posts ON posts USING GIN (textsearchable_index_col);

ALTER TABLE users
    ADD COLUMN textsearchable_index_col tsvector
        GENERATED ALWAYS AS ( to_tsvector('english', coalesce(name, '') || ' ' || coalesce(about, ''))) STORED;

CREATE INDEX textsearch_idx_users ON users USING GIN (textsearchable_index_col);

-- We are storing ID as a string for the sake of convenience,
-- since we are only using it for search results.

CREATE OR REPLACE VIEW search_items AS
    SELECT
        text 'users' AS origin_table,
        id::text,
        name AS title,
        about AS content,
        textsearchable_index_col AS searchable_element
    FROM
        users
UNION ALL
    SELECT
        text 'posts' AS origin_table,
        concat(topic_id, '#', id)::text AS id,
        subject AS title,
        content,
        textsearchable_index_col AS searchable_element
    FROM
        posts
UNION ALL
    SELECT
        text 'boards' AS origin_table,
        id::text,
        name AS title,
        description AS content,
        textsearchable_index_col AS searchable_element
    FROM
        boards;

CREATE OR REPLACE FUNCTION search_with_highlights(search_term text)
    RETURNS TABLE (
                      origin_table text,
                      id text,
                      title text,
                      content text,
                      highlighted_title text,
                      highlighted_content text,
                      rank float4
                  ) AS $$
DECLARE
    query tsquery := to_tsquery('english', search_term);
BEGIN
    RETURN QUERY
        SELECT
            si.origin_table,
            si.id::text,
            si.title,
            si.content,
            ts_headline('english', si.title, query, 'StartSel=<mark>, StopSel=</mark>, MaxFragments=1') AS highlighted_title,
            ts_headline('english', si.content, query, 'StartSel=<mark>, StopSel=</mark>, MaxFragments=3, FragmentDelimiter=..., MaxWords=13, MinWords=3') AS highlighted_content,
            ts_rank(si.searchable_element, query) AS rank
        FROM
            search_items si
        WHERE
            query @@ si.searchable_element
        ORDER BY
            rank DESC;
END;
$$ LANGUAGE plpgsql;
`,
	"schema_version_2": `alter table settings
    add column footer text not null default ''`,
	"schema_version_3": `-- Admins are allowed to post even when locked.
CREATE OR REPLACE FUNCTION check_is_locked() RETURNS TRIGGER AS
$$
DECLARE
    _is_locked bool;
    _forum_id  int;
    _board_id  int;
    _posts     int;
    _is_admin  bool;
BEGIN
    SELECT is_admin into _is_admin from users where id = NEW.user_id;
    IF (_is_admin) THEN
        RETURN NEW;
    end if;
    SELECT board_id into _board_id from topics where id = NEW.topic_id;
    SELECT forum_id into _forum_id from boards where id = _board_id;
    SELECT is_locked into _is_locked from forums where id = _forum_id;
    IF (_is_locked) THEN
        RETURN NULL;
    end if;
    SELECT is_locked into _is_locked from boards where id = _board_id;
    IF (_is_locked) THEN
        RETURN NULL;
    end if;
    SELECT count(id) into _posts from posts where topic_id = NEW.topic_id LIMIT 1;
    IF (_posts = 0) THEN
        RETURN NEW;
    end if;
    SELECT is_locked into _is_locked from topics where id = NEW.topic_id;
    IF (_is_locked) THEN
        RETURN NULL;
    end if;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

alter table forums
    add column is_locked bool not null default false;`,
	"schema_version_4": `alter table settings
    add column per_page int not null default 50;`,
	"schema_version_5": `alter table settings
    add column url text not null default '';`,
	"schema_version_6": `CREATE OR REPLACE VIEW topics_summary as
select t.Id as topic_id,
       p.subject,
       p.content,
       t.posts_count,
       t.updated_at,
       u.id as user_id,
       u.name,
       t.board_id,
       t.is_sticky,
       t.is_locked,
       t.post_id,
       p.created_at
from topics t
         left join posts p on t.post_id = p.id
         left join users u on p.user_id = u.id
order by t.is_sticky desc, t.updated_at desc;`,
	"schema_version_7": `CREATE OR REPLACE PROCEDURE remove_board(_bid int) as
$$
BEGIN
    alter table posts
        disable trigger ALL;
    alter table topics
        disable trigger ALL;
    delete
    from posts
    where topic_id in (select id
                       from topics
                       where board_id = _bid);
    delete from topics where board_id = _bid;
    delete from boards where id = _bid;
    alter table posts
        enable trigger ALL;
    alter table topics
        enable trigger ALL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE remove_forum(_fid int) as
$$
DECLARE
    board int;
BEGIN
    FOR board in SELECT id FROM boards where forum_id = _fid
        LOOP
            call remove_board(board);
        END LOOP;
    delete from forums where id = _fid;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION remove_topic(_tid int) returns void as
$$
DECLARE
    pid int;
BEGIN
    -- Get the original post id
    select post_id into pid from topics where id = _tid;
    -- Delete all the posts of the topic except the original one
    delete from posts where topic_id = _tid and id != pid;
    -- Delete original post
    delete from posts where id = pid;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE PROCEDURE remove_user(_uid int) as
$$
BEGIN
    -- Delete all topics
    perform remove_topic(topic_id) from topics_summary where user_id = _uid;
    -- Delete all the posts
    delete from posts where user_id = _uid;
    -- Delete associated key
    delete from keys where user_id = _uid;
    -- Delete user
    delete from users where id = _uid;
END;
$$ LANGUAGE plpgsql;`,
	"schema_version_8": `ALTER TABLE settings
    ADD lang text default 'en'::text NOT NULL;`,
	"schema_version_9": `ALTER TABLE users
    ADD picture_alt TEXT NOT NULL DEFAULT '';`,
}
