-- posts テーブル
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    user_name VARCHAR(32) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- post_tags テーブル
CREATE TABLE post_tags (
    post_id UUID NOT NULL,
    name VARCHAR(16) NOT NULL,
    PRIMARY KEY (post_id, name),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

-- post_reactions テーブル
CREATE TABLE post_reactions (
    post_id UUID NOT NULL,
    reaction_id TINYINT NOT NULL,
    user_name VARCHAR(32) NOT NULL,
    created_at DATETIME,
    PRIMARY KEY (post_id, reaction_id, user_name),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);