package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Age       int       `json:"age,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Timestamp struct {
	time.Time
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.Unix(), 10)), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	ts, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	t.Time = time.Unix(ts, 0)
	return nil
}

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Debug    bool           `json:"debug"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type DatabaseConfig struct {
	URL            string `json:"url"`
	MaxConnections int    `json:"max_connections"`
}

func (c *Config) Validate() error {
	if c.Server.Host == "" {
		return fmt.Errorf("server host required")
	}
	if c.Server.Port == 0 {
		return fmt.Errorf("server port required")
	}
	if c.Database.URL == "" {
		return fmt.Errorf("database URL required")
	}
	return nil
}

type Event struct {
	Type    string
	Payload map[string]interface{}
	Time    time.Time
}

func (e Event) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":      e.Type,
		"data":      e.Payload,
		"timestamp": e.Time.Unix(),
	})
}

func main() {
	user := User{
		ID:        1,
		Username:  "john",
		Email:     "john@example.com",
		Age:       30,
		CreatedAt: time.Now(),
	}
	
	data, _ := json.MarshalIndent(user, "", "  ")
	fmt.Println(string(data))
	
	var decoded User
	json.Unmarshal(data, &decoded)
	fmt.Printf("%+v\n", decoded)
}
