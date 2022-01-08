-- create schema version table
create table schema_version (
    version text not null
);

create table users (
    name text primary key CHECK (name <> '' and length(name) <= 15),
    hash text not null CHECK (hash <> ''),
    about TEXT not null DEFAULT ''
);

create table boards (
    id integer primary key autoincrement,
    name text,
    description text
);

create table topics (
    id integer primary key autoincrement,
    board_id integer,
    first_post_id integer not null,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    foreign key (first_post_id) references posts(id),
    foreign key (board_id) references boards(id)
);

create table posts (
    id integer primary key autoincrement,
    author text not null,
    subject text not null check ( length(subject) < 120 ),
    content text not null check ( length(content) < 50000 ),
    topic_id integer,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    foreign key (topic_id) references topics(id),
    foreign key (author) references users(name)
)