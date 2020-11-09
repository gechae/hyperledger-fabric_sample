. scripts/network/conf

#create cert-channel
#CHANNEL_NAME=doro-channel
#CHANNEL_NAME="dm-channel dn-channel"
for ch in $CHANNEL_NAME ; do 
	docker exec $ENV_STR_PEER0 cli peer channel create -o $ORDERER_ADDRESS0 -c $ch -f ./channel-artifacts/$ch/$ch.tx --tls true --cafile $CA_FILE

	echo "======================================="
	echo "         End to create $ch channel         "
	echo "======================================="
done
