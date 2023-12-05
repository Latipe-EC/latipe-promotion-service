package repos

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"latipe-promotion-services/domain/entities"
	"latipe-promotion-services/pkgs/pagable"
	"time"
)

type VoucherRepository struct {
	voucherCollection     *mongo.Collection
	voucherLogsCollection *mongo.Collection
}

func NewVoucherRepos(db *mongo.Database) VoucherRepository {
	voucherCol := db.Collection("vouchers")
	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"voucher_code": 1,
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := voucherCol.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		panic("error creating unique index:" + err.Error())

	}

	logCol := db.Collection("voucher_using_logs")

	log.Info("voucher code unique index created successfully")
	return VoucherRepository{voucherCollection: voucherCol, voucherLogsCollection: logCol}
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

	filter, err := query.ConvertQueryToFilter()
	if err != nil {
		return nil, err
	}

	opts := options.Find().SetLimit(int64(query.GetSize())).SetSkip(int64(query.GetPage() - 1))
	cursor, err := dr.voucherCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &delis); err != nil {
		return nil, err
	}
	return delis, err
}

func (dr VoucherRepository) GetVoucherForUser(ctx context.Context, query *pagable.Query) ([]entities.Voucher, error) {
	var delis []entities.Voucher

	opts := options.Find().SetLimit(int64(query.GetSize())).SetSkip(int64(query.GetPage() - 1))

	filter := bson.M{
		"stated_time": bson.M{"$lt": time.Now()},
		"ended_time":  bson.M{"$gt": time.Now()},
		"status":      entities.ACTIVE, // Thêm điều kiện status = 1
	}
	cursor, err := dr.voucherCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &delis); err != nil {
		return nil, err
	}
	return delis, err
}

func (dr VoucherRepository) Total(ctx context.Context, query *pagable.Query) (int64, error) {
	opts := options.Count().SetHint("_id_")
	filter, err := query.ConvertQueryToFilter()
	if err != nil {
		return 0, err
	}

	count, err := dr.voucherCollection.CountDocuments(context.TODO(), filter, opts)
	if err != nil {
		return -1, err
	}

	return count, err
}

func (dr VoucherRepository) CreateVoucher(ctx context.Context, voucher *entities.Voucher) (string, error) {
	voucher.CreatedAt = time.Now()
	voucher.UpdatedAt = time.Now()

	lastId, err := dr.voucherCollection.InsertOne(ctx, voucher)
	if err != nil {
		return "", err
	}
	return lastId.InsertedID.(primitive.ObjectID).Hex(), err
}

func (dr VoucherRepository) CreateLogUseVoucher(ctx context.Context, voucher *entities.VoucherUsingLog) error {
	voucher.CreatedAt = time.Now()

	_, err := dr.voucherLogsCollection.InsertOne(ctx, voucher)
	if err != nil {
		return err
	}
	return err
}

func (dr VoucherRepository) UpdateStatus(ctx context.Context, voucher *entities.Voucher) error {

	update := bson.D{
		{"$set", bson.D{
			{"status", voucher.Status},
			{"updated_at", time.Now()},
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

func (dr VoucherRepository) UpdateVoucherCounts(ctx context.Context, vouchers []*entities.Voucher) error {
	for _, i := range vouchers {
		update := bson.D{
			{"$set", bson.D{
				{"voucher_counts", i.VoucherCounts},
			}},
		}
		data, err := dr.voucherCollection.UpdateByID(ctx, i.ID, update)
		if err != nil {
			return err
		}
		if data.ModifiedCount == 0 {
			return errors.New("not change")
		}
	}

	return nil
}
