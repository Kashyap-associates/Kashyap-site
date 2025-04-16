package config

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type DataType struct {
	Hero struct {
		Title     string `json:"title"`
		Sub_title string `json:"sub_title"`
	} `json:"hero"`

	About struct {
		Image   string `json:"image"`
		Intro   string `json:"details"`
		Story   string `json:"story"`
		Patners struct {
			Intro   string
			Members []struct {
				Image         string `json:"image"`
				Name          string `json:"name"`
				Details       string `json:"details"`
				Background    string `json:"background"`
				Startup_story string `json:"startup_story"`

				Links struct {
					Email    string `json:"email"`
					Linkedin string `json:"linkedin"`
					Whatsapp string `json:"whatsapp"`
				} `json:"links"`

				Feedback struct {
					Intro string `json:"intro"`
					Quote struct {
						Name string `json:"name"`
						Said string `json:"said"`
					} `json:"quote"`
				} `json:"feedback"`
			} `json:"members"`
		} `json:"patners"`

		Team struct {
			Intro   string `json:"intro"`
			Members []struct {
				Image    string `json:"image"`
				Name     string `json:"name"`
				Position string `json:"position"`
				Link     string `json:"link"`
			} `json:"members"`
		} `json:"team"`

		Choose_us string `json:"choose_us"`
	} `json:"about"`

	Calculator string `json:"calculator"`

	Links struct {
		Linkedin string `json:"linkedin"`
		Telegram string `json:"telegram"`
		Email    string `json:"email"`
	} `json:"links"`

	Services struct {
		Intro       string   `json:"intro"`
		Short_intro string   `json:"short_intro"`
		Options     []string `json:"options"`
		Data        []struct {
			Image        string   `json:"image"`
			Name         string   `json:"name"`
			Time         string   `json:"time"`
			Deliverable  []string `json:"deliverable"`
			Regulation   string   `json:"regulation"`
			Catagory     string   `json:"catagory"`
			Availability string   `json:"availability"`
			Audience     string   `json:"audience"`
			Details      string   `json:"details"`
			Scope        string   `json:"scope"`
		} `json:"data"`
	} `json:"services"`

	Others struct {
		Details string `json:"details"`
		Tasks   []struct {
			Name    string `json:"name"`
			Title   string `json:"title"`
			Details string `json:"details"`
		} `json:"tasks"`
	} `json:"others"`

	Contacts struct {
		Intro    string `json:"intro"`
		Addr     string `json:"addr"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Responce string `json:"responce"`
	} `json:"contacts"`

	Annotations []struct {
		Short string `json:"short"`
		Long  string `json:"long"`
	} `json:"annotations"`
}

var (
	//go:embed default.json
	config_json string
	config_data DataType
	mu          sync.RWMutex
	s_key       = "admin-data"
	Login_key   = "admin-login"
)

func init() {
	defer update()
	if _, err := Get(s_key); err != nil {
		if errors.Is(err, redis.Nil) {
			Set(config_json, s_key)
			return
		}
		slog.Error(err.Error())
		return
	}
	if _, err := Get(Login_key); err != nil {
		if errors.Is(err, redis.Nil) {
			Set(fmt.Sprint(
				os.Getenv("ADMIN_USERNAME"),
				"--", os.Getenv("ADMIN_PASSWORD"),
			), Login_key)
			return
		}
		slog.Error(err.Error())
		return
	}
}

func connect() *redis.Client {
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     "db-redis:6379",
			Password: "",
			DB:       0,
		})
	return rdb
}

func Read() *DataType {
	return &config_data
}

func Get(key string) (string, error) {
	var err error
	db := connect()
	defer db.Close()
	var data string
	if data, err = db.Get(context.Background(), key).Result(); err != nil {
		return "", err
	}
	return data, nil
}

func update() {
	var err error
	data, err := Get(s_key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			slog.Warn("No admin data found in Redis.")
			return
		}
		slog.Error(err.Error())
		return
	}
	if data != "" {
		mu.Lock()
		defer mu.Unlock()
		if err := json.Unmarshal([]byte(data), &config_data); err != nil {
			slog.Error(err.Error())
			return
		}
	}
}

func Set(data, key string) {
	db := connect()
	defer db.Close()
	defer update()
	if err := db.Set(context.Background(), key, data, time.Hour*99999).Err(); err != nil {
		slog.Error(err.Error())
	}
}
