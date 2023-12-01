create table users(
    id bigint primary key auto_increment,
    name varchar(255) not null,
    email varchar(255) not null unique,
    email_verified_at timestamp default null,
    password varchar(255) not null,
    phone varchar(255)  default null,
    team_id int default null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp,
    deleted_at timestamp default null
)engine=InnoDB;