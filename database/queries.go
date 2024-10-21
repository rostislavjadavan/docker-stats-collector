package database

func (d *Database) GetImages() ([]string, error) {
	rows, err := d.Conn.Query(`SELECT distinct image from container_stats`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		results = append(results, name)
	}

	return results, nil
}

func (d *Database) CleanupOldData() (int64, error) {
	result, err := d.Conn.Exec("DELETE FROM container_stats WHERE timestamp < date('now', '-30 days')")
	if err != nil {
		return -1, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	if rowsAffected > 0 {
		_, err = d.Conn.Exec("VACUUM")
		if err != nil {
			return -1, err
		}
	}

	return rowsAffected, nil
}
