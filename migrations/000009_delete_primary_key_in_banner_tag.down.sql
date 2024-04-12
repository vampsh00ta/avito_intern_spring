begin;
drop table banner_tag;
create table banner_tag(
                           banner_id  integer references banner(id) on DELETE cascade ,
                           tag_id  bigint references tag(id) on DELETE cascade ,
                           feature  bigint references feature(id) on DELETE cascade ,
                           primary key (tag_id,feature)
);
drop table banner_tag_history;
create table banner_tag_history(
                                   banner_history_id  integer references banner_history(id) on delete cascade ,
                                   tag_id  bigint references tag(id) ,
                                   feature_id bigint references feature(id),
                                   primary key(banner_history_id,tag_id,feature_id)
);
commit;