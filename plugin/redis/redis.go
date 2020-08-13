package redis

import (
	"context"
	"fmt"
	"github.com/Gitforxuyang/eva/config"
	"github.com/Gitforxuyang/eva/server"
	error2 "github.com/Gitforxuyang/eva/util/error"
	"github.com/Gitforxuyang/eva/util/logger"
	"github.com/Gitforxuyang/eva/util/trace"
	"github.com/Gitforxuyang/eva/util/utils"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"sync"
	"time"
)

type EvaRedis interface {
	Ping(ctx context.Context) (string, error)
	Del(ctx context.Context, keys ...string) (int64, error)
	Dump(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error)
	Persist(ctx context.Context, key string) (bool, error)
	PExpire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	PExpireAt(ctx context.Context, key string, tm time.Time) (bool, error)
	PTTL(ctx context.Context, key string) (time.Duration, error)
	Sort(ctx context.Context, key string, sort *redis.Sort) ([]string, error)
	SortStore(ctx context.Context, key, store string, sort *redis.Sort) (int64, error)
	SortInterfaces(ctx context.Context, key string, sort *redis.Sort) ([]interface{}, error)
	TTL(ctx context.Context, key string) (time.Duration, error)
	Decr(ctx context.Context, key string) (int64, error)
	DecrBy(ctx context.Context, key string, decrement int64) (int64, error)
	Get(ctx context.Context, key string) (string, error)
	GetRange(ctx context.Context, key string, start, end int64) (string, error)
	GetSet(ctx context.Context, key string, value interface{}) (string, error)
	Incr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error)
	IncrByFloat(ctx context.Context, key string, value float64) (float64, error)
	MGet(ctx context.Context, keys ...string) ([]interface{}, error)
	MSet(ctx context.Context, values ...interface{}) (string, error)
	MSetNX(ctx context.Context, values ...interface{}) (bool, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	SetRange(ctx context.Context, key string, offset int64, value string) (int64, error)
	StrLen(ctx context.Context, key string) (int64, error)
	HDel(ctx context.Context, key string, fields ...string) (int64, error)
	HExists(ctx context.Context, key, field string) (bool, error)
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error)
	HIncrByFloat(ctx context.Context, key, field string, incr float64) (float64, error)
	HLen(ctx context.Context, key string) (int64, error)
	HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error)
	HSet(ctx context.Context, key string, values ...interface{}) (int64, error)
	HMSet(ctx context.Context, key string, values ...interface{}) (bool, error)
	HSetNX(ctx context.Context, key, field string, value interface{}) (bool, error)
	HVals(ctx context.Context, key string) ([]string, error)
	BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error)
	BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error)
	BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (string, error)
	LIndex(ctx context.Context, key string, index int64) (string, error)
	LInsert(ctx context.Context, key, op string, pivot, value interface{}) (int64, error)
	LInsertBefore(ctx context.Context, key string, pivot, value interface{}) (int64, error)
	LInsertAfter(ctx context.Context, key string, pivot, value interface{}) (int64, error)
	LLen(ctx context.Context, key string) (int64, error)
	LPop(ctx context.Context, key string) (string, error)
	LPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	LPushX(ctx context.Context, key string, values ...interface{}) (int64, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error)
	LSet(ctx context.Context, key string, index int64, value interface{}) (string, error)
	LTrim(ctx context.Context, key string, start, stop int64) (string, error)
	RPop(ctx context.Context, key string) (string, error)
	RPopLPush(ctx context.Context, source, destination string) (string, error)
	RPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	RPushX(ctx context.Context, key string, values ...interface{}) (int64, error)
	SAdd(ctx context.Context, key string, members ...interface{}) (int64, error)
	SCard(ctx context.Context, key string) (int64, error)
	SDiff(ctx context.Context, keys ...string) ([]string, error)
	SDiffStore(ctx context.Context, destination string, keys ...string) (int64, error)
	SInter(ctx context.Context, keys ...string) ([]string, error)
	SInterStore(ctx context.Context, destination string, keys ...string) (int64, error)
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)
	SMembers(ctx context.Context, key string) ([]string, error)
	SMembersMap(ctx context.Context, key string) (map[string]struct{}, error)
	SMove(ctx context.Context, source, destination string, member interface{}) (bool, error)
	SPop(ctx context.Context, key string) (string, error)
	SPopN(ctx context.Context, key string, count int64) ([]string, error)
	SRandMember(ctx context.Context, key string) (string, error)
	SRandMemberN(ctx context.Context, key string, count int64) ([]string, error)
	SRem(ctx context.Context, key string, members ...interface{}) (int64, error)
	SUnion(ctx context.Context, keys ...string) ([]string, error)
	SUnionStore(ctx context.Context, destination string, keys ...string) (int64, error)
	ZAdd(ctx context.Context, key string, members ...*redis.Z) (int64, error)
	ZAddNX(ctx context.Context, key string, members ...*redis.Z) (int64, error)
	ZAddXX(ctx context.Context, key string, members ...*redis.Z) (int64, error)
	ZAddCh(ctx context.Context, key string, members ...*redis.Z) (int64, error)
	ZAddNXCh(ctx context.Context, key string, members ...*redis.Z) (int64, error)
	ZAddXXCh(ctx context.Context, key string, members ...*redis.Z) (int64, error)
	ZIncr(ctx context.Context, key string, member *redis.Z) (float64, error)
	ZIncrNX(ctx context.Context, key string, member *redis.Z) (float64, error)
	ZIncrXX(ctx context.Context, key string, member *redis.Z) (float64, error)
	ZCard(ctx context.Context, key string) (int64, error)
	ZCount(ctx context.Context, key, min, max string) (int64, error)
	ZLexCount(ctx context.Context, key, min, max string) (int64, error)
	ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error)
	ZInterStore(ctx context.Context, destination string, store *redis.ZStore) (int64, error)
	ZPopMax(ctx context.Context, key string, count ...int64) ([]redis.Z, error)
	ZPopMin(ctx context.Context, key string, count ...int64) ([]redis.Z, error)
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error)
	ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error)
	ZRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error)
	ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) ([]redis.Z, error)
	ZRank(ctx context.Context, key, member string) (int64, error)
	ZRem(ctx context.Context, key string, members ...interface{}) (int64, error)
	ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error)
	ZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error)
	ZRemRangeByLex(ctx context.Context, key, min, max string) (int64, error)
	ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error)
	ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error)
	ZRevRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error)
	ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) ([]redis.Z, error)
	ZRevRank(ctx context.Context, key, member string) (int64, error)
	ZScore(ctx context.Context, key, member string) (float64, error)
	ZUnionStore(ctx context.Context, dest string, store *redis.ZStore) (int64, error)
	Lock(ctx context.Context, key string, ttl time.Duration, opt *LockOption) (*Lock, error)
}

