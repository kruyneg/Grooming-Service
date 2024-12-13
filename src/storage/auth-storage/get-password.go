package authstorage

func (s *AuthStorage) GetPassword(login, role string) (int64, string, error) {
	var (
		id int64
		passHash string
	)
	row := s.db.QueryRow(`
		SELECT user_id, password_hash
		FROM auth
		WHERE login = $1 AND 
			role = $2
	`, login, role)
	if row.Err() != nil {
		return 0, "", row.Err()
	}
	row.Scan(&id, &passHash)
	return id, passHash, nil
}