. scripts/network/conf

CHANNEL_NAME=$1
CHAINCODE_NAME=$2
CHAINCODE_VERSION=$3
CHAINCODE_PATH=github.com/chaincode/$CHAINCODE_NAME/go/
CHAINCODE_PACKAGE_FILE=${CHAINCODE_NAME}_v${CHAINCODE_VERSION}.out
CHAINCODE_PACKAGE_FILE_PATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/$CHAINCODE_PACKAGE_FILE
# docker exec -it cli cp /opt/gopath/src/$CHAINCODE_PATH$CHAINCODE_NAME.go /opt/gopath/src/$CHAINCODE_PATH$CHAINCODE_NAME.go_v$3

#packaging chaincode
docker exec $ENV_STR_PEER0 cli peer chaincode package $CHAINCODE_PACKAGE_FILE_PATH -n $CHAINCODE_NAME -v $CHAINCODE_VERSION -l golang -p $CHAINCODE_PATH

#chaincode install by package
docker exec $ENV_STR_PEER0 cli peer chaincode install $CHAINCODE_PACKAGE_FILE_PATH
