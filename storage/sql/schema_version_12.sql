-- Replacing to_tsquery with websearch_to_tsquery, fixing issues with spaces and such
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
    query tsquery := websearch_to_tsquery('english', search_term);
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
