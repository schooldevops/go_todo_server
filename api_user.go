package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// todo handler 등록하기
func userHandler() []handler {
	handlers := []handler{
		{path: "/users/{id}", fun: GetUserByID, methods: []string{"GET"}},
		{path: "/users", fun: GetAllUsers, methods: []string{"GET"}},
		{path: "/users", fun: CreateUsers, methods: []string{"POST"}},
		{path: "/users", fun: UpdateUsers, methods: []string{"PUT"}},
		{path: "/users/{id}", fun: DeleteUsers, methods: []string{"DELETE"}},
	}

	return handlers
}

//	GetUserByID todo 아이디로 조회
func GetUserByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	userID := vars["id"]

	var userEntity User
	db := ConnectGormDB()

	db.Where("id = ?", userID).First(&userEntity)

	responseWithJSON(w, http.StatusOK, convertUserEntityToDto(&userEntity))
}

//	todo 목록 전체 조회
func GetAllUsers(w http.ResponseWriter, r *http.Request) {

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

//	todo 신규 생성
func CreateUsers(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var userParam users
	err := decoder.Decode(&userParam)
	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusBadRequest, "Couldn't parse request data.")
	}

	// fmt.Println("params: ", todoParam)

	db := ConnectGormDB()

	userEntity := convertUserDtoToEntity(&userParam)

	db.Create(&userEntity)

	responseWithJSON(w, http.StatusOK, convertUserEntityToDto(&userEntity))

}

//	todo 업데이트
func UpdateUsers(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var usersParam users
	err := decoder.Decode(&usersParam)
	if err != nil {
		fmt.Println(err)
		responseWithError(w, http.StatusBadRequest, "Couldn't parse request data.")
	}

	if usersParam.ID == "" {
		responseWithError(w, http.StatusBadRequest, "Couldn't parse request data.")
	}

	fmt.Println("Request Param: ", usersParam)

	var userEntity User
	db := ConnectGormDB()

	db.Where("id = ?", usersParam.ID).First(&userEntity)

	if usersParam.Name != "" {
		userEntity.Name = usersParam.Name
	}

	if usersParam.Birth != "" {
		userEntity.Birth = usersParam.Birth
	}

	db.Save(&userEntity)

	responseWithJSON(w, http.StatusOK, convertUserEntityToDto(&userEntity))

}

//	todo 삭제
func DeleteUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	db := ConnectGormDB()

	var userEntity User

	db.Where("id = ?", userID).Delete(&userEntity)

	responseWithJSON(w, http.StatusOK, map[string]string{"id": userID})

}

func convertUserDtoToEntity(userDto *users) User {
	var userEntity User

	userEntity.ID = userDto.ID
	userEntity.Birth = userDto.Birth
	userEntity.Name = userDto.Name
	userEntity.CreatedAt = time.Now()

	return userEntity
}

func convertUserEntityToDto(userEntity *User) users {
	var dtoUser users

	dtoUser.ID = userEntity.ID
	dtoUser.Birth = userEntity.Birth
	dtoUser.Name = userEntity.Name
	dtoUser.CreatedAt = userEntity.CreatedAt.Format(time.RFC3339)

	return dtoUser

}
