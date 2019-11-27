#!/bin/bash

kill -9 $(pidof cubebeat)
ps aux | grep cubebeat
mage build
