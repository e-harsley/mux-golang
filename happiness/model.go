package happiness

type (
	DocumentModel interface {
		GetModelName() string
	}
	BaseModel struct {
		ID        string `json:"id" bson:"_id" mgoType:"id"`
		CreatedAt string `json:"created_at" bson:"createdAt" mgoType:"date"`
		UpdatedAt string `json:"updated_at" bson:"updatedAt" mgoType:"date"`
	}
)
