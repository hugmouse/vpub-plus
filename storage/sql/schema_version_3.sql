-- Admins are allowed to post even when locked.
CREATE OR REPLACE FUNCTION check_is_locked() RETURNS TRIGGER AS $$
DECLARE
    _is_locked  bool;
    _forum_id   int;
    _board_id   int;
    _posts      int;
    _is_admin   bool;
BEGIN
    SELECT is_admin  into _is_admin from users where id=NEW.user_id;
    IF (_is_admin) THEN
        RETURN NEW;
    end if;
    SELECT board_id into _board_id from topics where id=NEW.topic_id;
    SELECT forum_id into _forum_id from boards where id=_board_id;
    SELECT is_locked into _is_locked from forums where id=_forum_id;
    IF (_is_locked) THEN
        RETURN NULL;
    end if;
    SELECT is_locked into _is_locked from boards where id=_board_id;
    IF (_is_locked) THEN
        RETURN NULL;
    end if;
    SELECT count(id) into _posts from posts where topic_id=NEW.topic_id LIMIT 1;
    IF (_posts = 0) THEN
        RETURN NEW;
    end if;
    SELECT is_locked into _is_locked from topics where id=NEW.topic_id;
    IF (_is_locked) THEN
        RETURN NULL;
    end if;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

alter table forums
    add column is_locked bool not null default false;