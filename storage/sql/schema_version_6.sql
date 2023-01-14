CREATE OR REPLACE VIEW topics_summary as
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
order by t.is_sticky desc, t.updated_at desc;