
use payment_gateway_lite;

delete from teams;
delete from users;
alter table users AUTO_INCREMENT = 1;
alter table teams AUTO_INCREMENT = 1;