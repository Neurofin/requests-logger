package logDb

import (
	logTypeEnum "logger/src/store/enum"
	loggerTypes "logger/src/store/types"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	Id        primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Type      logTypeEnum.LogType    `json:"type" bson:"type"`
	Data      map[string]interface{} `json:"data" bson:"data"`
	TraceId   string                 `json:"traceId" bson:"traceId"`
	Timestamp time.Time              `json:"timestamp" bson:"timestamp"`
	loggerTypes.Timestamps
}