func LockOptions(retryDuration time.Duration, max int) *LockOption {
	return &LockOption{retryDuration, max}
}

type LockOption struct {
	retryDuration time.Duration //每次重试的间隔
	maxRetry      int           //最大重试次数
}
type evaRedis struct {
	client *redis.Client
	name   string
	addr   string
	db     int
	tmp    []byte
	tmpMu  sync.Mutex
}

var (
	lock   sync.Mutex
	client map[string]*evaRedis = make(map[string]*evaRedis)
)

func GetRedisClient(name string) EvaRedis {
	lock.Lock()
	defer lock.Unlock()
	if client[name] == nil {
		conf := config.GetConfig()
		c := conf.GetRedis(name)
		rdb := redis.NewClient(&redis.Options{
			Addr:         c.Addr,
			Password:     c.Password, // no password set
			DB:           c.DB,       // use default DB
			PoolSize:     c.PoolSize,
			MinIdleConns: c.MinIdleConns,
			DialTimeout:  time.Second * time.Duration(c.DialTimeout),
			ReadTimeout:  time.Second * time.Duration(c.ReadTimeout),
			WriteTimeout: time.Second * time.Duration(c.WriteTimeout),
		})
		rdb.AddHook(&redisHook{tracer: trace.GetTracer(), log: logger.GetLogger(),
			name: name, addr: c.Addr, db: c.DB, appName: conf.GetName(), traceConfig: conf.GetTraceConfig()})
		client[name] = &evaRedis{client: rdb, name: name, addr: c.Addr, db: c.DB}
		//注册关闭服务时的优雅关闭
		server.RegisterShutdownFunc(func() {
			rdb.Close()
		})
	}
	return client[name]
}

type redisHook struct {
	tracer      *trace.Tracer
	log         logger.EvaLogger
	name        string
	addr        string
	db          int
	appName     string
	traceConfig *config.TraceConfig
}

