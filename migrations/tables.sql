create table tag(
    id bigint primary key
);
create table feature(
    id bigint primary key
);
create table banner(
    id serial primary key ,
    content text,
    is_active boolean,
    created_at timestamp default  now() not null,
    updated_at timestamp default  now()  not null
);
create table banner_tag(
    banner_id  integer references banner(id) on DELETE cascade ,
    tag_id  bigint references tag(id) on DELETE cascade ,
    feature_id  bigint references feature(id) on DELETE cascade ,
    constraint banner_tag_feature_uc unique (tag_id,feature_id)
);
create table banner_history(
    id serial primary key ,
    banner_id integer references banner(id) on DELETE cascade,
    content text,
    is_active boolean,
    created_at timestamp  not null,
    updated_at timestamp  not null
);

create table banner_tag_history(
    banner_history_id  integer references banner_history(id) on delete cascade ,
    tag_id  bigint references tag(id) ,
    feature_id bigint references feature(id),
    constraint banner_tag_history_feature_uc unique (banner_history_id,tag_id,feature_id)
);

