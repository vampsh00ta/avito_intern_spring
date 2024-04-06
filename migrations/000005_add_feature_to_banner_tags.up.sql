
begin;
    drop table banner_tag;
    alter table banner drop constraint  banner_feature_id_fkey;
    alter table banner drop column  feature_id;

    create table banner_tag(
        banner_id  integer references banner(id) on DELETE cascade ,
        tag_id  bigint references tag(id) on DELETE cascade ,
        feature  bigint references feature(id) on DELETE cascade ,
        primary key (tag_id,feature)
    );
commit;