func (m *redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	err := utils.ContextDie(ctx)
	if err != nil {
		return ctx, err
	}
	if m.traceConfig.Redis {
		ctx, span, err := m.tracer.StartRedisClientSpanFromContext(ctx, fmt.Sprintf("redis.%s", cmd.Name()))
		if err != nil {
			m.log.Error(ctx, "链路错误", logger.Fields{"err": utils.StructToMap(err)})
		}
		ext.DBStatement.Set(span, cmd.String())
		ext.DBInstance.Set(span, m.name)
		ext.PeerAddress.Set(span, m.addr)
		span.SetTag("DB", m.db)
		ctx = context.WithValue(ctx, "span", &span)
	}
	return ctx, nil
}

func (m *redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) (err error) {
	if m.traceConfig.Redis {
		s := ctx.Value("span")
		sp, _ := s.(*opentracing.Span)
		span := *sp
		defer span.Finish()
		defer func() {
			if err != nil {
				ext.Error.Set(span, true)
				span.LogFields(log.String("event", "error"))
				span.LogFields(
					log.Object("evaError", utils.StructToJson(err)),
				)
			}
		}()
		err = cmd.Err()
		if err != nil {
			err = error2.RedisError.SetAppId(m.appName).SetDetail(err.Error())
		}
	}
	return err
}

func (m *redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	fmt.Println("BeforeProcessPipeline")
	return ctx, nil
}

func (m *redisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	fmt.Println("AfterProcessPipeline")
	return nil
}

func (m *evaRedis) Ping(ctx context.Context) (string, error) {
	return m.client.Ping(ctx).Result()
}

func (m *evaRedis) Del(ctx context.Context, keys ...string) (int64, error) {
	return m.client.Del(ctx, keys...).Result()
}

func (m *evaRedis) Dump(ctx context.Context, key string) (string, error) {
	return m.client.Dump(ctx, key).Result()
}

func (m *evaRedis) Exists(ctx context.Context, keys ...string) (int64, error) {
	return m.client.Exists(ctx, keys...).Result()
}

func (m *evaRedis) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return m.client.Expire(ctx, key, expiration).Result()
}

func (m *evaRedis) ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error) {
	return m.client.ExpireAt(ctx, key, tm).Result()
}

func (m *evaRedis) Persist(ctx context.Context, key string) (bool, error) {
	return m.client.Persist(ctx, key).Result()
}

func (m *evaRedis) PExpire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return m.client.PExpire(ctx, key, expiration).Result()
}

func (m *evaRedis) PExpireAt(ctx context.Context, key string, tm time.Time) (bool, error) {
	return m.client.PExpireAt(ctx, key, tm).Result()
}

func (m *evaRedis) PTTL(ctx context.Context, key string) (time.Duration, error) {
	return m.client.PTTL(ctx, key).Result()
}

func (m *evaRedis) Sort(ctx context.Context, key string, sort *redis.Sort) ([]string, error) {
	return m.client.Sort(ctx, key, sort).Result()
}

func (m *evaRedis) SortStore(ctx context.Context, key, store string, sort *redis.Sort) (int64, error) {
	return m.client.SortStore(ctx, key, store, sort).Result()
}

func (m *evaRedis) SortInterfaces(ctx context.Context, key string, sort *redis.Sort) ([]interface{}, error) {
	return m.client.SortInterfaces(ctx, key, sort).Result()
}

func (m *evaRedis) TTL(ctx context.Context, key string) (time.Duration, error) {
	return m.client.TTL(ctx, key).Result()
}

func (m *evaRedis) Decr(ctx context.Context, key string) (int64, error) {
	return m.client.Decr(ctx, key).Result()
}

func (m *evaRedis) DecrBy(ctx context.Context, key string, decrement int64) (int64, error) {
	return m.client.DecrBy(ctx, key, decrement).Result()
}

func (m *evaRedis) Get(ctx context.Context, key string) (string, error) {
	return m.client.Get(ctx, key).Result()
}

func (m *evaRedis) GetRange(ctx context.Context, key string, start, end int64) (string, error) {
	return m.client.GetRange(ctx, key, start, end).Result()
}

func (m *evaRedis) GetSet(ctx context.Context, key string, value interface{}) (string, error) {
	return m.client.GetSet(ctx, key, value).Result()
}

func (m *evaRedis) Incr(ctx context.Context, key string) (int64, error) {
	return m.client.Incr(ctx, key).Result()
}

