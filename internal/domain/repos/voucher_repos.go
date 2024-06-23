package repos

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"latipe-promotion-services/internal/domain/entities"
	"latipe-promotion-services/internal/domain/message"
	"latipe-promotion-services/pkgs/mongodb"
	"latipe-promotion-services/pkgs/pagable"
	"strings"
	"time"
)

type VoucherRepository struct {
	voucherCollection     *mongo.Collection
	voucherLogsCollection *mongo.Collection
}

func NewVoucherRepos(client *mongodb.MongoClient) *VoucherRepository {
	voucherCol := client.GetDB().Collection("vouchers")
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

	model := mongo.IndexModel{Keys: bson.D{{"voucher_code", "text"}}}
	name, err := voucherCol.Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		panic(err)
	}
	fmt.Println("Name of index created: " + name)

	logCol := client.GetDB().Collection("voucher_using_logs")

	log.Info("voucher code unique index created successfully")
	return &VoucherRepository{voucherCollection: voucherCol, voucherLogsCollection: logCol}
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

	err := dr.voucherCollection.FindOne(ctx, bson.M{"voucher_code": strings.ToUpper(voucherCode)}).Decode(&entity)
	if err != nil {
		return nil, err
	}

	return &entity, err
}

func (dr VoucherRepository) GetAll(ctx context.Context, voucherCode string, query *pagable.Query) ([]entities.Voucher, int, error) {
	var delis []entities.Voucher

	filter, err := query.ConvertQueryToFilter()
	if err != nil {
		return nil, 0, err
	}

	if voucherCode != "" {
		filter["$text"] = bson.M{"$search": voucherCode}
	}

	opts := options.Find().SetLimit(int64(query.GetSize())).SetSkip(int64(query.GetOffset()))
	cursor, err := dr.voucherCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	if err = cursor.All(ctx, &delis); err != nil {
		return nil, 0, err
	}

	total, err := dr.voucherCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return delis, int(total), err
}

func (dr VoucherRepository) GetComingVoucher(ctx context.Context, query *pagable.Query) ([]entities.Voucher, int, error) {
	var delis []entities.Voucher

	tomorrow := time.Now().Add(24 * time.Hour)
	tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())

	filter := bson.M{"stated_time": bson.M{"$gte": tomorrow}, "status": entities.ACTIVE}

	opts := options.Find().SetLimit(int64(query.GetSize())).SetSkip(int64(query.GetOffset()))
	cursor, err := dr.voucherCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	if err = cursor.All(ctx, &delis); err != nil {
		return nil, 0, err
	}

	total, err := dr.voucherCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return delis, int(total), err
}

func (dr VoucherRepository) GetVoucherForUser(ctx context.Context, voucherCode string, query *pagable.Query) ([]entities.Voucher, int, error) {
	var delis []entities.Voucher

	opts := options.Find().SetLimit(int64(query.GetSize())).SetSkip(int64(query.GetOffset()))

	filter, err := query.ConvertQueryToFilter()
	if err != nil {
		return nil, 0, err
	}

	filter["stated_time"] = bson.M{"$lt": time.Now()}
	filter["ended_time"] = bson.M{"$gte": time.Now()}
	filter["status"] = entities.ACTIVE // Thêm điều kiện status = 1

	if voucherCode != "" {
		filter["$text"] = bson.M{"$search": voucherCode}
	}

	cursor, err := dr.voucherCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	if err = cursor.All(ctx, &delis); err != nil {
		return nil, 0, err
	}

	total, err := dr.voucherCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return delis, int(total), err
}

