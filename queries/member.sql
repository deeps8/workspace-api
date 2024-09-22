DROP TABLE IF EXISTS member;

CREATE TABLE member (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    space_id UUID NOT NULL REFERENCES workspace(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id,space_id)
)