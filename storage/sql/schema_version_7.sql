CREATE OR REPLACE PROCEDURE remove_board(_bid int) as
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
$$ LANGUAGE plpgsql;