func (m *evaRedis) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return m.client.IncrBy(ctx, key, value).Result()
}

func (m *evaRedis) IncrByFloat(ctx context.Context, key string, value float64) (float64, error) {
	return m.client.IncrByFloat(ctx, key, value).Result()
}

func (m *evaRedis) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return m.client.MGet(ctx, keys...).Result()
}

func (m *evaRedis) MSet(ctx context.Context, values ...interface{}) (string, error) {
	return m.client.MSet(ctx, values...).Result()
}

func (m *evaRedis) MSetNX(ctx context.Context, values ...interface{}) (bool, error) {
	return m.client.MSetNX(ctx, values...).Result()
}

func (m *evaRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return m.client.Set(ctx, key, value, expiration).Result()
}

func (m *evaRedis) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return m.client.SetNX(ctx, key, value, expiration).Result()
}

func (m *evaRedis) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return m.client.SetXX(ctx, key, value, expiration).Result()
}

func (m *evaRedis) SetRange(ctx context.Context, key string, offset int64, value string) (int64, error) {
	return m.client.SetRange(ctx, key, offset, value).Result()
}

func (m *evaRedis) StrLen(ctx context.Context, key string) (int64, error) {
	return m.client.StrLen(ctx, key).Result()
}

func (m *evaRedis) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return m.client.HDel(ctx, key, fields...).Result()
}

func (m *evaRedis) HExists(ctx context.Context, key, field string) (bool, error) {
	return m.client.HExists(ctx, key, field).Result()
}

func (m *evaRedis) HGet(ctx context.Context, key, field string) (string, error) {
	return m.client.HGet(ctx, key, field).Result()
}

func (m *evaRedis) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return m.client.HGetAll(ctx, key).Result()
}

func (m *evaRedis) HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	return m.client.HIncrBy(ctx, key, field, incr).Result()
}

func (m *evaRedis) HIncrByFloat(ctx context.Context, key, field string, incr float64) (float64, error) {
	return m.client.HIncrByFloat(ctx, key, field, incr).Result()
}

func (m *evaRedis) HLen(ctx context.Context, key string) (int64, error) {
	return m.client.HLen(ctx, key).Result()
}

func (m *evaRedis) HMGet(ctx context.Context, key string, fields ...string) ([]interface{}, error) {
	return m.client.HMGet(ctx, key, fields...).Result()
}

func (m *evaRedis) HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return m.client.HSet(ctx, key, values...).Result()
}

func (m *evaRedis) HMSet(ctx context.Context, key string, values ...interface{}) (bool, error) {
	return m.client.HMSet(ctx, key, values...).Result()
}

func (m *evaRedis) HSetNX(ctx context.Context, key, field string, value interface{}) (bool, error) {
	return m.client.HSetNX(ctx, key, field, value).Result()
}

func (m *evaRedis) HVals(ctx context.Context, key string) ([]string, error) {
	return m.client.HVals(ctx, key).Result()
}

func (m *evaRedis) BLPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	return m.client.BLPop(ctx, timeout, keys...).Result()
}

func (m *evaRedis) BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error) {
	return m.client.BRPop(ctx, timeout, keys...).Result()
}

func (m *evaRedis) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) (string, error) {
	return m.client.BRPopLPush(ctx, source, destination, timeout).Result()
}

func (m *evaRedis) LIndex(ctx context.Context, key string, index int64) (string, error) {
	return m.client.LIndex(ctx, key, index).Result()
}

func (m *evaRedis) LInsert(ctx context.Context, key, op string, pivot, value interface{}) (int64, error) {
	return m.client.LInsert(ctx, key, op, pivot, value).Result()
}

func (m *evaRedis) LInsertBefore(ctx context.Context, key string, pivot, value interface{}) (int64, error) {
	return m.client.LInsertBefore(ctx, key, pivot, value).Result()
}

func (m *evaRedis) LInsertAfter(ctx context.Context, key string, pivot, value interface{}) (int64, error) {
	return m.client.LInsertAfter(ctx, key, pivot, value).Result()
}

func (m *evaRedis) LLen(ctx context.Context, key string) (int64, error) {
	return m.client.LLen(ctx, key).Result()
}

func (m *evaRedis) LPop(ctx context.Context, key string) (string, error) {
	return m.client.LPop(ctx, key).Result()
}

func (m *evaRedis) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return m.client.LPush(ctx, key, values...).Result()
}

