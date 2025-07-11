-- name: CreateUser :one
insert into users (login, pwd_hash) values ($1, $2) returning id;

-- name: GetUserByID :one
select * from users where id = $1;

-- name: DeleteUserByID :execrows
delete from users where id = $1;

-- name: UpdateUserPwd :one
update users set pwd_hash = $2 where id = $1 returning *;