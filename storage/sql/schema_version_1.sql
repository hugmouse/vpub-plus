-- create schema version table
create table schema_version (
    version text not null
);

create table users (
    id integer primary key autoincrement,
    name text unique CHECK (name <> '' and length(name) <= 15),
    hash text not null CHECK (hash <> ''),
    about TEXT not null DEFAULT '',
    is_admin boolean default false
);

create table settings (
    name text not null,
    css text not null default ''
);

insert into settings (name) values ('vpub');

create table boards (
    id integer primary key autoincrement,
    name text not null check ( name <> '' and length(name) < 120 ),
    topics integer not null default 0,
    posts integer not null default 0,
    description text,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

create table topics (
    id integer primary key autoincrement,
    board_id integer,
    first_post_id integer not null,
    replies integer not null default 0,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    foreign key (first_post_id) references posts(id),
    foreign key (board_id) references boards(id)
);

create table posts (
    id integer primary key autoincrement,
    user_id text not null,
    subject text not null check ( length(subject) <= 120 ),
    content text not null check ( length(content) <= 50000 ),
    topic_id integer,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    foreign key (topic_id) references topics(id),
    foreign key (user_id) references users(id)
);