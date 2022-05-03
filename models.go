package main

type Page struct {
	Success bool         `json:"success"`
	Prev    int          `json:"prev"`
	Current int          `json:"current"`
	Next    int          `json:"next"`
	Last    int          `json:"last"`
	Data    []*Datapoint `json:"data"`
}

type Datapoint struct {
	DBId        int
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Exp         int    `json:"exp"`
	Fame        int    `json:"fame"`
	Job         int    `json:"job"`
	Image       string `json:"image"`
	Restriction int    `json:"restriction_flag"`
}

type Snapshot struct {
	Timestamp int64
	Players   map[int]*Datapoint
}

//{"id":696,"name":"Yoshino","level":71,"exp":672651,"fame":68,"job":131,"image":"128da59b-0ab1-48b4-8ad4-713e08b4893b","restriction_flag":0}
