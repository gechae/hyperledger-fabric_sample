. scripts/network/conf


docker exec $ENV_STR_PEER0 cli peer chaincode package $CHAINCODE_PACKAGE_FILE_PATH -n $CHAINCODE_NAME -v $CHAINCODE_VERSION -l golang -p $CHAINCODE_PATH
docker exec $ENV_STR_PEER0 cli peer chaincode install $CHAINCODE_PACKAGE_FILE_PATH
docker exec $ENV_STR_PEER0 cli peer chaincode instantiate -o $ORDERER_ADDRESS0 -C $CHANNEL_NAME -n $CHAINCODE_NAME -l golang -v $CHAINCODE_VERSION -c '{"Args":["Init"]}' --tls true --cafile $CA_FILE --collections-config $COLLECTIONS_CONFIG









