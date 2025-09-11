create table users (
   id           uuid    primary key not null,
   emailaddress varchar unique      default '',
   name         varchar             default '',
   enabled      boolean             default true,
   created_at   integer             default (strftime('%s', 'now'))
);