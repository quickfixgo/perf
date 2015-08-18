#!/bin/sh

SAMPLE_SIZE=$1
[ -z $SAMPLE_SIZE ] && SAMPLE_SIZE=1000

# ./inbound/inbound -cpuprofile=inbound.prof -fixconfig=cfg/inbound.cfg  -samplesize=$SAMPLE_SIZE &
./inbound/inbound -fixconfig=cfg/inbound.cfg  -samplesize=$SAMPLE_SIZE &
RETVAL=$?
[ $RETVAL -ne 0 ] && exit $RETVAL

INBOUND_PID=$!

trap "kill $INBOUND_PID" 1 2 3 9 15

./outbound/outbound -fixconfig=cfg/outbound.cfg -samplesize=$SAMPLE_SIZE 
RETVAL=$?
[ $RETVAL -ne 0 ] && kill $INBOUND_PID

wait $INBOUND_PID