func (m *evaRedis) LPushX(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return m.client.LPushX(ctx, key, values...).Result()
}

func (m *evaRedis) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return m.client.LRange(ctx, key, start, stop).Result()
}

func (m *evaRedis) LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error) {
	return m.client.LRem(ctx, key, count, value).Result()
}

func (m *evaRedis) LSet(ctx context.Context, key string, index int64, value interface{}) (string, error) {
	return m.client.LSet(ctx, key, index, value).Result()
}

func (m *evaRedis) LTrim(ctx context.Context, key string, start, stop int64) (string, error) {
	return m.client.LTrim(ctx, key, start, stop).Result()
}

func (m *evaRedis) RPop(ctx context.Context, key string) (string, error) {
	return m.client.RPop(ctx, key).Result()
}

func (m *evaRedis) RPopLPush(ctx context.Context, source, destination string) (string, error) {
	return m.client.RPopLPush(ctx, source, destination).Result()
}

func (m *evaRedis) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return m.client.RPush(ctx, key, values...).Result()
}

func (m *evaRedis) RPushX(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return m.client.RPushX(ctx, key, values...).Result()
}

func (m *evaRedis) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return m.client.SAdd(ctx, key, members...).Result()
}

func (m *evaRedis) SCard(ctx context.Context, key string) (int64, error) {
	return m.client.SCard(ctx, key).Result()
}

func (m *evaRedis) SDiff(ctx context.Context, keys ...string) ([]string, error) {
	return m.client.SDiff(ctx, keys...).Result()
}

func (m *evaRedis) SDiffStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	return m.client.SDiffStore(ctx, destination, keys...).Result()
}

func (m *evaRedis) SInter(ctx context.Context, keys ...string) ([]string, error) {
	return m.client.SInter(ctx, keys...).Result()
}

func (m *evaRedis) SInterStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	return m.client.SInterStore(ctx, destination, keys...).Result()
}

func (m *evaRedis) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return m.client.SIsMember(ctx, key, member).Result()
}

func (m *evaRedis) SMembers(ctx context.Context, key string) ([]string, error) {
	return m.client.SMembers(ctx, key).Result()
}

func (m *evaRedis) SMembersMap(ctx context.Context, key string) (map[string]struct{}, error) {
	return m.client.SMembersMap(ctx, key).Result()
}

func (m *evaRedis) SMove(ctx context.Context, source, destination string, member interface{}) (bool, error) {
	return m.client.SMove(ctx, source, destination, member).Result()
}

func (m *evaRedis) SPop(ctx context.Context, key string) (string, error) {
	return m.client.SPop(ctx, key).Result()
}

func (m *evaRedis) SPopN(ctx context.Context, key string, count int64) ([]string, error) {
	return m.client.SPopN(ctx, key, count).Result()
}

func (m *evaRedis) SRandMember(ctx context.Context, key string) (string, error) {
	return m.client.SRandMember(ctx, key).Result()
}

func (m *evaRedis) SRandMemberN(ctx context.Context, key string, count int64) ([]string, error) {
	return m.client.SRandMemberN(ctx, key, count).Result()
}

func (m *evaRedis) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return m.client.SRem(ctx, key, members...).Result()
}

func (m *evaRedis) SUnion(ctx context.Context, keys ...string) ([]string, error) {
	return m.client.SUnion(ctx, keys...).Result()
}

func (m *evaRedis) SUnionStore(ctx context.Context, destination string, keys ...string) (int64, error) {
	return m.client.SUnionStore(ctx, destination, keys...).Result()
}

func (m *evaRedis) ZAdd(ctx context.Context, key string, members ...*redis.Z) (int64, error) {
	return m.client.ZAdd(ctx, key, members...).Result()
}

func (m *evaRedis) ZAddNX(ctx context.Context, key string, members ...*redis.Z) (int64, error) {
	return m.client.ZAddNX(ctx, key, members...).Result()
}

func (m *evaRedis) ZAddXX(ctx context.Context, key string, members ...*redis.Z) (int64, error) {
	return m.client.ZAddXX(ctx, key, members...).Result()
}

func (m *evaRedis) ZAddCh(ctx context.Context, key string, members ...*redis.Z) (int64, error) {
	return m.client.ZAddCh(ctx, key, members...).Result()
}

func (m *evaRedis) ZAddNXCh(ctx context.Context, key string, members ...*redis.Z) (int64, error) {
	return m.client.ZAddNXCh(ctx, key, members...).Result()
}

