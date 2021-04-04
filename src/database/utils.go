package postgreRepo

func (repo *PostgreRepo) IsEmailUniq(userID int, email string) bool {
	return true
}
