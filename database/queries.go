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

//type ContainerUsage struct {
//	ContainerId   string
//	ContainerName string
//	Image         string
//	MemoryUsage   int64
//	MemoryLimit   int64
//	Timestamp     time.Time
//}

//func (d *Database) GetUsage(fromTs time.Time, toTs time.Time) {
//	query := `
//		SELECT id, name, image, timestamp, memory_usage, memory_limit FROM container_stats
//		WHERE timestamp between ? and ?
//	`
//	rows, err := d.Conn.Query(query, fromTs, toTs)
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var results []ContainerUsage
//	for rows.Next() {
//		var data ContainerUsage
//		err := rows.Scan(
//			&data.ContainerId,
//			&data.ContainerName,
//			&data.Image,
//			&data.Timestamp,
//			&data.MemoryUsage,
//			&data.MemoryLimit,
//		)
//		if err != nil {
//			return nil, err
//		}
//		results = append(results, data)
//	}
//
//	if err = rows.Err(); err != nil {
//		return nil, err
//	}
//
//	return results, nil
//}