func (m *evaRedis) ZAddXXCh(ctx context.Context, key string, members ...*redis.Z) (int64, error) {
	return m.client.ZAddXXCh(ctx, key, members...).Result()
}

func (m *evaRedis) ZIncr(ctx context.Context, key string, member *redis.Z) (float64, error) {
	return m.client.ZIncr(ctx, key, member).Result()
}

func (m *evaRedis) ZIncrNX(ctx context.Context, key string, member *redis.Z) (float64, error) {
	return m.client.ZIncrNX(ctx, key, member).Result()
}

func (m *evaRedis) ZIncrXX(ctx context.Context, key string, member *redis.Z) (float64, error) {
	return m.client.ZIncrXX(ctx, key, member).Result()
}

func (m *evaRedis) ZCard(ctx context.Context, key string) (int64, error) {
	return m.client.ZCard(ctx, key).Result()
}

func (m *evaRedis) ZCount(ctx context.Context, key, min, max string) (int64, error) {
	return m.client.ZCount(ctx, key, min, max).Result()
}

func (m *evaRedis) ZLexCount(ctx context.Context, key, min, max string) (int64, error) {
	return m.client.ZLexCount(ctx, key, min, max).Result()
}

func (m *evaRedis) ZIncrBy(ctx context.Context, key string, increment float64, member string) (float64, error) {
	return m.client.ZIncrBy(ctx, key, increment, member).Result()
}

func (m *evaRedis) ZInterStore(ctx context.Context, destination string, store *redis.ZStore) (int64, error) {
	return m.client.ZInterStore(ctx, destination, store).Result()
}

func (m *evaRedis) ZPopMax(ctx context.Context, key string, count ...int64) ([]redis.Z, error) {
	return m.client.ZPopMax(ctx, key, count...).Result()
}

func (m *evaRedis) ZPopMin(ctx context.Context, key string, count ...int64) ([]redis.Z, error) {
	return m.client.ZPopMin(ctx, key, count...).Result()
}

func (m *evaRedis) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return m.client.ZRange(ctx, key, start, stop).Result()
}

func (m *evaRedis) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return m.client.ZRangeWithScores(ctx, key, start, stop).Result()
}

func (m *evaRedis) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return m.client.ZRangeByScore(ctx, key, opt).Result()
}

func (m *evaRedis) ZRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return m.client.ZRangeByLex(ctx, key, opt).Result()
}

func (m *evaRedis) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) ([]redis.Z, error) {
	return m.client.ZRangeByScoreWithScores(ctx, key, opt).Result()
}

func (m *evaRedis) ZRank(ctx context.Context, key, member string) (int64, error) {
	return m.client.ZRank(ctx, key, member).Result()
}

func (m *evaRedis) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return m.client.ZRem(ctx, key, members...).Result()
}

func (m *evaRedis) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	return m.client.ZRemRangeByRank(ctx, key, start, stop).Result()
}

func (m *evaRedis) ZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error) {
	return m.client.ZRemRangeByScore(ctx, key, min, max).Result()
}

func (m *evaRedis) ZRemRangeByLex(ctx context.Context, key, min, max string) (int64, error) {
	return m.client.ZRemRangeByLex(ctx, key, min, max).Result()
}

func (m *evaRedis) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return m.client.ZRevRange(ctx, key, start, stop).Result()
}

func (m *evaRedis) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return m.client.ZRevRangeWithScores(ctx, key, start, stop).Result()
}

func (m *evaRedis) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return m.client.ZRevRangeByScore(ctx, key, opt).Result()
}

func (m *evaRedis) ZRevRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	return m.client.ZRevRangeByLex(ctx, key, opt).Result()
}

func (m *evaRedis) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) ([]redis.Z, error) {
	return m.client.ZRevRangeByScoreWithScores(ctx, key, opt).Result()
}

func (m *evaRedis) ZRevRank(ctx context.Context, key, member string) (int64, error) {
	return m.client.ZRevRank(ctx, key, member).Result()
}

func (m *evaRedis) ZScore(ctx context.Context, key, member string) (float64, error) {
	return m.client.ZScore(ctx, key, member).Result()
}

func (m *evaRedis) ZUnionStore(ctx context.Context, dest string, store *redis.ZStore) (int64, error) {
	return m.client.ZUnionStore(ctx, dest, store).Result()
}
