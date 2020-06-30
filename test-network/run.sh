#!/bin/bash

./network.sh down
./network.sh up
./network.sh createChannel
./network.sh createChannel -c channel2
./network.sh deployCC
