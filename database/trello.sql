create table if not exists USERS (
  id serial primary key,
  email text unique,
  username text,
  hashed_password text,
  path_to_avatar text
);

create table if not exists BOARDS (
  id serial primary key,
  name text,
  description text
);

create table if not exists ROWS (
  id serial primary key,
  name text,
  position integer not null
);

create table if not exists TASKS (
  id serial primary key,
  name text,
  description text,
  position integer not null,
  check (name != '')
);

create table if not exists Owners (
  user_id integer not null,
  board_id integer not null,

  foreign key (user_id) references USERS (id) on delete cascade,
  foreign key (board_id) references BOARDS (id) on delete cascade
);

create table if not exists Administrators (
  user_id integer not null,
  board_id integer not null,

  foreign key (user_id) references USERS (id) on delete cascade,
  foreign key (board_id) references BOARDS (id) on delete cascade
);

create table if not exists Members (
  user_id integer not null,
  board_id integer not null,

  foreign key (user_id) references USERS (id) on delete cascade,
  foreign key (board_id) references BOARDS (id) on delete cascade
);

create table if not exists Boards_Rows (
  board_id integer not null,
  row_id integer not null,

  foreign key (board_id) references BOARDS (id) on delete cascade,
  foreign key (row_id) references ROWS (id) on delete cascade
);

create table if not exists Rows_Tasks (
  row_id integer not null,
  task_id integer not null,

  foreign key (row_id) references ROWS (id) on delete cascade,
  foreign key (task_id) references TASKS (id) on delete cascade
);

create table if not exists Tasks_Users (
  task_id integer not null,
  user_id integer not null,

  foreign key (task_id) references TASKS (id) on delete cascade,
  foreign key (user_id) references USERS (id) on delete cascade
);

create table if not exists Attachments (
  id serial primary key,
  file_name text,
  path text,

  task_id integer not null,
  foreign key (task_id) references TASKS (id) on delete cascade
);



create table if not exists Sessions (
  session_token text not null unique,
  user_id integer not null,

  foreign key (user_id) references USERS (id) on delete cascade
);
