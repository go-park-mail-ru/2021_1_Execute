create table if not exists USERS (
  id serial primary key,
  email text,
  username text,
  hashed_password text,
  path_to_avatar text
);

create table if not exists BOARDS (
  id integer not null primary key,
  name text,
  description text
);

create table if not exists ROWS (
  id integer not null primary key,
  name text,
  position integer
);

create table if not exists TASKS (
  id integer not null primary key,
  name text,
  description text
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
