package main

type LoadAverage interface {
	GetData() error
	Output() error
}
