package server

import (
	"log"
)

func RunQuery(query string, params ...interface{}) ([]map[string]interface{}, error) {

	rows, err := DB.Query(query, params...)
	if err != nil {
		log.Printf("Erreur lors de l'exécution de la requête : %v", err)
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	var results []map[string]interface{}
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Printf("Erreur lors du scan des résultats : %v", err)
			return nil, err
		}
		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			row[col] = val
		}
		results = append(results, row)
	}

	return results, nil
}
