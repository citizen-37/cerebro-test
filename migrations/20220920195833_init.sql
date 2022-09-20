-- +goose Up
-- +goose StatementBegin
create table users
(
    id         uuid primary key,
    name       varchar(255) not null,
    balance    bigint default 0,
    created_at timestamptz  not null,
    updated_at timestamptz  not null,
    check ( balance >= 0 )
);
create unique index users_name_uniq on users (name);
create index users_name_idx on users (name);

create table transactions
(
    id             uuid primary key,
    user_id        uuid references users (id),
    external_id    varchar(266) not null,
    amount         bigint      default 0,
    created_at     timestamptz default now(),
    rolled_back_at timestamptz  null,
    payload        jsonb        null
);
create unique index transactions_uniq on transactions (user_id, external_id);
create index transactions_external_id_idx on transactions (external_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table transactions;
drop table users;
-- +goose StatementEnd
