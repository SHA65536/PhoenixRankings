package main

import (
	_ "embed"
	"encoding/json"
)

//go:embed queries/create_tables.sql
var query_create_tables string

//go:embed queries/create_indexes.sql
var query_create_indexes string

//go:embed queries/get_last_two.sql
var query_get_last_two string

//go:embed queries/create_point.sql
var query_create_point string

//go:embed queries/update_point.sql
var query_update_point string

type Page struct {
	Success bool         `json:"success"`
	Prev    int          `json:"prev"`
	Current int          `json:"current"`
	Next    int          `json:"next"`
	Last    int          `json:"last"`
	Data    []*Datapoint `json:"data"`
}

type Datapoint struct {
	DBId        int64
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Rank        int
	Level       int    `json:"level"`
	Exp         int    `json:"exp"`
	Fame        int    `json:"fame"`
	Job         int    `json:"job"`
	Image       string `json:"image"`
	Restriction int    `json:"restriction_flag"`
}

func (p *Page) Parse(data []byte) error {
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}
	for idx, player := range p.Data {
		player.Rank = ((p.Current - 1) * 5) + idx + 1
	}
	return err
}

type Snapshot struct {
	Timestamp int64
	Players   map[int]*Datapoint
}

//{"id":696,"name":"Yoshino","level":71,"exp":672651,"fame":68,"job":131,"image":"128da59b-0ab1-48b4-8ad4-713e08b4893b","restriction_flag":0}
