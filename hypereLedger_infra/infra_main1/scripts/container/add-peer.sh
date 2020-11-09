#!/bin/bash
#. scripts/add-conf

function addPeer() {
	echo "######## Adding PEER ############"
	registerIdentities $1
    	getCACerts
	enrollPeer
}
function addOrderer() {
	echo "####### Adding ORDERER ###########"
	registerIdentities $1
    	getCACerts
	enrollOrderer
}
function main() {	
    echo $1
    #echo "Beginning registering orderer and peer identities ..."
    if [ "${1}" = "orderer" ]; then
		addOrderer $1
    elif [ "${1}" = "peer" ]; then
		addPeer $1
    fi
    #registerIdentities
    #getCACerts
    #enrollOrderer
    #enrollPeer
    #makeConfigTxYaml
    #generateChannelArtifacts

    chown -R $PROD_USER:$PROD_USER /crypto-config
    chown -R $PROD_USER:$PROD_USER /root

}

# Enroll the CA administrator
function enrollCAAdmin() {
	echo "############################################ $CA_NAME : $CA_HOST_PORT"
    waitPort "$CA_NAME to start" 20 $CA_HOST_PORT
    echo "Enrolling with $CA_NAME as bootstrap identity ..."
    export FABRIC_CA_CLIENT_HOME=$HOME/cas/$CA_NAME
	echo "################################ $FABRIC_CA_CLIENT_HOME"
    export FABRIC_CA_CLIENT_TLS_CERTFILES=$CA_CHAINFILE
	echo "################################ $FABRIC_CA_CLIECNT_TLS_CERTFILES"
    fabric-ca-client enroll -d -u http://$CA_ADMIN_USER_PASS@$CA_HOST_PORT
}

function registerIdentities() {
    echo "Registering identities ..."
    if [ "${1}" = "orderer" ]; then
    	registerOrdererIdentities
    elif [ "${1}" = "peer" ]; then
    	registerPeerIdentities
    fi
}

function registerOrdererIdentities() {
    initOrgVars $ORG
    enrollCAAdmin

    initOrdererVars $ORG $NUM
    echo "Registering $ORDERER_NAME with $CA_NAME"
    fabric-ca-client register -d --id.name $ORDERER_NAME --id.secret $ORDERER_PASS --id.type orderer

    #echo "Registering admin identity with $CA_NAME"
    #fabric-ca-client register -d --id.name $ADMIN_NAME --id.secret $ADMIN_PASS --id.attrs "admin=true:ecert"
}

function registerPeerIdentities() {
	echo "############################## RegisterPeer Identities"
    initOrgVars $ORG
    enrollCAAdmin

    initPeerVars $ORG $NUM
    echo "Registering $PEER_NAME with $CA_NAME"
    fabric-ca-client register -d --id.name $PEER_NAME --id.secret $PEER_PASS --id.type peer

    #echo "################################## Registering admin identity with $CA_NAME"
    #$BINPATH/fabric-ca-client register -d --id.name $ADMIN_NAME --id.secret $ADMIN_PASS --id.attrs "hf.Registrar.Roles=client,hf.Registrar.Attributes=*,hf.Revoker=true,hf.GenCRL=true,admin=true:ecert,abac.init=true:ecert"


}

function getCACerts() {
    echo "Getting CA certificates ..."
    #for ORG in $ORGS; do
        initOrgVars $ORG
        echo "Getting CA certs for organization $ORG and storing in $ORG_MSP_DIR"
        export FABRIC_CA_CLIENT_TLS_CERTFILES=$CA_CHAINFILE
        fabric-ca-client getcacert -d -u http://$CA_HOST_PORT -M $ORG_MSP_DIR
        finishMSPSetup $ORG_MSP_DIR

        if [ $ADMINCERTS ]; then
            switchToAdminIdentity
        fi
    #done
}

