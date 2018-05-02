#!/bin/bash
echo Starting looping script...
while true
do
	python3 bot.py
	echo Bot crashed, restarting...
done
echo Loop terminated
