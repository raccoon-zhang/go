```sql
create database if not exists school;
use school;
create table if not exists student(
    id int unique auto_increment,
    name varchar(255) primary key,
    age int,
    password varchar(255) not null
);
create index name_password on student(name,password);
```
