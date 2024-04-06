begin;
drop table banner_tag;


alter table banner add column  feature_id bigint;
alter table banner add constraint  banner_feature_id_fkey
    foreign key (feature_id) references feature(id) on DELETE cascade;

create table banner_tag (
    banner_id  integer references banner(id) on DELETE cascade ,
    tag_id  bigint references tag(id) on DELETE cascade ,
    primary key (banner_id, tag_id)
);


commit;