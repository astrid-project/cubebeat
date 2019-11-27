#!/bin/bash

# Simple script that kills the cubebeat process in execution and build the new binary executile file from the code
kill -9 $(pidof cubebeat)
mage build
