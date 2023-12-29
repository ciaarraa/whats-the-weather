package cache

import (
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"os"

	"git.mills.io/prologic/bitcask"
	"github.com/google/uuid"
)

var db *bitcask.Bitcask

type Cache struct {
	location string
	database *bitcask.Bitcask
}

func (cache *Cache) open() {
	db, err := bitcask.Open(cache.location)
	if err != nil {
		fmt.Print(err)
	}
	cache.database = db
}

func (cache *Cache) close() {
	cache.database.Close()
}

func NewCache() Cache {
	//db, err := bitcask.Open("tmp/db")
	//if err != nil {
	//fmt.Print(err)
	//}
	checkCacheFolder()
	return Cache{location: "tmp/db"}
}

func (cache *Cache) add(object []byte, key string) {
	cache.open()
	hashKey := getHashKey(key)
	if cache.database.Has([]byte(hashKey)) {
		return
	}
	id := uuid.New().String()
	if _, err := os.Stat(".cache/" + id); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(".cache/" + id)
		if err != nil {
			log.Println(err)
		}
		err = os.WriteFile(".cache/"+id, []byte(object), 0644)
		if err != nil {
			fmt.Println(err)
		}
	}

	err := cache.database.Put([]byte(key), []byte(id))
	if err != nil {
		fmt.Print(err)
	}
	cache.close()
}

func (cache *Cache) get(key string) string {
	cache.open()
	hashKey := getHashKey(key)
	if cache.database.Has([]byte(hashKey)) {
		fmt.Print("cache hit!")
		val, err := db.Get([]byte(hashKey))
		if err != nil {
			log.Fatalln(err)
		}
		object, err := os.ReadFile(".cache/" + string(val))
		if err != nil {
			log.Fatalln(err)
		}
		return string(object)
	}
	cache.close()
	return ""
}

func checkCacheFolder() {
	fmt.Println("Looking for cache")
	if _, err := os.Stat(".cache"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(".cache", os.ModePerm)
		if err != nil {
			fmt.Println("An error has occured: ", err)
		}
	}
}

func getHashKey(key string) string {
	keyHash := fnv.New64()
	keyHash.Write([]byte(key))
	return fmt.Sprintf("%v", keyHash.Sum64())
}

func addToCache(db *bitcask.Bitcask, object []byte, key string) {
	hashKey := getHashKey(key)
	if db.Has([]byte(hashKey)) {
		return
	}
	id := uuid.New().String()
	if _, err := os.Stat(".cache/" + id); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(".cache/" + id)
		if err != nil {
			log.Println(err)
		}
		err = os.WriteFile(".cache/"+id, []byte(object), 0644)
		if err != nil {
			fmt.Println(err)
		}
	}

	err := db.Put([]byte(key), []byte(id))
	if err != nil {
		fmt.Print(err)
	}
}

func keyInCache(db *bitcask.Bitcask, key string) bool {
	key = getHashKey(key)
	if db.Has([]byte(key)) {
		return true
	}
	return false
}

func getFromCache(db *bitcask.Bitcask, key string) string {
	hashKey := getHashKey(key)
	fmt.Print(db.Has([]byte(hashKey)))
	if db.Has([]byte(key)) {
		fmt.Print("cache hit!")
		val, err := db.Get([]byte(key))
		if err != nil {
			log.Fatalln(err)
		}
		object, err := os.ReadFile(".cache/" + string(val))
		if err != nil {
			log.Fatalln(err)
		}
		return string(object)
	}
	return ""
}

func Sample() {
	cache := NewCache()
	cache.add([]byte("string"), "122")
	cache.get("122")
	cache.close()

}
