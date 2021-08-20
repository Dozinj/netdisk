package utils

import "netdisk/global"

type redis struct {}

func GetRedis()*redis{return new(redis)}

func (r *redis)HExists(key string,filed string)bool{
	return global.Redis.HExists(key,filed).Val()
}

func (r *redis)HGet(key string,filed string)string{
	return global.Redis.HGet(key, filed).Val()
}

func (r *redis)HSet(key string, field string, value interface{})error{
	return global.Redis.HSet(key,field,value).Err()
}