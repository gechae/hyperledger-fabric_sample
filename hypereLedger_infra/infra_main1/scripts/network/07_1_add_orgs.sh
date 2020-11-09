. scripts/network/conf

CHANNEL=$1
ADD_ORG=$2
ANCHOR_PORT=$3
ADD_DAT=$4

echo "test2"
echo $ANCHOR_PORT $ADD_DAT
#copy add-org script in to cli container
docker cp ./scripts/container/add-org.sh cli:/opt/gopath/src/github.com/hyperledger/fabric/peer

#run add-org script
docker exec ${ENV_STR_PEER0} cli /bin/bash /opt/gopath/src/github.com/hyperledger/fabric/peer/add-org.sh $CHANNEL $ADD_ORG $ANCHOR_PORT $ADD_DAT

#delete add-oorg script in cli container
docker exec cli rm -rf /opt/gopath/src/github.com/hyperledger/fabric/peer/add-org.sh
