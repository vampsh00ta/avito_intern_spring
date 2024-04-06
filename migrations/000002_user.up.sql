
begin;
create table customer(
    id serial primary key ,
    username varchar(255),
    admin boolean
);
insert into customer(username,admin) values ('notadmin',false), ('admin',true);
commit ;