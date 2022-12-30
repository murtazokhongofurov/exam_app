CREATE TABLE IF NOT EXISTS review (
    id UUID NOT NULL PRIMARY KEY,
    post_id UUID NOT NULL, 
    owner_id UUID NOT NULL,
    name TEXT NOT NULL, 
    rating INTEGER NOT NULL CHECK(rating >= 1 AND rating <= 5),
    description TEXT NOT NULL, 
    created_at TIME NOT NULL DEFAULT NOW(), 
    updated_at TIME NOT NULL DEFAULT NOW(),
    deleted_at TIME
);
