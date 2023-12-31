create table teams(
    id bigint primary key auto_increment,
    owner_id bigint default null unique,
    name varchar(255) not null,
    logo varchar(255) default null,
    address varchar(255) default null,
    callback_url varchar(255) default null,
    return_url varchar(255) default null,
    bank_id varchar(255) default null,
    status enum('unverified','verified','pending','rejected') default 'unverified',
    status_disbursement enum('unverified','verified','pending','rejected') default 'unverified',
    rate enum('REGULER','CUSTOM') default 'REGULER',
    uuid varchar(36) default null,
    secret varchar(255) default null,
    ssl_verification tinyint default 0,
    other_bank varchar(255) default null,
    settle_time enum('DEFAULT','CUSTOM') default 'DEFAULT',
    fee_charged_to enum('MERCHANT','USER') default 'MERCHANT',
    campaign_name varchar(255) default null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp,
    deleted_at timestamp default null,

    foreign key (owner_id) references users(id)
)engine = InnoDB;