function enrollOrderer() {

       initOrdererVars $ORG $NUM

       fabric-ca-client enroll -d --enrollment.profile tls -u $ENROLLMENT_URL -M /tmp/$ORG/peer$NUM/tls --csr.hosts $ORDERER_HOST

       mkdir -p $TLSDIR
       cp /tmp/$ORG/peer$NUM/tls/signcerts/* $TLSDIR/server.crt
       cp /tmp/$ORG/peer$NUM/tls/keystore/* $TLSDIR/server.key

       fabric-ca-client enroll -d -u $ENROLLMENT_URL -M $MSPDIR

       mkdir -p /crypto-config/ordererOrganizations/${ORDERER_HOST:9}/users/Admin@${ORDERER_HOST:9}/
       mkdir -p /crypto-config/ordererOrganizations/${ORDERER_HOST:9}/msp
       cp -r /root/orgs/$ORG/msp/ /crypto-config/ordererOrganizations/${ORDERER_HOST:9}/
       cp -r /root/orgs/$ORG/admin/msp/ /crypto-config/ordererOrganizations/${ORDERER_HOST:9}/users/Admin@${ORDERER_HOST:9}/

       finishMSPSetup $MSPDIR
       copyAdminCert $MSPDIR $ORDERER_HOST
}

function enrollPeer() {
	echo "################################ enrollPeer ############################################"
        initPeerVars $ORG $NUM

        fabric-ca-client enroll -d --enrollment.profile tls -u $ENROLLMENT_URL -M /tmp/$ORG/peer$NUM/tls --csr.hosts $PEER_HOST
	echo "############################ $TLSDIR"
        mkdir -p $TLSDIR
        cp /tmp/$ORG/peer$NUM/tls/signcerts/* $TLSDIR/server.crt
        cp /tmp/$ORG/peer$NUM/tls/keystore/* $TLSDIR/server.key

	echo "############################ $ENROLLMENT_URL : $MSPDIR"
        fabric-ca-client enroll -d -u $ENROLLMENT_URL -M $MSPDIR

        mkdir -p crypto-config/peerOrganizations/${PEER_HOST:6}/users/Admin@${PEER_HOST:6}
        mkdir -p crypto-config/peerOrganizations/${PEER_HOST:6}/msp
        cp -r /root/orgs/$ORG/msp/ /crypto-config/peerOrganizations/${PEER_HOST:6}/
        cp -r /root/orgs/$ORG/admin/msp/ /crypto-config/peerOrganizations/${PEER_HOST:6}/users/Admin@${PEER_HOST:6}/

        finishMSPSetup $MSPDIR
	echo "###################### $PEER_HOST"
        copyAdminCert $MSPDIR $PEER_HOST
}

function makeConfigTxYaml() {
    {
        echo "
Organizations:"

        for ORG in $ORDERER_ORGS; do
            printOrdererOrg $ORG
        done

        for ORG in $PEER_ORGS; do
            printPeerOrg $ORG 0
        done

        echo "
Capabilities:
    Channel: &ChannelCapabilities
        V1_4_3: true
        V1_3: false
        V1_1: false

    Orderer: &OrdererCapabilities
        V1_4_2: true
        V1_1: false

    Application: &ApplicationCapabilities
        V1_4_2: true
        V1_3: false
        V1_2: false
        V1_1: false

Application: &ApplicationDefaults
    Organizations:
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: \"ANY Readers\"
        Writers:
            Type: ImplicitMeta
            Rule: \"ANY Writers\"
        Admins:
            Type: Signature
            Rule: \"OR('doroMSP.admin')\"

    Capabilities:
        <<: *ApplicationCapabilities

Orderer: &OrdererDefaults

    OrdererType: solo

    Addresses:
        - orderer0.orgorderer.com:37060

    BatchTimeout: 1s

    BatchSize:
        MaxMessageCount: 20
        AbsoluteMaxBytes: 80 KB
        PreferredMaxBytes: 20 KB

    Kafka:
        Brokers:
            - 127.0.0.1:9092

    EtcdRaft:
        Consenters:
            - Host: orderer0.orgorderer.com
              Port: 37060
              ClientTLSCert: ./crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer0.orgorderer.com/tls/server.crt
              ServerTLSCert: ./crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer0.orgorderer.com/tls/server.crt
            - Host: orderer1.orgorderer.com
              Port: 37060
              ClientTLSCert: ./crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer1.orgorderer.com/tls/server.crt
              ServerTLSCert: ./crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer1.orgorderer.com/tls/server.crt
            - Host: orderer2.orgorderer.com
              Port: 36060
              ClientTLSCert: ./crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer2.orgorderer.com/tls/server.crt
              ServerTLSCert: ./crypto-config/ordererOrganizations/orgorderer.com/orderers/orderer2.orgorderer.com/tls/server.crt

    Organizations:
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: \"ANY Readers\"
        Writers:
            Type: ImplicitMeta
            Rule: \"ANY Writers\"
        Admins:
            Type: Signature
            Rule: \"OR('doroMSP.admin')\"

        BlockValidation:
            Type: ImplicitMeta
            Rule: \"ANY Writers\"

Channel: &ChannelDefaults
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: \"ANY Readers\"
        Writers:
            Type: ImplicitMeta
            Rule: \"ANY Writers\"
        Admins:
            Type: Signature
            Rule: \"OR('doroMSP.admin')\"

    Capabilities:
        <<: *ChannelCapabilities


Profiles:"
        echo "
    OneOrgChannel:
        Consortium: SampleConsortium
        <<: *ChannelDefaults
        Application:
            <<: *ApplicationDefaults
            Organizations:"
        for ORG in $PEER_ORGS; do
            echo "                - *${ORG}"
        done

        echo "

    EtcdRaftNetwork:
        <<: *ChannelDefaults
        Capabilities:
            <<: *ChannelCapabilities
        Orderer:
            <<: *OrdererDefaults
            OrdererType: etcdraft
            EtcdRaft:
                Consenters:"
        for ORG in $ORDERER_ORGS; do
            for ((i = 0; i < $NUM_ORDERERS; i++)); do
                initOrdererVars $ORG $i
                echo "                - Host: ${ORDERER_HOST}
                  Port: ${ORDERER_PORT}
                  ClientTLSCert: ${TLSDIR}/server.crt
                  ServerTLSCert: ${TLSDIR}/server.crt"
            done
        done
        echo " 
            Addresses:"
        for ORG in $ORDERER_ORGS; do
            for ((i = 0; i < $NUM_ORDERERS; i++)); do
                initOrdererVars $ORG $i
                echo "                - ${ORDERER_HOST}:${ORDERER_PORT}"
            done
        done
        echo "
            Organizations:"
        for ORG in $ORDERER_ORGS; do
            echo "            - *${ORG}"
        done
        echo "
            Capabilities:
                <<: *OrdererCapabilities
        Application:
            <<: *ApplicationDefaults
            Organizations:"
        for ORG in $ORDERER_ORGS; do
            echo "            - <<: *${ORG}"
        done
        echo "
        Consortiums:
            SampleConsortium:
                Organizations:"
        for ORG in $PEER_ORGS; do
            echo "                - *${ORG}"
        done

    } >/root/data/configtx.yaml
}

function printOrg() {
    echo "
  - &$ORG

    Name: $ORG

    # ID to load the MSP definition as
    ID: $ORG_MSP_ID

    # MSPDir is the filesystem path which contains the MSP configuration
    MSPDir: $ORG_MSP_DIR"
}

function printOrdererOrg() {
    initOrgVars $1
    printOrg
}

function printPeerOrg() {
    initPeerVars $1 $2
    printOrg
    echo "
    AnchorPeers:
       # AnchorPeers defines the location of peers which can be used
       # for cross org gossip communication.  Note, this value is only
       # encoded in the genesis block in the Application section context
       - Host: $PEER_HOST
         Port: 37050"

}

function generateChannelArtifacts() {
    export FABRIC_CFG_PATH=/root/data/

    mkdir -p /root/data/$CHANNEL_NAME_1
    mkdir -p /root/data/$CHANNEL_NAME_2

    echo "Generating channel configuration transaction at $CHANNEL_TX_FILE"
    /root/data/bin/configtxgen -profile OneOrgChannel -outputCreateChannelTx /root/data/$CHANNEL_NAME_1/$CHANNEL_NAME_1.tx -channelID $CHANNEL_NAME_1
    /root/data/bin/configtxgen -profile OneOrgChannel -outputCreateChannelTx /root/data/$CHANNEL_NAME_2/$CHANNEL_NAME_2.tx -channelID $CHANNEL_NAME_2
    #   /root/data/bin/configtxgen -profile OneOrgChannel -outputCreateChannelTx /root/data/$CHANNEL_NAME_1/$CHANNEL_NAME_1.tx -channelID cert_channel
    #   /root/data/bin/configtxgen -profile OneOrgChannel -outputCreateChannelTx /root/data/$CHANNEL_NAME_2/$CHANNEL_NAME_2.tx -channelID tsa_channel
    if [ "$?" -ne 0 ]; then
        echo "Failed to generate channel configuration transaction"
        exit 1
    fi
    sleep 2
    for ORG in $PEER_ORGS; do
        initOrgVars $ORG
        echo "Generating anchor peer update transaction for $ORG at $ANCHOR_TX_FILE"
        /root/data/bin/configtxgen -profile OneOrgChannel -outputAnchorPeersUpdate /root/data/$CHANNEL_NAME_1/${ORG}MSPanchors.tx \
            -channelID $CHANNEL_NAME_1 -asOrg $ORG
        /root/data/bin/configtxgen -profile OneOrgChannel -outputAnchorPeersUpdate /root/data/$CHANNEL_NAME_2/${ORG}MSPanchors.tx \
            -channelID $CHANNEL_NAME_2 -asOrg $ORG
        if [ "$?" -ne 0 ]; then
            echo "Failed to generate anchor peer update for $ORG"
            exit 1
        fi
    done
    echo "Generating orderer genesis block at $GENESIS_BLOCK_FILE"

    /root/data/bin/configtxgen -profile EtcdRaftNetwork -outputBlock /root/data/genesis.block
    /root/data/bin/configtxgen -profile EtcdRaftNetwork -outputBlock /root/data/genesis.block

    sleep 2

}

set -e

SDIR=$(dirname "$0")
source $SDIR/env.sh

ORG=$1
FLAG=$2
NUM=$3
echo "$NUM : $ORG : $FLAG"
main $FLAG
