package postgresRepo

func (repo *PostgresRepo) IsEmailUniq(userID int, email string) bool {
	return true
}
