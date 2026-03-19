-- Add indexes for foreign keys and common query patterns

CREATE INDEX IF NOT EXISTS idx_topics_board_id ON topics(board_id);
CREATE INDEX IF NOT EXISTS idx_topics_board_id_updated_at ON topics(board_id, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_topic_id ON posts(topic_id);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_boards_forum_id ON boards(forum_id);
