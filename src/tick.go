package main

import (
  "time"
)

var searchtick chan int

func searchtickupdate()  {
  time.Sleep(500 * time.Millisecond)
  searchtick <- 1
  searchtickupdate()
}
