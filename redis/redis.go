package redis

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"time"
	"errors"
    "log"
)

var (
    connectionPool *redis.Pool
)

func init() {
    log.Println("Redis init")
    
	connectionPool = &redis.Pool{
        MaxIdle: 10,
		MaxActive: 15,
        IdleTimeout: 240 * time.Second,
        Dial: func () (redis.Conn, error) {
            c, err := redis.Dial("tcp", "redis:6379")
            
            if err != nil {
                log.Fatal(err)
            }
            
            return c, err
        },
    }
}

func Set(key string, value interface{}) error {
	c := connectionPool.Get()
	defer c.Close()
	
	valueAsJSON, _ := json.Marshal(value)
	
	_, err := c.Do("SET", key, valueAsJSON)
	
	if err != nil {
		return err
	}
	
	return nil
}

func Get(key string, buffer interface{}) error {
	c := connectionPool.Get()
	defer c.Close()
	
	valueAsJSON, err := c.Do("GET", key)
	
	if err != nil {
		return err
	}
	
	if valueAsJSON == nil {
		return nil
	}
	
	bytes, ok := valueAsJSON.([]byte)
	
	if !ok {
		return errors.New("The data retrived from redis was of an unexpected type. Key: " + key)
	}
	
	err = json.Unmarshal(bytes, buffer)
	
	if err != nil {
		return err
	}
	
	return nil
}