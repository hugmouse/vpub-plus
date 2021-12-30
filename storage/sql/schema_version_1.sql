-- create schema version table
create table schema_version (
    version text not null
);

-- create users table
create table users
(
    name text primary key CHECK (name <> ''),
    hash text not null CHECK (hash <> ''),
    about TEXT not null DEFAULT ''
);

-- create posts table
create table posts
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    author TEXT references users(name) NOT NULL,
    title TEXT NOT NULL CHECK (title <> ''),
    content TEXT NOT NULL CHECK (content <> ''),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

-- create replies table
create table replies
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    author TEXT references users(name) NOT NULL,
    content TEXT NOT NULL CHECK (content <> ''),
    post_id int references posts(id) NOT NULL,
    parent_id int references replies(id),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    foreign key (post_id) references posts(id) ON DELETE RESTRICT,
    foreign key (parent_id) references replies(id) ON DELETE RESTRICT
);

-- create notification table
create table notifications
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    author text references users(name),
    reply_id int references replies(id),
    foreign key (author) references users(name) ON DELETE RESTRICT,
    foreign key (reply_id) references replies(id) ON DELETE RESTRICT
);

-- indices
CREATE INDEX idx_replies_parent_id ON replies(parent_id);
CREATE INDEX idx_replies_post_id ON replies(post_id);
