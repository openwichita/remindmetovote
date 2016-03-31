create table subscriptions (
  id serial,
  number varchar(12) not null
  CONSTRAINT subscriptions_number_key UNIQUE (number)
);

create table message_log (
  id serial,
  from_number varchar(12) not null,
  to_number varchar(12) not null,
  body text not null,
  created_at timestamp not null
);

ALTER TABLE subscriptions
ADD CONSTRAINT subscriptions_number_key UNIQUE (number)