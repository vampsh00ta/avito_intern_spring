begin;
insert into tag (id) values (1),(2),(3),(4) on  conflict do nothing;
insert into feature (id) values (1),(2),(3),(4) on  conflict do nothing;
commit;