func (dr VoucherRepository) Total(ctx context.Context, query *pagable.Query) (int64, error) {
	opts := options.Count().SetHint("_id_")
	filter, err := query.ConvertQueryToFilter()
	if err != nil {
		return 0, err
	}

	count, err := dr.voucherCollection.CountDocuments(ctx, filter, opts)
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

func (dr VoucherRepository) CreateUsingVoucherLog(ctx context.Context, log *entities.VoucherUsingLog) error {
	log.CreatedAt = time.Now()

	_, err := dr.voucherLogsCollection.InsertOne(ctx, log)
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

func (dr VoucherRepository) UpdateVoucherCounts(ctx context.Context, vouchers *entities.Voucher) error {

	update := bson.D{
		{"$set", bson.D{
			{"voucher_counts", vouchers.VoucherCounts},
		}},
	}
	data, err := dr.voucherCollection.UpdateByID(ctx, vouchers.ID, update)
	if err != nil {
		return err
	}

	if data.ModifiedCount == 0 {
		return errors.New("not change")
	}

	return nil
}

func (dr VoucherRepository) UpdateUsingVoucherLog(ctx context.Context, usingLog *entities.VoucherUsingLog) error {

	update := bson.D{
		{"$set", bson.D{
			{"checkout_purchase.order_ids", usingLog.CheckoutPurchase.OrderIDs},
			{"updated_at", time.Now()},
		}},
	}
	data, err := dr.voucherLogsCollection.UpdateByID(ctx, usingLog.ID, update)
	if err != nil {
		return err
	}

	if data.ModifiedCount == 0 {
		return errors.New("not change")
	}

	return nil
}

func (dr VoucherRepository) UpdateUsingStatus(ctx context.Context, voucherLog *entities.VoucherUsingLog) error {

	update := bson.D{
		{"$set", bson.D{
			{"status", voucherLog.Status},
		}},
	}
	data, err := dr.voucherLogsCollection.UpdateByID(ctx, voucherLog.ID, update)
	if err != nil {
		return err
	}

	if data.ModifiedCount == 0 {
		return errors.New("not change")
	}

	return nil
}

func (dr VoucherRepository) GetVoucherOfStore(ctx context.Context, storeId string,
	voucherCode string, query *pagable.Query) ([]entities.Voucher, int, error) {
	var delis []entities.Voucher

	filter, err := query.ConvertQueryToFilter()
	if err != nil {
		return nil, 0, err
	}

	filter["owner_voucher"] = storeId

	if voucherCode != "" {
		filter["$text"] = bson.M{"$search": voucherCode}
	}

	opts := options.Find().SetLimit(int64(query.GetSize())).SetSkip(int64(query.GetOffset()))
	cursor, err := dr.voucherCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	if err = cursor.All(ctx, &delis); err != nil {
		return nil, 0, err
	}

	total, err := dr.voucherCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return delis, int(total), err
}

func (dr VoucherRepository) CheckUsableVoucherOfUser(ctx context.Context, userId string, voucherCode string) (int, error) {
	filter := bson.M{"user_id": userId, "voucher_code": voucherCode, "status": message.COMMIT_SUCCESS}
	count, err := dr.voucherLogsCollection.CountDocuments(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	return int(count), nil
}

func (dr VoucherRepository) FindVoucherLogByOrderID(ctx context.Context, orderId string) ([]*entities.VoucherUsingLog, error) {
	var results []*entities.VoucherUsingLog
	filter := bson.M{"checkout_purchase.order_ids": orderId}
	cursor, err := dr.voucherLogsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, err
}

func (dr VoucherRepository) FindVoucherLogByVoucherCodeAndCheckoutID(ctx context.Context, voucherCode string, checkoutId string) (*entities.VoucherUsingLog, error) {
	var results entities.VoucherUsingLog
	filter := bson.M{
		"checkout_purchase.checkout_id": checkoutId,
		"voucher_code":                  voucherCode,
	}

	err := dr.voucherLogsCollection.FindOne(ctx, filter).Decode(&results)
	if err != nil {
		return nil, err
	}

	return &results, err
}

// CountAllVoucherCreatedInCurrentMonthByStoreId counts all vouchers created in the current month by store ID.
func (dr VoucherRepository) CountAllVoucherCreatedInCurrentMonthByStoreId(ctx context.Context, storeId string) (int, error) {
	now := time.Now()
	year, month, _ := now.Date()
	location := now.Location()

	// Define the start and end of the current month
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, location)
	endDate := startDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// Create filter for vouchers created in the current month and by the specified store ID
	filter := bson.M{
		"created_at": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
		"owner_voucher": storeId,
	}

	// Count the number of vouchers that match the filter
	voucherCount, err := dr.voucherCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(voucherCount), nil
}
