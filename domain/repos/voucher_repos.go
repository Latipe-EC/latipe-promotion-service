package repos

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"latipe-promotion-services/domain/entities"
	"latipe-promotion-services/pkgs/pagable"
	"time"
)

type VoucherRepository struct {
	voucherCollection *mongo.Collection
}

func NewVoucherRepos(db *mongo.Database) VoucherRepository {
	col := db.Collection("vouchers")
	return VoucherRepository{voucherCollection: col}
}

func (dr VoucherRepository) GetById(ctx context.Context, Id string) (*entities.Voucher, error) {
	var entity entities.Voucher
	id, _ := primitive.ObjectIDFromHex(Id)

	err := dr.voucherCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&entity)
	if err != nil {
		return nil, err
	}

	return &entity, err
}

func (dr VoucherRepository) GetByCode(ctx context.Context, voucherCode string) (*entities.Voucher, error) {
	var entity entities.Voucher

	err := dr.voucherCollection.FindOne(ctx, bson.M{"voucher_code": voucherCode}).Decode(&entity)
	if err != nil {
		return nil, err
	}

	return &entity, err
}
func (dr VoucherRepository) GetAll(ctx context.Context, query *pagable.Query) ([]entities.Voucher, error) {
	var delis []entities.Voucher

	opts := options.Find().SetLimit(int64(query.GetSize())).SetSkip(int64(query.GetPage()))
	cursor, err := dr.voucherCollection.Find(ctx, bson.D{{}}, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &delis); err != nil {
		return nil, err
	}
	return delis, err
}

func (dr VoucherRepository) Total(ctx context.Context, query *pagable.Query) (int64, error) {

	count, err := dr.voucherCollection.CountDocuments(context.TODO(), bson.D{}, nil)
	if err != nil {
		return -1, err
	}

	return count, err
}

func (dr VoucherRepository) CreateVoucher(ctx context.Context, voucher *entities.Voucher) (string, error) {
	voucher.CreateAt = time.Now()
	voucher.UpdateAt = time.Now()

	lastId, err := dr.voucherCollection.InsertOne(ctx, voucher)
	if err != nil {
		return "", err
	}
	return lastId.InsertedID.(primitive.ObjectID).Hex(), err
}

func (dr VoucherRepository) UpdateStatus(ctx context.Context, voucher *entities.Voucher) error {

	update := bson.D{
		{"$set", bson.D{
			{"status", voucher.Status},
			{"update_at", time.Now()},
		}},
	}
	data, err := dr.voucherCollection.UpdateByID(ctx, voucher.ID, update)
	if err != nil {
		return err
	}

	if data.ModifiedCount == 0 {
		return errors.New("not change")
	}

	return nil
}
