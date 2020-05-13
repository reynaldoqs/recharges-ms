package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"audiman/projects/recharger-x/entity"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

// NewMongoRepository creates a mongo repository
func NewMongoRepository(mongoURL, mongoDB string, mongoTimeout int) (RechargeRepository, error) {
	repo := &mongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}
	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		err := errors.Wrap(err, "repository.rechargesMongoRepo.NewMongoRepository")
		return nil, err
	}

	repo.client = client

	return repo, nil
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		err := errors.Wrap(err, "repository.rechargesMongoRepo.newMongoClient")
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		err := errors.Wrap(err, "repository.rechargesMongoRepo.newMongoClient")
		return nil, err
	}
	return client, nil
}

func (r *mongoRepository) Save(recharge *entity.Recharge) (*entity.Recharge, error) {

	if r.client == nil {
		fmt.Println("es nil carajo")
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.database).Collection("recharges")

	result, err := collection.InsertOne(
		ctx,
		bson.M{
			"phoneNumber": recharge.PhoneNumber,
			"company":     recharge.Company,
			"cardNumber":  recharge.CardNumber,
			"status":      recharge.Status,
			"mount":       recharge.Mount,
			"idResolver":  recharge.IDResolver,
			"createdAt":   recharge.CreatedAt,
			"resolvedAt":  recharge.ResolvedAt,
		},
	)

	if err != nil {
		err := errors.Wrap(err, "repository.rechargesMongoRepo.Save")
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		recharge.ID = oid.Hex()
	}
	return recharge, nil
}

/*
func (r *mongoRepository) Update(recharge *entity.Recharge) error {

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.database).Collection("recharges")

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id": recharge.ID}}
	update := bson.D{{"$set",
		bson.D{{
			"phoneNumber": recharge.PhoneNumber,
			"company":     recharge.Company,
			"cardNumber":  recharge.CardNumber,
			"status":      recharge.Status,
			"idResolver":  recharge.IDResolver,
			"createdAt":   recharge.CreatedAt,
			"resolvedAt":  recharge.ResolvedAt,
		}},
	}}

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		err := errors.Wrap(err, "repository.rechargesMongoRepo.Update")
		return err
	}

	return nil
}
*/
func (r *mongoRepository) UpdateRecharger(rmsg *entity.RechargeMsg) error {

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.database).Collection("recharges")

	id, err := primitive.ObjectIDFromHex(string(rmsg.IDRecharge))
	if err != nil {
		err := errors.Wrap(err, "repository.rechargesMongoRepo.UpdateRecharger")
		return err
	}

	// FIX: this date needs to come from recharger microservice
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{"$set", bson.M{"status": 3, "idResolver": rmsg.IDResolver, "resolvedAt": rmsg.ResolvedAt}}}

	_, err = collection.UpdateOne(ctx, filter, update, nil)
	if err != nil {
		err := errors.Wrap(err, "repository.rechargesMongoRepo.UpdateRecharger")
		return err
	}

	return nil
}

func (r *mongoRepository) FindAll() ([]*entity.Recharge, error) {

	collection := r.client.Database(r.database).Collection("recharges")
	findOptions := options.Find()

	var results []*entity.Recharge

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		err := errors.Wrap(err, "repository.rechargesMongoRepo.FindAll")
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var elem entity.Recharge
		err := cur.Decode(&elem)
		if err != nil {
			err := errors.Wrap(err, "repository.rechargesMongoRepo.FindAll")
			return nil, err
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		err := errors.Wrap(err, "repository.rechargesMongoRepo.FindAll")
		return nil, err
	}
	cur.Close(context.TODO())
	return results, nil
}
