. scripts/network/conf

CHANNEL_NAME=$1
CHAINCODE_NAME=$2
CHAINCODE_VERSION=$3
CHAINCODE_PATH=github.com/chaincode/$CHAINCODE_NAME/go/
CHAINCODE_PACKAGE_FILE=${CHAINCODE_NAME}_v${CHAINCODE_VERSION}.out
CHAINCODE_PACKAGE_FILE_PATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/$CHAINCODE_PACKAGE_FILE

docker exec $ENV_STR_PEER0 cli peer chaincode install $CHAINCODE_PACKAGE_FILE_PATH










