package happiness

import (
	"encoding/json"
	"fmt"
)

var (
	requestErrorWraper = func(err error) any {
		return &StandardResponse{
			Data:    err.Error(),
			Message: "Error",
		}
	}
)

type IRequest interface {
	Validate() error
}

type BaseRequest struct {
}

func (br BaseRequest) Validate() error {
	fmt.Println("checkkk")
	return nil
}

type ExtraParameters struct {
	UserID string `json:"user_id"`
}

func (er ExtraParameters) Validate() error {
	fmt.Println("checkkk")
	return nil
}

func BindDataOperationStruct(data interface{}, output interface{}) error {
	jsonString, _ := json.Marshal(data)

	err := json.Unmarshal([]byte(jsonString), &output)

	if err != nil {

		return err
	}

	fmt.Println(output)

	return nil

}
