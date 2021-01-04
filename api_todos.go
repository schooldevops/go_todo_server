package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// todo handler 등록하기
func todoHandler() []handler {
	handlers := []handler{
		{path: "/todos/{id}", fun: GetTodoByID, methods: []string{"GET"}},
		{path: "/todos", fun: GetAllTodos, methods: []string{"GET"}},
		{path: "/todos", fun: CreateTodos, methods: []string{"POST"}},
		{path: "/todos", fun: UpdateTodos, methods: []string{"PUT"}},
		{path: "/todos/{id}", fun: DeleteTodos, methods: []string{"DELETE"}},
	}

	return handlers
}

//	todo 삭제
func DeleteTodos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID := vars["id"]

	db := Connect()

	selectQueryString := `
		SELECT * 
		FROM todos
		WHERE id=?
	`

	deleteQueryString := `
		DELETE FROM todos
		WHERE id=?
	`
	var id sql.NullInt64
	var userID sql.NullInt64
	var title sql.NullString
	var priority sql.NullString
	var status sql.NullString
	var completionLevel sql.NullInt64
	var createdAt sql.NullTime
	var modifiedAt sql.NullTime
	var doneAt sql.NullTime

	err := db.QueryRow(selectQueryString, todoID).Scan(&id, &userID, &title, &priority, &status, &completionLevel, &createdAt, &modifiedAt, &doneAt)

	switch {
	case err == sql.ErrNoRows:
		responseWithError(w, http.StatusBadRequest, "There are not any todos by id="+todoID)
		return
	case err != nil:
		responseWithError(w, http.StatusInternalServerError, "Fail to select todos.")
		return
	default:
		res, err := db.Exec(deleteQueryString, todoID)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, "Fail to delete todos.")
			return
		}
		count, err := res.RowsAffected()
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, "Fail to delete todos.")
			return
		}
		if count == 1 {
			var todoEntry todos
			todoEntry.ID = id.Int64
			todoEntry.UserID = userID.Int64
			todoEntry.TITLE = title.String
			todoEntry.Priority = priority.String
			todoEntry.Status = status.String
			todoEntry.CompletionLevel = completionLevel.Int64
			todoEntry.CreatedAt = createdAt.Time.String()
			todoEntry.ModifiedAt = modifiedAt.Time.String()
			todoEntry.DoneAt = doneAt.Time.String()

			responseWithJSON(w, http.StatusOK, todoEntry)
			return
		}
	}

}

//	todo 업데이트
func UpdateTodos(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var todoParam todos
	err := decoder.Decode(&todoParam)
	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusBadRequest, "Couldn't parse request data.")
	}

	if todoParam.ID <= 0 {
		responseWithError(w, http.StatusBadRequest, "Couldn't parse request data.")
	}

	db := Connect()

	queryString := `
		UPDATE todos SET 
			title=?,
			priority=?,
			status=?,
			completion_level=?,
			modified_at=?,
			done_at=?
		WHERE id=?;
	`
	statement, err := db.Prepare(queryString)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "File to update to todos.")
		return
	}
	defer statement.Close()

	nowTime := timeToString(time.Now())
	todoParam.ModifiedAt = nowTime

	if todoParam.DoneAt == "" {
		todoParam.DoneAt = nowTime
	}

	fmt.Println("todoParam: ", todoParam)

	res, err := statement.Exec(todoParam.TITLE, todoParam.Priority,
		todoParam.Status, todoParam.CompletionLevel,
		todoParam.ModifiedAt, todoParam.DoneAt,
		todoParam.ID)

	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusInternalServerError, "Fail to update todos.")
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		responseWithJSON(w, http.StatusOK, todoParam)
	}
}

//	todo 신규 생성
func CreateTodos(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var todoParam todos
	err := decoder.Decode(&todoParam)
	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusBadRequest, "Couldn't parse request data.")
	}

	// fmt.Println("params: ", todoParam)

	db := Connect()

	queryString := `
		INSERT INTO todos (user_id, title, 	priority, 	status, completion_level, 	created_at)
		VALUES 			  (?, 		?, 		?, 			?, 		?, 					?)
	`

	statement, err := db.Prepare(queryString)
	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusInternalServerError, "Fail to prepare query to insert.")
	}
	defer statement.Close()

	todoParam.CreatedAt = timeToString(time.Now())
	res, err := statement.Exec(todoParam.UserID, todoParam.TITLE, todoParam.Priority, todoParam.Status, todoParam.CompletionLevel, todoParam.CreatedAt)
	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusInternalServerError, "Fail to insert todos.")
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		id, _ := res.LastInsertId()
		todoParam.ID = int64(id)
		responseWithJSON(w, http.StatusOK, todoParam)
	}
}

//	todo 목록 전체 조회
func GetAllTodos(w http.ResponseWriter, r *http.Request) {

	db := Connect()

	rows, err := db.Query("SELECT * FROM todos;")
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error occured when query to db.")
	}
	defer rows.Close()

	var todoList []todos
	for rows.Next() {
		todoEntry := mappingTodos(rows)
		todoList = append(todoList, todoEntry)
	}

	responseWithJSON(w, http.StatusOK, todoList)
}

//	todo 아이디로 조회
func GetTodoByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	todoID := vars["id"]

	db := Connect()

	rows, err := db.Query("SELECT * FROM todos WHERE id = ?;", todoID)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error occured when query to db.")
	}
	defer rows.Close()

	var todoEntry todos
	for rows.Next() {
		todoEntry = mappingTodos(rows)
	}

	responseWithJSON(w, http.StatusOK, todoEntry)
}

func mappingTodos(rows *sql.Rows) todos {
	var todoEntry todos

	var id sql.NullInt64
	var userID sql.NullInt64
	var title sql.NullString
	var priority sql.NullString
	var status sql.NullString
	var completionLevel sql.NullInt64
	var createdAt sql.NullTime
	var modifiedAt sql.NullTime
	var doneAt sql.NullTime

	rows.Scan(&id, &userID, &title, &priority, &status, &completionLevel, &createdAt, &modifiedAt, &doneAt)

	todoEntry.ID = id.Int64
	todoEntry.UserID = userID.Int64
	todoEntry.TITLE = title.String
	todoEntry.Priority = priority.String
	todoEntry.Status = status.String
	todoEntry.CompletionLevel = completionLevel.Int64
	todoEntry.CreatedAt = createdAt.Time.String()
	todoEntry.ModifiedAt = modifiedAt.Time.String()
	todoEntry.DoneAt = doneAt.Time.String()

	fmt.Println("todoEntry: ", id, userID, title, priority, status, completionLevel, createdAt)
	return todoEntry
}
