CREATE TABLE groups (
    id   SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL CHECK (name <> '' AND length(name) <= 50)
);

CREATE TABLE group_members (
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id  INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, user_id)
);

CREATE INDEX idx_group_members_user_id ON group_members(user_id);

ALTER TABLE forums
    ADD COLUMN group_id              INTEGER REFERENCES groups(id) ON DELETE SET NULL,
    ADD COLUMN restricted_visibility TEXT NOT NULL DEFAULT 'hidden'
        CHECK (restricted_visibility IN ('hidden', 'visible'));

CREATE OR REPLACE VIEW forums_summary AS
SELECT b.id   AS board_id,
       f.id   AS forum_id,
       f.name AS forum_name,
       b.name AS board_name,
       b.description,
       b.topics_count,
       b.posts_count,
       b.updated_at,
       f.group_id,
       f.restricted_visibility
FROM boards b
         LEFT JOIN forums f ON f.id = b.forum_id
ORDER BY f.position, b.position, f.id;

CREATE OR REPLACE VIEW search_items AS
    SELECT
        text 'users' AS origin_table,
        id::text,
        name AS title,
        about AS content,
        textsearchable_index_col AS searchable_element,
        NULL::integer AS forum_group_id
    FROM users
UNION ALL
    SELECT
        text 'posts' AS origin_table,
        concat(p.topic_id, '#', p.id)::text AS id,
        p.subject AS title,
        p.content,
        p.textsearchable_index_col AS searchable_element,
        f.group_id AS forum_group_id
    FROM posts p
    JOIN topics t ON t.id = p.topic_id
    JOIN boards b ON b.id = t.board_id
    JOIN forums f ON f.id = b.forum_id
UNION ALL
    SELECT
        text 'boards' AS origin_table,
        b.id::text,
        b.name AS title,
        b.description AS content,
        b.textsearchable_index_col AS searchable_element,
        f.group_id AS forum_group_id
    FROM boards b
    JOIN forums f ON f.id = b.forum_id;

DROP FUNCTION IF EXISTS search_with_highlights(text);

CREATE OR REPLACE FUNCTION search_with_highlights(search_term text)
    RETURNS TABLE (
        origin_table       text,
        id                 text,
        title              text,
        content            text,
        highlighted_title  text,
        highlighted_content text,
        rank               float4,
        forum_group_id     integer
    ) AS $$
DECLARE
    query tsquery := websearch_to_tsquery('english', search_term);
BEGIN
    RETURN QUERY
        SELECT
            si.origin_table,
            si.id::text,
            si.title,
            si.content,
            ts_headline('english', si.title, query,
                'StartSel=<mark>, StopSel=</mark>, MaxFragments=1') AS highlighted_title,
            ts_headline('english', si.content, query,
                'StartSel=<mark>, StopSel=</mark>, MaxFragments=3, FragmentDelimiter=..., MaxWords=13, MinWords=3') AS highlighted_content,
            ts_rank(si.searchable_element, query) AS rank,
            si.forum_group_id
        FROM search_items si
        WHERE query @@ si.searchable_element
        ORDER BY rank DESC;
END;
$$ LANGUAGE plpgsql;
