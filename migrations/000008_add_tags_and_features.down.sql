begin;
delete from tag where id in (1,2,3,4);
delete from feature where id in (1,2,3,4);
commit;

