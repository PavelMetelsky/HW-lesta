package db

type errorDb string

func (err errorDb) Error() string { return string(err) }

const errorTableEmpty = errorDb("the table is empty")
