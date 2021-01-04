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
func alarmTargetHandler() []handler {
	handlers := []handler{
		{path: "/alarms/{alarmId}/targets", fun: GetAllAlarmTargets, methods: []string{"GET"}},
		{path: "/alarms/{alarmId}/targets/{targetId}", fun: GetAlarmTargetByTargetID, methods: []string{"GET"}},
		{path: "/alarms/{alarmId}/targets/{targetId}", fun: DeleteAlarmTarget, methods: []string{"PUT"}},
		{path: "/alarms/{alarmId}/targets", fun: AddAlarmTarget, methods: []string{"POST"}},
		{path: "/alarms/{alarmId}/targets", fun: ModifyAlarmTarget, methods: []string{"PUT"}},
	}

	return handlers
}

//	알람 대상 삭제
func DeleteAlarmTarget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetID := vars["targetId"]
	alarmID := vars["alarmId"]

	db := Connect()

	queryString := `
		DELETE FROM alarm_target
		WHERE id = ?
		AND todo_alarm_id = ?
	`

	statement, err := db.Prepare(queryString)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Fail to delete to todos alarm_target.")
	}

	res, err := statement.Exec(targetID, alarmID)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Fail to delete to todos alarm_target.")
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		responseWithJSON(w, http.StatusOK, map[string]int64{"delete:": rowsAffected})
	}

}

//	알람 대상 삭제
func AddAlarmTarget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alarmID := vars["alarmId"]

	decoder := json.NewDecoder(r.Body)
	var alarmTargetParam alarmTarget
	err := decoder.Decode(&alarmTargetParam)
	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusBadRequest, "Couldn't parse request data.")
	}

	db := Connect()

	queryString := `
		INSERT INTO alarm_target (todo_alarm_id,   phone, email, user_id, active_yn, created_at)
		VALUES				     (?,		 	   ?,	  ?,	 ?, 	  ?,		 ?)
	`

	statement, err := db.Prepare(queryString)
	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusInternalServerError, "Fail to prepare query to insert.")
	}
	defer statement.Close()

	alarmTargetParam.CreatedAt = timeToString(time.Now())
	res, err := statement.Exec(alarmID, alarmTargetParam.Phone, alarmTargetParam.Email, alarmTargetParam.UserID, alarmTargetParam.ActiveYn, alarmTargetParam.CreatedAt)
	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusInternalServerError, "Fail to insert todos.")
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		id, _ := res.LastInsertId()
		alarmTargetParam.ID = int64(id)
		responseWithJSON(w, http.StatusOK, alarmTargetParam)
	}
}

//	알람 대상 삭제
func ModifyAlarmTarget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alarmID := vars["alarmId"]

	decoder := json.NewDecoder(r.Body)
	var alarmTargetParam alarmTarget
	err := decoder.Decode(&alarmTargetParam)
	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusBadRequest, "Couldn't parse request data.")
	}

	if alarmTargetParam.ID <= 0 {
		fmt.Println("Not allowed id is null.")
		responseWithError(w, http.StatusBadRequest, "Couldn't parse request data.")
	}

	db := Connect()

	queryString := `
		UPDATE alarm_target SET
			phone=?,
			email=?,
			user_id=?,
			active_yn=?,
			modified_at=?
		WHERE todo_alarm_id = ?
		AND id = ?
	`

	statement, err := db.Prepare(queryString)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Fail to update to todos.")
		return
	}
	defer statement.Close()

	nowTime := timeToString(time.Now())
	alarmTargetParam.ModifiedAt = nowTime

	fmt.Println("todoAlarmParam: ", alarmTargetParam)
	res, err := statement.Exec(alarmTargetParam.Phone, alarmTargetParam.Email,
		alarmTargetParam.UserID, alarmTargetParam.ActiveYn, alarmTargetParam.ModifiedAt,
		alarmID, alarmTargetParam.ID,
	)

	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusInternalServerError, "Fail to update todos.")
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		responseWithJSON(w, http.StatusOK, alarmTargetParam)
	}
}

//	todo 목록 전체 조회
func GetAllAlarmTargets(w http.ResponseWriter, r *http.Request) {

	db := Connect()

	rows, err := db.Query("SELECT * FROM alarm_target;")
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Error occurred when query to db.")
	}
	defer rows.Close()

	var alarmTargetList []alarmTarget

	for rows.Next() {
		todoAlarmEntry := mappingAlarmsTarget(rows)
		alarmTargetList = append(alarmTargetList, todoAlarmEntry)
	}

	responseWithJSON(w, http.StatusOK, alarmTargetList)
}

//	todo 아이디로 조회
func GetAlarmTargetByTargetID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	alarmID := vars["alarmId"]
	targetID := vars["targetId"]

	db := Connect()

	fmt.Println("alarmID: ", alarmID, "targetID:", targetID)
	rows, err := db.Query("SELECT * FROM alarm_target WHERE id = ? AND todo_alarm_id = ?;", targetID, alarmID)
	if err != nil {
		fmt.Println("Error: ", err)
		responseWithError(w, http.StatusInternalServerError, "Error occured when query to db.")
	}
	defer rows.Close()

	var alarmTargetList []alarmTarget

	for rows.Next() {
		alarmTargetEntry := mappingAlarmsTarget(rows)
		alarmTargetList = append(alarmTargetList, alarmTargetEntry)
	}

	responseWithJSON(w, http.StatusOK, alarmTargetList)
}

func mappingAlarmsTarget(rows *sql.Rows) alarmTarget {
	var alarmTargetEntry alarmTarget

	var id sql.NullInt64
	var todoAlarmID sql.NullInt64
	var phone sql.NullString
	var email sql.NullString
	var userID sql.NullInt64
	var activeYN sql.NullString
	var createdAt sql.NullTime
	var modifiedAt sql.NullTime

	rows.Scan(&id, &todoAlarmID, &phone, &email, &userID, &activeYN, &createdAt, &modifiedAt)

	alarmTargetEntry.ID = id.Int64
	alarmTargetEntry.TodoAlarmID = todoAlarmID.Int64
	alarmTargetEntry.Phone = phone.String
	alarmTargetEntry.Email = email.String
	alarmTargetEntry.UserID = userID.Int64
	alarmTargetEntry.ActiveYn = activeYN.String
	alarmTargetEntry.CreatedAt = createdAt.Time.String()
	alarmTargetEntry.ModifiedAt = modifiedAt.Time.String()

	fmt.Println("alarmTargetEntry: ", id, todoAlarmID, phone, email, userID, activeYN, createdAt, modifiedAt)
	return alarmTargetEntry
}
