package util

import (
	// "log"
	"reportgenengine/settings"
	"time"
)

var Limmiter chan struct{}

func StartTokenBucket(){


	setLimit(settings.MySettings.RATE_LIMIT)
	//refill tokens
	ticker := time.NewTicker(1*time.Second)

	go func() {
		for range ticker.C{
			for i:=0 ;i<cap(Limmiter);i++{
				select {
				case Limmiter <- struct{}{}:
				default :
					//bucket full
				}
			}
		}
	}()
}


func setLimit(limit int){
	Limmiter = make(chan struct{},limit)
}