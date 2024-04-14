package auth

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	rstore "github.com/rbcervilla/redisstore/v9"
)

func initRedis() {
	var err error

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	store, err = rstore.NewRedisStore(context.Background(), client)
	if err != nil {
		log.Fatal("failed to create redis store: ", err)
	}

	store.KeyPrefix("session_token")
}
