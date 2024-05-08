package redis

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisHelper struct {
	Client *redis.Client
}

func NewRedisHelper(client *redis.Client) *RedisHelper {
	return &RedisHelper{
		Client: client,
	}
}

func InitRedis(ctx context.Context, options *redis.Options) (*redis.Client, error) {
	rdbclient := redis.NewClient(options)
	_, err := rdbclient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdbclient, nil
}

func (h *RedisHelper) Close() error {
	if h.Client != nil {
		err := h.Client.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *RedisHelper) Exists(key string) (bool, error) {
	if h.Client == nil {
		return false, errors.New("Redis Client is null")
	}
	indicator, err := h.Client.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	if indicator <= 0 {
		return false, nil
	}
	return true, nil
}

func (h *RedisHelper) Get(key string) (interface{}, error) {
	if h.Client == nil {
		return nil, errors.New("Redis Client is null")
	}
	data, err := h.Client.Get(context.Background(), key).Result()
	if err != nil && err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var value interface{}
	err = json.Unmarshal([]byte(data), &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (h *RedisHelper) MGet(key ...string) (interface{}, error) {
	if h.Client == nil {
		return nil, errors.New("Redis Client is null")
	}
	data, err := h.Client.MGet(context.Background(), key...).Result()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (h *RedisHelper) HGet(key string, field string) (interface{}, error) {
	if h.Client == nil {
		return nil, errors.New("Redis Client is null")
	}
	data, err := h.Client.HGet(context.Background(), key, field).Result()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (h *RedisHelper) HMGet(key string, field ...string) ([]interface{}, error) {
	if h.Client == nil {
		return nil, errors.New("Redis Client is null")
	}
	data, err := h.Client.HMGet(context.Background(), key, field...).Result()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (h *RedisHelper) HSet(key string, value ...interface{}) (interface{}, error) {
	if h.Client == nil {
		return nil, errors.New("Redis Client is null")
	}
	data, err := h.Client.HSet(context.Background(), key, value...).Result()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (h *RedisHelper) HMSet(key string, field ...interface{}) (interface{}, error) {
	if h.Client == nil {
		return nil, errors.New("Redis Client is null")
	}
	data, err := h.Client.HMSet(context.Background(), key, field...).Result()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (h *RedisHelper) HGetAll(key string) (interface{}, error) {
	if h.Client == nil {
		return nil, errors.New("Redis Client is null")
	}
	data, err := h.Client.HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// return new value of key after increase old value
func (h *RedisHelper) IncreaseInt(key string, value int) (int, error) {
	if h.Client == nil {
		return 0, errors.New("Redis Client is null")
	}
	res := 0
	//
	ctx := context.Background()
	err := h.Client.Watch(ctx, func(tx *redis.Tx) error {
		n, err := tx.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}

		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			res = n + value
			pipe.Set(ctx, key, res, time.Duration(300)*time.Second)
			return nil
		})
		return err
	}, key)
	//
	if err != nil {
		return 0, err
	}
	return res, nil
}

// return key after increase min value by pattern
func (h *RedisHelper) IncreaseMinValue(keys []string, value int) (string, error) {
	if h.Client == nil {
		return "", errors.New("Redis Client is null")
	}
	res := 0
	//
	key := ""
	ctx := context.Background()
	err := h.Client.Watch(ctx, func(tx *redis.Tx) error {
		key = keys[0]
		minValue, err := tx.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}

		var n int
		for _, k := range keys {
			n, err = tx.Get(ctx, k).Int()
			if err != nil && err != redis.Nil {
				return err
			}

			if minValue > n {
				minValue = n
				key = k
			}
		}

		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			res = minValue + value
			pipe.Set(ctx, key, res, -1)
			return nil
		})
		return err
	}, key)
	//
	if err != nil {
		return "", err
	}
	return key, nil
}

func (h *RedisHelper) GetInterface(key string, value interface{}) (interface{}, error) {
	if h.Client == nil {
		return nil, errors.New("Redis Client is null")
	}
	data, err := h.Client.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	typeValue := reflect.TypeOf(value)
	kind := typeValue.Kind()

	var outData interface{}
	switch kind {
	case reflect.Ptr, reflect.Struct, reflect.Slice:
		outData = reflect.New(typeValue).Interface()
	default:
		outData = reflect.Zero(typeValue).Interface()
	}
	err = json.Unmarshal([]byte(data), &outData)
	if err != nil {
		return nil, err
	}

	switch kind {
	case reflect.Ptr, reflect.Struct, reflect.Slice:
		outDataValue := reflect.ValueOf(outData)

		if reflect.Indirect(reflect.ValueOf(outDataValue)).IsZero() {
			return nil, errors.New("Get redis nil result")
		}
		if outDataValue.IsZero() {
			return outDataValue.Interface(), nil
		}
		return outDataValue.Elem().Interface(), nil
	}
	var outValue interface{} = outData
	if reflect.TypeOf(outData).ConvertibleTo(typeValue) {
		outValueConverted := reflect.ValueOf(outData).Convert(typeValue)
		outValue = outValueConverted.Interface()
	}
	return outValue, nil
}

func (h *RedisHelper) Set(key string, value interface{}, expiration time.Duration) error {
	if h.Client == nil {
		return errors.New("Redis Client is null")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = h.Client.Set(context.Background(), key, string(data), expiration).Result()
	if err != nil && err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (h *RedisHelper) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	if h.Client == nil {
		return false, errors.New("Redis Client is null")
	}
	var isSuccessful bool
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}

	isSuccessful, err = h.Client.SetNX(context.Background(), key, string(data), expiration).Result()
	if err != nil {
		return false, err
	}
	return isSuccessful, nil
}

func (h *RedisHelper) Del(key string) error {
	if h.Client == nil {
		return errors.New("Redis Client is null")
	}
	_, err := h.Client.Del(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (h *RedisHelper) Expire(key string, expiration time.Duration) error {
	if h.Client == nil {
		return errors.New("Redis Client is null")
	}
	_, err := h.Client.Expire(context.Background(), key, expiration).Result()
	if err != nil {
		return err
	}
	return nil
}

func (h *RedisHelper) DelMulti(keys ...string) error {
	if h.Client == nil {
		return errors.New("Redis Client is null")
	}
	var err error
	pipeline := h.Client.TxPipeline()
	pipeline.Del(context.Background(), keys...)
	_, err = pipeline.Exec(context.Background())
	return err
}

func (h *RedisHelper) GetKeysByPattern(pattern string) ([]string, uint64, error) {
	if h.Client == nil {
		return nil, 0, errors.New("Redis Client is null")
	}
	var (
		keys   []string
		cursor uint64 = 0
		limit  int64  = 100
		err    error
	)

	for {
		var temp_keys []string
		temp_keys, cursor, err = h.Client.Scan(context.Background(), cursor, pattern, limit).Result()
		if err != nil {
			return nil, 0, err
		}

		keys = append(keys, temp_keys...)
		if cursor == 0 {
			break
		}
	}

	return keys, cursor, nil
}

func (h *RedisHelper) RenameKey(oldkey, newkey string) error {
	if h.Client == nil {
		return errors.New("Redis Client is null")
	}
	var err error
	_, err = h.Client.Rename(context.Background(), oldkey, newkey).Result()
	return err
}

func (h *RedisHelper) GetType(key string) (string, error) {
	if h.Client == nil {
		return "", errors.New("Redis Client is null")
	}
	typeK, err := h.Client.Type(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return typeK, nil
}

func (h *RedisHelper) GetWithContext(ctx context.Context, key string) (interface{}, error) {
	if h.Client == nil {
		return "", errors.New("Redis Client is null")
	}
	data, err := h.Client.Get(ctx, key).Result()
	if err != nil && err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var value interface{}
	err = json.Unmarshal([]byte(data), &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (h *RedisHelper) GetInterfaceWithContext(ctx context.Context, key string, value interface{}) (interface{}, error) {
	if h.Client == nil {
		return "", errors.New("Redis Client is null")
	}
	data, err := h.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	typeValue := reflect.TypeOf(value)
	kind := typeValue.Kind()

	var outData interface{}
	switch kind {
	case reflect.Ptr, reflect.Struct, reflect.Slice:
		outData = reflect.New(typeValue).Interface()
	default:
		outData = reflect.Zero(typeValue).Interface()
	}
	err = json.Unmarshal([]byte(data), &outData)
	if err != nil {
		return nil, err
	}

	switch kind {
	case reflect.Ptr, reflect.Struct, reflect.Slice:
		outDataValue := reflect.ValueOf(outData)

		if reflect.Indirect(reflect.ValueOf(outDataValue)).IsZero() {
			return nil, errors.New("Get redis nil result")
		}
		if outDataValue.IsZero() {
			return outDataValue.Interface(), nil
		}
		return outDataValue.Elem().Interface(), nil
	}
	var outValue interface{} = outData
	if reflect.TypeOf(outData).ConvertibleTo(typeValue) {
		outValueConverted := reflect.ValueOf(outData).Convert(typeValue)
		outValue = outValueConverted.Interface()
	}
	return outValue, nil
}

func (h *RedisHelper) SetWithContext(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if h.Client == nil {
		return errors.New("Redis Client is null")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = h.Client.Set(ctx, key, string(data), expiration).Result()
	if err != nil && err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}
