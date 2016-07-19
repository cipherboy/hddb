package main

import (
	"database/sql"
	"time"
)

import _ "github.com/mattn/go-sqlite3"

var DB *sql.DB = nil

type FileList struct {
	FileName string
	Ignored  bool
}

type iFileHashes struct {
	Size      uint64
	MD5       []byte
	SHA1      []byte
	SHA256    []byte
	Tiger     []byte
	Whirlpool []byte
	FileName  []byte
	ScanDate  time.Time
}

type FileHashes struct {
	Size      uint64
	MD5       string
	SHA1      string
	SHA256    string
	Tiger     string
	Whirlpool string
	FileName  string
	ScanDate  time.Time
}
