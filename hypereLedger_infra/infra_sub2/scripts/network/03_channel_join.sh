. scripts/network/conf

docker exec $ENV_STR_PEER0 cli peer channel join -b $CHANNEL_NAME.block


docker exec cli  peer channel fetch 0 $CHANNEL_NAME.block  -o $ORDERER_ADDRESS0 -c $CHANNEL_NAME --tls --cafile $CA_FILE

docker exec $ENV_STR_PEER0 cli peer channel join -b $CHANNEL_NAME.block


echo "======================================="
echo "        End to join $CHANNEL_NAME      "
echo "       by peer0.orgnonghyupit.com      "
echo "======================================="
