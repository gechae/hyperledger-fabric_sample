#!/bin/bash

ORGNAME=$1
SERVICEFLAG=$2
ADDNUM=$3

if [ $# -ne 3 ]; then
	echo "Usage 99_2_add_fabric.sh <ORGNAME> <ORDERER|PEER flag> <ADDNUM>"
	echo "example) 99_2_add-peer.sh doro peer 1"
else 
	docker exec -it setup /scripts/add-peer.sh $ORGNAME $SERVICEFLAG $ADDNUM
fi
