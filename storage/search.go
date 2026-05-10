package storage

import (
	"database/sql"
	"log"
	"vpub/model"
)

func (s *Storage) Search(query string) ([]model.Search, error) {
	var searchResults []model.Search

	rows, err := s.db.Query(`
		SELECT origin_table, id, title, highlighted_title, highlighted_content, rank, forum_group_id
		FROM search_with_highlights($1)
	`, query)
	if err != nil {
		return []model.Search{}, err
	}

	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}(rows)

	for rows.Next() {
		var search model.Search
		var groupID sql.NullInt64
		if err := rows.Scan(
			&search.OriginTable,
			&search.ID,
			&search.Title,
			&search.HighlightedTitle,
			&search.HighlightedContent,
			&search.Rank,
			&groupID,
		); err != nil {
			return []model.Search{}, err
		}
		if groupID.Valid {
			v := groupID.Int64
			search.ForumGroupID = &v
		}
		searchResults = append(searchResults, search)
	}

	return searchResults, err
}
