package clienthelper

import (
	"net/http"
	wr "../models"
)

func UploadList(list *wr.WordsList, res http.ResponseWriter){
	// todo implement save to db
	// if ok
	res.Write([]byte("Successfully uploaded list of word for the user"))	
}
