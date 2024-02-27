```sql
create database if not exists gptWeb;
use gptWeb;
create table if not exists user(
    id int unique auto_increment,
    name varchar(255) primary key,
    age int,
    password varchar(255) not null
);
create index name_password on user(name,password);
```
