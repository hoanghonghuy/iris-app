CREATE TABLE IF NOT EXISTS post_comments (
	comment_id      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	post_id         uuid NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
	author_user_id  uuid NOT NULL REFERENCES users(user_id) ON DELETE RESTRICT,
	content         text NOT NULL,
	created_at      timestamptz NOT NULL DEFAULT now(),
	CHECK (length(trim(content)) > 0)
);

CREATE INDEX IF NOT EXISTS idx_post_comments_post_created
	ON post_comments (post_id, created_at DESC);

CREATE TABLE IF NOT EXISTS post_interactions (
	interaction_id  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	post_id         uuid NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
	user_id         uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
	action_type     varchar(20) NOT NULL CHECK (action_type IN ('like', 'share')),
	created_at      timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_post_interactions_like_user
	ON post_interactions (post_id, user_id, action_type)
	WHERE action_type = 'like';

CREATE INDEX IF NOT EXISTS idx_post_interactions_post_action
	ON post_interactions (post_id, action_type, created_at DESC);
