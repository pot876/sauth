create table "credentials" (
    "realm_id" uuid not null default gen_random_uuid(),
    "user_id" uuid not null default gen_random_uuid(),
    "pwdhash" varchar(255) not null default '',
    "login" varchar(255) not null default '',
    "role" varchar(255) not null default '' 
);
create unique index credentials_realmid_login_unique_idx on credentials (realm_id, login);