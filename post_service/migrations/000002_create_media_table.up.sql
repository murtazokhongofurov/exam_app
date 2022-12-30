CREATE TABLE IF NOT EXISTS medias(
    id UUID PRIMARY KEY NOT NULL,
    post_id UUID NOT NULL REFERENCES posts(id),
    name TEXT NOT NULL,
    link TEXT NOT NULL,
    type TEXT NOT NULL
);