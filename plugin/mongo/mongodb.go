package mongo

import (
	"context"
	"fmt"
	"github.com/Gitforxuyang/eva/config"
	"github.com/Gitforxuyang/eva/server"
	"github.com/Gitforxuyang/eva/util/logger"
	"github.com/Gitforxuyang/eva/util/trace"
	"github.com/Gitforxuyang/eva/util/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type EvaMongo interface {
	Database(name string) *Database
}

type Database struct {
	//clientName string
	mongo  *evaMongo
	dbName string
	db     *mongo.Database
}

type Collection struct {
	//clientName     string
	//dbName         string
	db             *Database
	collectionName string
	collection     *mongo.Collection
}

type evaMongo struct {
	cli   *mongo.Client
	name  string
	trace *trace.Tracer
	log   logger.EvaLogger
}

func (m *evaMongo) Database(name string) *Database {
	db := m.cli.Database(name)
	return &Database{dbName: name, db: db, mongo: m}
}

func (m *Database) Collection(name string) *Collection {
	c := m.db.Collection(name)
	return &Collection{collectionName: name, collection: c, db: m}
}
func StartSpan(ctx context.Context, m *Collection, f string) (context.Context, opentracing.Span, error) {
	err := utils.ContextDie(ctx)
	if err != nil {
		return nil, nil, err
	}
	ctx, span, err := m.db.mongo.trace.StartMongoClientSpanFromContext(ctx, fmt.Sprintf("mongo.%s", f))
	if err != nil {
		m.db.mongo.log.Error(ctx, "mongo startspan err", logger.Fields{
			"err": err,
		})
	}
	ext.DBStatement.Set(span, f)
	ext.DBInstance.Set(span, m.db.mongo.name)
	span.SetTag("database", m.db.dbName)
	span.SetTag("collection", m.collectionName)
	return ctx, span, nil
}

func FinishSpan(span opentracing.Span, err error) {
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.String("event", "error"))
		span.LogFields(
			log.Object("evaError", utils.StructToJson(err)),
		)
	}
	span.Finish()
}

func (m *Collection) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (r *mongo.InsertOneResult, err error) {
	ctx, span, e := StartSpan(ctx, m, "InsertOne")
	if e != nil {
		return nil, e
	}
	defer func() {
		FinishSpan(span, err)
	}()
	r, err = m.collection.InsertOne(ctx, document, opts...)
	return r, err
}

func (m *Collection) InsertMany(ctx context.Context, documents []interface{},
	opts ...*options.InsertManyOptions) (r *mongo.InsertManyResult, err error) {
	ctx, span, e := StartSpan(ctx, m, "InsertMany")
	if e != nil {
		return nil, e
	}
	defer func() {
		FinishSpan(span, err)
	}()
	r, err = m.collection.InsertMany(ctx, documents, opts...)
	return r, err
}

func (m *Collection) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (r *mongo.UpdateResult, err error) {
	ctx, span, e := StartSpan(ctx, m, "UpdateOne")
	if e != nil {
		return nil, e
	}
	defer func() {
		FinishSpan(span, err)
	}()
	r, err = m.collection.UpdateOne(ctx, filter, update, opts...)
	return r, err
}

func (m *Collection) UpdateMany(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (r *mongo.UpdateResult, err error) {
	ctx, span, e := StartSpan(ctx, m, "UpdateMany")
	if e != nil {
		return nil, e
	}
	defer func() {
		FinishSpan(span, err)
	}()
	r, err = m.collection.UpdateMany(ctx, filter, update, opts...)
	return r, err
}

func (m *Collection) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (r *mongo.Cursor, err error) {
	ctx, span, e := StartSpan(ctx, m, "Find")
	if e != nil {
		return nil, e
	}
	defer func() {
		FinishSpan(span, err)
	}()
	r, err = m.collection.Find(ctx, filter, opts...)
	return r, err
}
func (m *Collection) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) (r *mongo.SingleResult, err error) {
	ctx, span, e := StartSpan(ctx, m, "FindOne")
	if e != nil {
		return nil, e
	}
	defer func() {
		FinishSpan(span, err)
	}()
	r = m.collection.FindOne(ctx, filter, opts...)
	if r.Err() != nil && r.Err() != mongo.ErrNoDocuments {
		err = r.Err()
	}
	return r, err
}
func (m *Collection) FindOneAndUpdate(ctx context.Context, filter interface{},
	update interface{}, opts ...*options.FindOneAndUpdateOptions) (r *mongo.SingleResult, err error) {
	ctx, span, e := StartSpan(ctx, m, "FindOneAndUpdate")
	if e != nil {
		return nil, e
	}
	defer func() {
		FinishSpan(span, err)
	}()
	r = m.collection.FindOneAndUpdate(ctx, filter, update, opts...)
	if r.Err() != nil && r.Err() != mongo.ErrNoDocuments {
		err = r.Err()
	}
	return r, err
}
func (m *Collection) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (r *mongo.DeleteResult, err error) {
	ctx, span, e := StartSpan(ctx, m, "DeleteOne")
	if e != nil {
		return nil, e
	}
	defer func() {
		FinishSpan(span, err)
	}()
	r, err = m.collection.DeleteOne(ctx, filter, opts...)
	return r, err
}

func (m *Collection) DeleteMany(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (r *mongo.DeleteResult, err error) {
	ctx, span, e := StartSpan(ctx, m, "DeleteMany")
	if e != nil {
		return nil, e
	}
	defer func() {
		FinishSpan(span, err)
	}()
	r, err = m.collection.DeleteMany(ctx, filter, opts...)
	return r, err
}

func (m *Collection) Aggregate(ctx context.Context, pipeline interface{},
	opts ...*options.AggregateOptions) (r *mongo.Cursor, err error) {
	ctx, span, e := StartSpan(ctx, m, "Aggregate")
	if e != nil {
		return nil, e
	}
	defer func() {
		FinishSpan(span, err)
	}()
	r, err = m.collection.Aggregate(ctx, pipeline, opts...)
	return r, err
}
func (m *Collection) CountDocuments(ctx context.Context, filter interface{},
	opts ...*options.CountOptions) (r int64, err error) {
	ctx, span, e := StartSpan(ctx, m, "CountDocuments")
	if e != nil {
		return 0, e
	}
	defer func() {
		FinishSpan(span, err)
	}()
	r, err = m.collection.CountDocuments(ctx, filter, opts...)
	return r, err
}

var (
	lock   sync.Mutex
	client map[string]*evaMongo = make(map[string]*evaMongo)
)

func GetMongoClient(name string) EvaMongo {
	lock.Lock()
	defer lock.Unlock()
	if client[name] == nil {
		conf := config.GetConfig().GetMongo(name)
		opts := options.Client().ApplyURI(conf.Url).
			SetMaxPoolSize(conf.MaxPoolSize).
			SetMinPoolSize(conf.MinPoolSize)
		cli, err := mongo.NewClient(opts)
		utils.Must(err)
		ctx, _ := context.WithTimeout(context.TODO(), time.Second*5)
		cli.Connect(ctx)
		server.RegisterShutdownFunc(func() {
			cli.Disconnect(context.TODO())
		})
		client[name] = &evaMongo{
			cli:   cli,
			name:  name,
			trace: trace.GetTracer(),
			log:   logger.GetLogger(),
		}
	}
	return client[name]
}
