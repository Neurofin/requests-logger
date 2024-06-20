package logics

import "go.mongodb.org/mongo-driver/bson/primitive"

func ProcessApplicationChecklistLogic(appId primitive.ObjectID) {
	// TODO: Get Application docs, see if all the documents are uploaded
	// TODO: Get all checklist items belonging to the documents and send them to query manager
}
