-- create schema version table
create table schema_version (
    version text not null
);

-- create users table
create table users
(
    name text primary key CHECK (name <> ''),
    hash text not null CHECK (hash <> ''),
    about TEXT not null DEFAULT '',
    theme TEXT not null DEFAULT ''
);

-- create posts table
create table posts
(
    id serial primary key,
    author TEXT references users(name) NOT NULL,
    title TEXT NOT NULL CHECK (title <> ''),
    content TEXT NOT NULL CHECK (content <> ''),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

-- create replies table
create table replies
(
    id serial primary key,
    author TEXT references users(name) NOT NULL,
    content TEXT NOT NULL CHECK (content <> ''),
    post_id int references posts(id) NOT NULL,
    parent_id int references replies(id),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

-- create notification table
create table notifications
(
    id serial primary key,
    author text references users(name),
    reply_id int references replies(id)
);

-- indices
CREATE INDEX idx_replies_parent_id ON replies(parent_id);
CREATE INDEX idx_replies_post_id ON replies(post_id);
