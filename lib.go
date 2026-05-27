package main

import (
    "encoding/json"
    "os"
)

func LoadConfigMap(path string) (map[string]interface{}, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var cfg map[string]interface{}
    err = json.Unmarshal(data, &cfg)
    return cfg, err
}