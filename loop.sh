#!/bin/bash
echo Starting looping script...
while true
do
	python3 bot.py
	echo Bot crashed, restarting in 5 seconds...
	sleep 5
done
echo Loop terminated
