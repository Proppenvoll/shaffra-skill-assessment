package main

import "database/sql"

type queryCreateUser func(UserWithoutId) (int, error)

func getQueryCreateUser(db *sql.DB) queryCreateUser {
	return func(user UserWithoutId) (int, error) {
		var userId int

		error := db.QueryRow(
			"INSERT INTO app_user (name, email, age) VALUES ($1, $2, $3) RETURNING app_user_id",
			user.Name,
			user.Email,
			user.Age,
		).Scan(&userId)

		return userId, error
	}
}

type queryGetUser func(id string) (User, error)

func getQueryGetUser(db *sql.DB) queryGetUser {
	return func(id string) (User, error) {
		var user User
		error := db.QueryRow(
			"SELECT * FROM app_user WHERE app_user_id = $1",
			id,
		).Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Age,
		)
		return user, error
	}
}

type queryReplaceUser func(user User) error

func getQueryReplaceUser(db *sql.DB) queryReplaceUser {
	return func(user User) error {
		_, error := db.Exec(
			`INSERT INTO app_user (app_user_id, name, email, age)
                 VALUES ($1, $2, $3, $4)
				 ON CONFLICT (app_user_id) DO UPDATE SET
                 name = $2, email = $3, age = $4`,
			user.Id,
			user.Name,
			user.Email,
			user.Age,
		)
		return error
	}
}

type queryDeleteUser func(id string) error

func getQueryDeleteUser(db *sql.DB) queryDeleteUser {
	return func(id string) error {
		_, error := db.Exec("DELETE FROM app_user WHERE app_user_id = $1", id)
		return error
	}
}
