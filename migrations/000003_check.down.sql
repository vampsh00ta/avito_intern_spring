begin;
alter table tag drop constraint chk_tag_id ;
alter table feature drop constraint chk_feature_id ;
commit;