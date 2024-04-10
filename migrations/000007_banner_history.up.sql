begin;
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
                     primary key(banner_history_id,tag_id,feature_id)
);
commit ;