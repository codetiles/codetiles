package main

import (
  "time"
)

// use <-searchtick in select to check when a user joins the queue
var searchtick chan int
var gametick chan int

func tickUser() {
  searchtick <- 0
}

func performGameTick() {
  time.Sleep(time.Second)
  gametick <- 0
  performGameTick()
}
