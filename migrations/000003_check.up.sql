begin;
alter table tag add constraint chk_tag_id check ( id > 0 );
alter table feature add constraint chk_feature_id check ( id > 0 );
commit ;
