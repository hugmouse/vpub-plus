package storage

//
//import (
//	"context"
//	"database/sql"
//	"vpub/model"
//)
//
//func (s *Storage) populateReply(rows *sql.Rows) (model.Reply, error) {
//	var reply model.Reply
//	var createdAtStr string
//	err := rows.Scan(&reply.Id, &reply.User, &reply.Content, &reply.PostId, &reply.ParentId, &reply.Comments, &createdAtStr)
//	reply.CreatedAt, err = parseCreatedAt(createdAtStr)
//	if err != nil {
//		return reply, err
//	}
//	return reply, nil
//}
//
//func (s *Storage) CreateReply(reply model.Reply) (int64, error) {
//	var lid int64
//	ctx := context.Background()
//	tx, err := s.db.BeginTx(ctx, nil)
//	if err != nil {
//		return lid, err
//	}
//	if err := tx.QueryRowContext(ctx, `INSERT INTO replies (author, content, post_id, parent_id) VALUES ($1, $2, $3, $4) RETURNING id`,
//		reply.User, reply.Content, reply.PostId, reply.ParentId).Scan(&lid); err != nil {
//		tx.Rollback()
//		return lid, err
//	}
//	var author string
//	if reply.ParentId != nil {
//		parent, err := s.ReplyById(*reply.ParentId)
//		if err != nil {
//			tx.Rollback()
//			return lid, nil
//		}
//		author = parent.User
//	} else {
//		parent, err := s.PostById(reply.PostId)
//		if err != nil {
//			tx.Rollback()
//			return lid, nil
//		}
//		author = parent.User
//		if _, err := tx.ExecContext(ctx, `update posts set updated_at=datetime('now') where id=$1`, reply.PostId); err != nil {
//			tx.Rollback()
//			return lid, err
//		}
//	}
//	if author != reply.User {
//		if _, err := tx.ExecContext(ctx, `INSERT into notifications (author, reply_id) values ($1, $2)`, author, lid); err != nil {
//			tx.Rollback()
//			return lid, err
//		}
//	}
//	err = tx.Commit()
//	return lid, err
//}
//
//func (s *Storage) RepliesByUsername(name string, page int64) ([]model.Reply, bool, error) {
//	q := `
//		SELECT
//    		r.id, r.author, r.content, r.post_id, r.parent_id, p.title from replies r
//		LEFT JOIN posts p on p.id = r.post_id
//		WHERE r.author=$1
//		ORDER BY r.id DESC
//		LIMIT 11
//		OFFSET $2`
//	rows, err := s.db.Query(q, name, page*10)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return nil, false, nil
//		}
//		return nil, false, err
//	}
//	var replies []model.Reply
//	for rows.Next() {
//		var reply model.Reply
//		err := rows.Scan(&reply.Id, &reply.User, &reply.Content, &reply.PostId, &reply.ParentId, &reply.PostTitle)
//		if err != nil {
//			return replies, false, err
//		}
//		replies = append(replies, reply)
//	}
//	if len(replies) > 10 {
//		return replies[0:10], true, err
//	}
//	return replies, false, nil
//}
//
//func (s *Storage) RepliesByPostId(postId int64) ([]model.Reply, error) {
//	q := `
//		SELECT
//		   r.id, r.author, r.content, r.post_id, r.parent_id, coalesce(count(s.id), 0), r.created_at from replies r
//		LEFT JOIN replies s ON r.id = s.parent_id
//		WHERE r.post_id=$1 AND r.parent_id IS NULL
//		GROUP BY r.id
//		ORDER BY r.id`
//	rows, err := s.db.Query(q, postId)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return nil, nil
//		}
//		return nil, err
//	}
//	var replies []model.Reply
//	for rows.Next() {
//		var reply model.Reply
//		var count int
//		var createdAtStr string
//		rows.Scan(&reply.Id, &reply.User, &reply.Content, &reply.PostId, &reply.ParentId, &count, &createdAtStr)
//		reply.CreatedAt, err = parseCreatedAt(createdAtStr)
//		if count > 0 {
//			reply.Thread, err = s.RepliesByParentId(reply.Id)
//			if err != nil {
//				return replies, err
//			}
//		}
//		replies = append(replies, reply)
//	}
//	return replies, nil
//}
//
//func (s *Storage) RepliesByParentId(parentId int64) ([]model.Reply, error) {
//	q := `
//		SELECT
//			r.id, r.author, r.content, r.post_id, r.parent_id, coalesce(count(s.id), 0), r.created_at from replies r
//		LEFT JOIN replies s ON r.id = s.parent_id
//		WHERE r.parent_id=$1
//		GROUP BY r.id
//		ORDER BY r.id DESC`
//	rows, err := s.db.Query(q, parentId)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return nil, nil
//		}
//		return nil, err
//	}
//	var replies []model.Reply
//	for rows.Next() {
//		reply, err := s.populateReply(rows)
//		if err != nil {
//			return replies, err
//		}
//		if reply.Comments > 0 {
//			reply.Thread, err = s.RepliesByParentId(reply.Id)
//			if err != nil {
//				return replies, err
//			}
//		}
//		replies = append(replies, reply)
//	}
//	return replies, nil
//}
//
//func (s *Storage) ReplyById(id int64) (model.Reply, error) {
//	var reply model.Reply
//	var createdAtStr string
//	err := s.db.QueryRow(
//		`SELECT id, author, content, post_id, parent_id, created_at from replies WHERE id=$1`, id).Scan(
//		&reply.Id,
//		&reply.User,
//		&reply.Content,
//		&reply.PostId,
//		&reply.ParentId,
//		&createdAtStr,
//	)
//	reply.CreatedAt, err = parseCreatedAt(createdAtStr)
//	return reply, err
//}
//
//func (s *Storage) UpdateReply(reply model.Reply) error {
//	stmt, err := s.db.Prepare(`UPDATE replies SET content = $1 WHERE id = $2;`)
//	if err != nil {
//		return err
//	}
//	_, err = stmt.Exec(reply.Content, reply.Id)
//	return nil
//}
//
//func (s *Storage) DeleteReply(id int64) error {
//	stmt, err := s.db.Prepare(`DELETE from replies WHERE id = $1;`)
//	if err != nil {
//		return err
//	}
//	_, err = stmt.Exec(id)
//	return err
//}
