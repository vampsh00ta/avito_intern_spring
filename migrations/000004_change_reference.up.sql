begin;
alter table banner_tag drop constraint banner_tag_tag_id_fkey;
alter table banner_tag  add constraint banner_tag_tag_id_fkey
   foreign key (tag_id) references tag(id) on DELETE cascade;
commit ;