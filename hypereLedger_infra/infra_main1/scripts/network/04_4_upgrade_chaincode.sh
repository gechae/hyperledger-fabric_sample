. scripts/network/conf

CHANNEL_NAME=$1
CHAINCODE_NAME=$2
CHAINCODE_VERSION=$3

#chaincode upgrade
docker exec $ENV_STR_PEER0 cli peer chaincode upgrade -o $ORDERER_ADDRESS0 -C $CHANNEL_NAME -n $CHAINCODE_NAME -l golang -v $CHAINCODE_VERSION -c '{"Args":["Init"]}' --tls true --cafile $CA_FILE
