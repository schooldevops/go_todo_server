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
func todoAlarmHandler() []handler {
	handlers := []handler{
		{path: "/todos/alarms", fun: GetAllAlarms, methods: []string{"GET"}},
		{path: "/todos/{id}/alarms", fun: GetAllTodoAlarmsByTodoId, methods: []string{"GET"}},
		{path: "/todos/{id}/alarms/{alarmId}", fun: GetTodoAlarmsById, methods: []string{"GET"}},
		{path: "/todos/{id}/alarms/{alarmId}", fun: DeleteTodoAlarmsById, methods: []string{"DELETE"}},
		{path: "/todos/{id}/alarms", fun: CreateTodoAlarm, methods: []string{"POST"}},
		{path: "/todos/{id}/alarms", fun: UpdateTodoAlarm, methods: []string{"PUT"}},
	}

	return handlers
}

//	Alarm 삭제한다.
func DeleteTodoAlarmsById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID := vars["id"]
	alarmID := vars["alarmId"]

	db := Connect()

	queryString := `
		DELETE FROM todo_alarm
		WHERE id = ?
		AND todo_id = ?
	`

	statement, err := db.Prepare(queryString)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Fail to delete to todos alarm.")
	}

	res, err := statement.Exec(alarmID, todoID)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Fail to delete to todos alarm.")
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		responseWithJSON(w, http.StatusOK, map[string]int64{"delete:": rowsAffected})
	}

}

//	todoAlarm 업데이트 하기.
func UpdateTodoAlarm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID := vars["id"]

	decoder := json.NewDecoder(r.Body)
	var todoAlarmParam todoAlarm
	err := decoder.Decode(&todoAlarmParam)
	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusBadRequest, "Couldn't parse request data.")
	}

	if todoAlarmParam.ID <= 0 {
		fmt.Println("Not allowed id is null.")
		responseWithError(w, http.StatusBadRequest, "Couldn't parse request data.")
	}

	db := Connect()

	queryString := `
		UPDATE todo_alarm SET
			period_type=?,
			alarm_type=?,
			alarm_date=?,
			alarm_time=?,
			active_yn=?,
			modified_at=?
		WHERE todo_id = ?
		AND id = ?
	`

	statement, err := db.Prepare(queryString)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Fail to update to todos.")
		return
	}
	defer statement.Close()

	nowTime := timeToString(time.Now())
	todoAlarmParam.ModifiedAt = nowTime

	fmt.Println("todoAlarmParam: ", todoAlarmParam)
	res, err := statement.Exec(todoAlarmParam.PeriodType, todoAlarmParam.AlarmType,
		todoAlarmParam.AlarmDate, todoAlarmParam.AlarmTime, todoAlarmParam.ActiveYn,
		todoAlarmParam.ModifiedAt, todoID, todoAlarmParam.ID)

	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusInternalServerError, "Fail to update todos.")
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		responseWithJSON(w, http.StatusOK, todoAlarmParam)
	}
}

//	todoAlarm 생성하기.
func CreateTodoAlarm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID := vars["id"]

	decoder := json.NewDecoder(r.Body)
	var todoAlarmParam todoAlarm
	err := decoder.Decode(&todoAlarmParam)
	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusBadRequest, "Couldn't parse request data.")
	}

	db := Connect()

	queryString := `
		INSERT INTO todo_alarm (todo_id, period_type, alarm_type, alarm_date, alarm_time, active_yn, created_at)
		VALUES				   (?,		 ?,			  ?,		  ?, 		  ?,		  ?, 		 ?)
	`

	statement, err := db.Prepare(queryString)
	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusInternalServerError, "Fail to prepare query to insert.")
	}
	defer statement.Close()

	todoAlarmParam.CreatedAt = timeToString(time.Now())
	res, err := statement.Exec(todoID, todoAlarmParam.PeriodType, todoAlarmParam.AlarmType, todoAlarmParam.AlarmDate, todoAlarmParam.AlarmTime, todoAlarmParam.ActiveYn, todoAlarmParam.CreatedAt)
	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusInternalServerError, "Fail to insert todos.")
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		id, _ := res.LastInsertId()
		todoAlarmParam.ID = int64(id)
		responseWithJSON(w, http.StatusOK, todoAlarmParam)
	}
}

//	todo alarm 을 아이디로 조회한다.
func GetTodoAlarmsById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID := vars["id"]
	todoAlarmID := vars["alarmId"]

	db := Connect()

	rows, err := db.Query("SELECT * FROM todo_alarm WHERE id = ? AND todo_id = ?;", todoAlarmID, todoID)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error occured when query to db.")
	}
	defer rows.Close()

	var todoAlarmEntry todoAlarm
	for rows.Next() {
		todoAlarmEntry = mappingTodoAlarms(rows)
	}

	responseWithJSON(w, http.StatusOK, todoAlarmEntry)
}

//	todo 목록 전체 조회
func GetAllAlarms(w http.ResponseWriter, r *http.Request) {

	db := Connect()

	rows, err := db.Query("SELECT * FROM todo_alarm;")
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error occurred when query to db.")
	}
	defer rows.Close()

	var todoAlarmList []todoAlarm

	for rows.Next() {
		todoAlarmEntry := mappingTodoAlarms(rows)
		todoAlarmList = append(todoAlarmList, todoAlarmEntry)
	}

	responseWithJSON(w, http.StatusOK, todoAlarmList)
}

//	todo 아이디로 조회
func GetAllTodoAlarmsByTodoId(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	todoID := vars["id"]

	db := Connect()

	rows, err := db.Query("SELECT * FROM todo_alarm WHERE todo_id = ?;", todoID)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error occured when query to db.")
	}
	defer rows.Close()

	var todoAlarmList []todoAlarm

	for rows.Next() {
		todoAlarmEntry := mappingTodoAlarms(rows)
		todoAlarmList = append(todoAlarmList, todoAlarmEntry)
	}

	responseWithJSON(w, http.StatusOK, todoAlarmList)
}

func mappingTodoAlarms(rows *sql.Rows) todoAlarm {
	var todoAlarmEntry todoAlarm

	var id sql.NullInt64
	var todoID sql.NullInt64
	var periodType sql.NullString
	var alarmType sql.NullString
	var alarmDate sql.NullString
	var alarmTime sql.NullString
	var activeYN sql.NullString
	var createdAt sql.NullTime
	var modifiedAt sql.NullTime

	rows.Scan(&id, &todoID, &periodType, &alarmType, &alarmDate, &alarmTime, &activeYN, &createdAt, &modifiedAt)

	todoAlarmEntry.ID = id.Int64
	todoAlarmEntry.TodoID = todoID.Int64
	todoAlarmEntry.PeriodType = periodType.String
	todoAlarmEntry.AlarmType = alarmType.String
	todoAlarmEntry.AlarmDate = alarmDate.String
	todoAlarmEntry.AlarmTime = alarmTime.String
	todoAlarmEntry.ActiveYn = activeYN.String
	todoAlarmEntry.CreatedAt = createdAt.Time.String()
	todoAlarmEntry.ModifiedAt = modifiedAt.Time.String()

	fmt.Println("todoAlarm: ", id, todoID, periodType, alarmType, alarmDate, alarmTime, activeYN, createdAt, modifiedAt)
	return todoAlarmEntry
}
