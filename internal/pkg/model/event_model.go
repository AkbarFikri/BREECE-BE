package model

import "time"

type FilterParam struct {
	search   string
	sort     string
	page     int
	tempat   string
	tanggal  time.Time
	kategori string
}
