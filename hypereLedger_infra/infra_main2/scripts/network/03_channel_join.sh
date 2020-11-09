. scripts/network/conf


for ch in $CHANNEL_NAME ; do
	docker exec $ENV_STR_PEER0 cli peer channel join -b $ch.block
	#docker exec $ENV_STR_PEER0 cli peer channel update -c $ch -o $ORDERER_ADDRESS0 -f ./channel-artifacts/$ch/doroMSPanchors.tx --tls true --cafile $CA_FILE
	echo $ORDERER_ADDRESS0
	docker exec cli  peer channel fetch 0 $ch.block  -o $ORDERER_ADDRESS0 -c $ch --tls --cafile $CA_FILE
	docker exec $ENV_STR_PEER0 cli peer channel join -b $ch.block
	echo "======================================="
	echo "        End to join $ch      "
	echo "       by peer0.orgdoro.com      "
	echo "======================================="
done



