ORG=doro
ADD_ORG=$2
CHANNEL_NAME=$1
ANCHOR_PORT=$3
ADD_DAT=$4
#create directory to config add org
mkdir /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config

#get channel's config block
peer channel fetch config /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/config_block.pb -o orderer1.orgorderer.com:9060 -c $CHANNEL_NAME --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/orgorderer.com/msp/tlscacerts/ca-orgorderer-bbc-ex-co-kr-9080.pem

#change channel's config file format from pb to json
configtxlator proto_decode --input /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/config_block.pb --type common.Block | jq .data.data[0].payload.data.config >/opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/config.json
cat /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/config.json
echo '{"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer1.org'$ADD_ORG'.'$ADD_DAT'","port": 9050}]},"version": "0"}}'
#add anchor peer config in to org's json
# jq '.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.org'$ADD_ORG'.com","port": 7051}]},"version": "0"}}' /opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/$ADD_ORG.json >/opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/addanchor_$ADD_ORG.json
# jq '.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer0.org'$ADD_ORG'.com","port": 36050}]},"version": "0"}}' /opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/$ADD_ORG.json >/opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/addanchor_$ADD_ORG.json
jq '.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "peer1.org'$ADD_ORG'.'$ADD_DAT'","port": '$ANCHOR_PORT'}]},"version": "0"}}' /opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts/$ADD_ORG.json >/opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/addanchor_$ADD_ORG.json
cat /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/addanchor_$ADD_ORG.json

#add org's config in to channel's config
jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups": {"'$ADD_ORG'":.[1]}}}}}' /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/config.json /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/addanchor_$ADD_ORG.json >/opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/added_$ADD_ORG.json

#change channel's config file format from json to pb
configtxlator proto_encode --input /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/config.json --type common.Config --output /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/config.pb

#change added orgs's config file format from json to pb
configtxlator proto_encode --input /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/added_$ADD_ORG.json --type common.Config --output /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/added_$ADD_ORG.pb

#extract the differences between existing channel's config and added org channel's config
configtxlator compute_update --channel_id $CHANNEL_NAME --original /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/config.pb --updated /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/added_$ADD_ORG.pb --output /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/config_update.pb

#change extracted config file formet from pb to json
configtxlator proto_decode --input /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/config_update.pb --type common.ConfigUpdate | jq . >/opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/config_update.json

#add header in extracted config file
echo '{"payload":{"header":{"channel_header":{"channel_id": "'$CHANNEL_NAME'", "type":2}},"data":{"config_update":'$(cat /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/config_update.json)'}}}' | jq . >/opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/config_update_in_envelope.json

#change final file format from json to pb
configtxlator proto_encode --input /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/config_update_in_envelope.json --type common.Envelope --output /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/config_update_in_envelope.pb

#channel update with final pb file
peer channel update -f /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config/config_update_in_envelope.pb -o orderer1.orgorderer.com:9060 -c $CHANNEL_NAME --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/orgorderer.com/msp/tlscacerts/ca-orgorderer-bbc-ex-co-kr-9080.pem

#delete created directory
rm -rf /opt/gopath/src/github.com/hyperledger/fabric/peer/add_config
