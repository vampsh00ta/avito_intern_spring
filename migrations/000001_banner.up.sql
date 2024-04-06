begin;
create table tag(
    id bigint primary key
);
create table feature(
    id bigint primary key
);
create table banner(
    id serial primary key ,
    feature_id bigint references feature(id) on DELETE cascade,
    content text,
    is_active boolean,
    created_at timestamp default  now() not null,
    updated_at timestamp default  now()  not null
);

create table banner_tag(
       banner_id  integer references banner(id) on DELETE cascade ,
       tag_id  bigint references tag(id) on DELETE cascade ,
       primary key (banner_id, tag_id)
);
commit ;