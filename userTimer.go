package main

func updateChannelTimer(channel string) bool {
	now := GetMillis()
	if now > RequestTimer[channel] {
		RequestTimer[channel] = GetMillis() + 60000 // add a minute
		RequestCount[channel] = 1
		return true
	}
	if RequestCount[channel] > 5 {
		return false
	}
	RequestCount[channel]++
	return true

}
