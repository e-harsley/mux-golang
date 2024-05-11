package happiness

import (
	"context"
	"encoding/json"
	"github.com/Lyearn/mgod"
	"github.com/Lyearn/mgod/schema/schemaopt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type IRepository[T DocumentModel] interface {
	FindOneAndUpdate(filter, update interface{}, opts ...*options.FindOneAndUpdateOptions) (T, error)
	GetDocToInsert(model T) (bson.D, error)
	BindDataOperationStruct(data interface{}) (T, error)
	InsertOne(model interface{}, opts ...*options.InsertOneOptions) (T, error)
	UpdateMany(filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	Find(filter interface{}, opts ...*options.FindOptions) ([]T, error)
	FindOne(filter interface{}, opts ...*options.FindOneOptions) (*T, error)
	DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	CountDocuments(filter interface{}, opts ...*options.CountOptions) (int64, error)
	Distinct(fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error)
	Aggregate(pipeline interface{}, opts ...*options.AggregateOptions) ([]bson.D, error)
}

type BaseRepository[T DocumentModel] struct {
	modelService mgod.EntityMongoModel[T]
}

func (cls BaseRepository[T]) BindDataOperationStruct(data interface{}) (T, error) {
	var mod T
	preMod := map[string]interface{}{}
	jsonString, _ := json.Marshal(data)

	if err := json.Unmarshal(jsonString, &preMod); err != nil {
		return mod, err
	}

	preMod["id"] = primitive.NewObjectID()

	jsonString, _ = json.Marshal(preMod)

	if err := json.Unmarshal(jsonString, &mod); err != nil {
		return mod, err
	}
	//fmt.Println("\n here", mod)
	return mod, nil

}

func NewBaseRepository[M DocumentModel](mod M) BaseRepository[M] {
	InitConnect()
	schemaOpts := schemaopt.SchemaOptions{
		Timestamps: true,
		Collection: mod.GetModelName(),
	}

	modelService, err := mgod.NewEntityMongoModel(mod, schemaOpts)

	if err != nil {
		log.Fatal(err)
	}

	return BaseRepository[M]{
		modelService: modelService,
	}
}

func (cls BaseRepository[M]) FindOneAndUpdate(filter, update interface{}, opts ...*options.FindOneAndUpdateOptions) (M, error) {

	updateMap := map[string]interface{}{}

	var model M

	jsonString, _ := json.Marshal(update)
	if err := json.Unmarshal(jsonString, &updateMap); err != nil {
		return model, err
	}

	upd := bson.D{
		{"$set", updateMap},
	}

	return cls.modelService.FindOneAndUpdate(context.TODO(), filter, upd, opts...)

}

func (cls BaseRepository[T]) GetDocToInsert(model T) (bson.D, error) {
	return cls.modelService.GetDocToInsert(context.TODO(), model)
}

func (cls BaseRepository[T]) InsertOne(model interface{}, opts ...*options.InsertOneOptions) (T, error) {
	return cls.modelService.InsertOne(context.TODO(), model, opts...)
}

func (cls BaseRepository[T]) UpdateMany(filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return cls.modelService.UpdateMany(context.TODO(), filter, update, opts...)
}

func (cls BaseRepository[T]) BulkWrite(bulkWrites []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	return cls.modelService.BulkWrite(context.TODO(), bulkWrites, opts...)
}

func (cls BaseRepository[T]) Find(filter interface{}, opts ...*options.FindOptions) ([]T, error) {
	return cls.modelService.Find(context.TODO(), filter, opts...)
}

func (cls BaseRepository[T]) FindOne(filter interface{}, opts ...*options.FindOneOptions) (*T, error) {
	return cls.modelService.FindOne(context.TODO(), filter, opts...)
}

func (cls BaseRepository[T]) DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return cls.modelService.DeleteOne(context.TODO(), filter, opts...)
}

func (cls BaseRepository[T]) DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return cls.modelService.DeleteMany(context.TODO(), filter, opts...)
}

func (cls BaseRepository[T]) CountDocuments(filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return cls.modelService.CountDocuments(context.TODO(), filter, opts...)
}

func (cls BaseRepository[T]) Distinct(fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	return cls.modelService.Distinct(context.TODO(), fieldName, filter, opts...)
}

func (cls BaseRepository[T]) Aggregate(pipeline interface{}, opts ...*options.AggregateOptions) ([]bson.D, error) {
	return cls.modelService.Aggregate(context.TODO(), pipeline, opts...)
}
