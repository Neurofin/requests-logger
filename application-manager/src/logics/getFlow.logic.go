package logics

import (
	"application-manager/src/models"
)

func GetFlow(uid string) (models.FlowModel, error) {
	flow := models.FlowModel{
		Uid: uid,
	}

	output, err := flow.GetFlow()
	if err != nil {
		return flow, err
	}

	flow = output.Data.(models.FlowModel)
	return flow, nil
}
