// WildMutt executes a SQL query on the local shard DB (shard executor)
func WildMutt(w http.ResponseWriter, r *http.Request) {
	var cmd WorkerCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if cmd.Cmd != "EXECUTE_QUERY" {
		http.Error(w, "invalid command", http.StatusBadRequest)
		return
	}

	if cmd.SQL == "" {
		http.Error(w, "sql cannot be empty", http.StatusBadRequest)
		return
	}

	// Try SELECT-like execution first
	rows, err := db.Query(cmd.SQL)
	if err == nil {
		defer rows.Close()

		cols, err := rows.Columns()
		if err != nil {
			sendJSON(w, WorkerResponse{
				Success: false,
				Message: "failed to read columns",
				ShardID: cmd.ShardID,
				Error:   err.Error(),
			})
			return
		}

		results := []map[string]interface{}{}

		for rows.Next() {
			values := make([]interface{}, len(cols))
			valuePtrs := make([]interface{}, len(cols))

			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				sendJSON(w, WorkerResponse{
					Success: false,
					Message: "failed to scan row",
					ShardID: cmd.ShardID,
					Error:   err.Error(),
				})
				return
			}

			rowMap := map[string]interface{}{}
			for i, col := range cols {
				val := values[i]
				if b, ok := val.([]byte); ok {
					rowMap[col] = string(b)
				} else {
					rowMap[col] = val
				}
			}

			results = append(results, rowMap)
		}

		sendJSON(w, WorkerResponse{
			Success: true,
			Message: "query executed successfully",
			ShardID: cmd.ShardID,
			Rows:    results,
		})
		return
	}

	// Otherwise treat it as INSERT/UPDATE/DELETE
	res, execErr := db.Exec(cmd.SQL)
	if execErr != nil {
		sendJSON(w, WorkerResponse{
			Success: false,
			Message: "query failed",
			ShardID: cmd.ShardID,
			Error:   execErr.Error(),
		})
		return
	}

	affected, _ := res.RowsAffected()

	sendJSON(w, WorkerResponse{
		Success:      true,
		Message:      "query executed successfully",
		ShardID:      cmd.ShardID,
		AffectedRows: affected,
	})
}