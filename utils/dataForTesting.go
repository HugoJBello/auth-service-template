package utils

import (
	"auth-service-template/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var usersExamples = []models.User{}


func init(){
	usersExamples, _ = GetUsersForTesting()
	for _, user := range(usersExamples){
		user.CreateOrUpdateWithPlainPw()
	}
}

func GetUsersForTesting() (users []models.User, err error) {
	if len(usersExamples) == 0 {
		root := "../utils/dataForTesting/user"
		files, _ := ioutil.ReadDir(root)

		for _, f := range files {
			user, err := readUsersJson(root + "/" + f.Name())
			if err != nil {
				fmt.Println(err)
				return usersExamples, err
			}
			usersExamples = append(usersExamples, user)
		}
	}

	return usersExamples, nil
}

func readUsersJson(path string) (user models.User, err error) {
	user = models.User{}
	file, err1 := ioutil.ReadFile(path)
	if err1 != nil {
		fmt.Println(err1)
		return user, err1
	}
	err2 := json.Unmarshal([]byte(file), &user)
	if err2 != nil {
		fmt.Println(err2)
		return user, err2
	}
	return user, nil

}
