CREATE TABLE IF NOT EXISTS products(
    id varchar(255) primary key,
    name varchar(60) not null,
    description text not null,
    image_url text,
    category varchar(50) not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    version int not null default 1
);
