package utils

import (
	"errors"
)

const (
	userType int = iota
	photoType
	likeType
	commentType
)

func validateDataByID(idType int, id int) error {
	var keys []int
	if idType == userType {
		//This should be taken from db
		users := map[int]User{}
		keys = make([]int, 0, len(users))
		for k := range users {
			keys = append(keys, k)
		}
	} else if idType == photoType {
		//This should be taken from db
		photos := map[int]Photo{}
		keys = make([]int, 0, len(photos))
		for k := range photos {
			keys = append(keys, k)
		}
	}

	if !arrayContains(keys, id) {
		return errors.New("StatusNotFound")
	}
	return nil
}

func arrayContains[T comparable](array []T, val T) bool {
	for _, v := range array {
		// checking if the array contains the given value
		if v == val {
			// changing the boolean variable
			return true
		}
	}
	return false
}

func remove[T interface{}](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}
