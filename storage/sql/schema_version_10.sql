ALTER TABLE boards
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
