#!/usr/bin/env bash

function printHelp() {

    echo
    echo
    echo "Usage: "
    echo "  main_network.sh -up <after start ca, setup container and copy crypto-config to other server, start fabric network> "
    echo "  main_network.sh -up [ca <ca, setup>] [o <orderer>] [p <peer>] [co <couchdb>] [ex <explorer>] [exdb <explorer db>]"
    echo "  main_network.sh -up [container_name <other containers>]"
    echo "  main_network.sh -down <all containers stop and network reset>"
    echo "  main_network.sh -down [ca <ca, setup>] [o <orderer>] [p <peer>] [co <couchdb>] [ex <explorer>] [exdb <explorer db>]"
    echo "  main_network.sh -down [container_name <other containers>]"
    echo
    # echo "  main_network.sh -install [channel_name] [chaincode_name] [version (optional)]"
    echo "  main_network.sh -install [channel_name] [chaincode_name] [version] <install chaincode to specific version>"
    echo
    echo "  main_network.sh -upgrade [channel_name] [chaincode_name] [version] <upgrade chaincode to specific version>"
    echo
    echo "  main_network.sh -addorg [channel_name] [org name] <add organization to the channel>"
    echo
    echo "  main_network.sh -removeorg [channel_name] [org name] <add organization to the channel>"
    echo
    echo "  main_network.sh -getconfig [channel_name] <get config in channel>"
    echo "  main_network.sh -setconfig [channel_name] <set config in channel after you modified existing config>"
    echo
    echo "start example: "
    echo "  main_network.sh -up "
    echo "  main_network.sh -up ca <ca, setup container start>"
    echo "  main_network.sh -up o <orderer container start>"
    echo "  main_network.sh -up p <peer container start>"
    echo "  main_network.sh -up cli <cli container start>"
    echo
    echo "shutdown example"
    echo "  main_network.sh -down "
    echo "  main_network.sh -down ca <ca, setup container stop>"
    echo "  doro_network.sh -down o <orderer container stop>"
    echo "  doro_network.sh -down p <peer container stop>"
    echo "  doro_network.sh -down cli <cli container stop>"
    echo
    echo "install example"
    # echo "  main_network.sh -install channel1 test-cc <upgrade to the next version ex) 1.0.00 > 2.0.0>"
    echo "  main_network.sh -install channel1 test-cc 3.2.11 <install to specific version>"
    # echo "  main_network.sh -install channel1 test-cc <upgrade to the next version ex) 1.0.00 > 2.0.0>"
    echo "  main_network.sh -install channel1 test-cc 3.2.11 <install to specific version>"
    echo
    echo "upgrade example"
    # echo "  main_network.sh -install channel1 test-cc <upgrade to the next version ex) 1.0.00 > 2.0.0>"
    echo "  main_network.sh -upgrade channel1 test-cc 3.2.11 <upgrade to specific version>"
    # echo "  main_network.sh -install channel1 test-cc <upgrade to the next version ex) 1.0.00 > 2.0.0>"
    echo "  main_network.sh -upgrade channe1l test-cc 3.2.11 <upgrade to specific version>"
    echo
    echo "addorg example"
    echo "  main_network.sh -addorg channel1 test"
    echo "  main_network.sh -addorg channel1 test"
    echo
    echo "removeorg example"
    echo "  main_network.sh -removeorg channel1 test"
    echo "  main_network.sh -removeorg channel1 test"
    echo
    echo "get config example"
    echo "  main_network.sh -getconfig cert-channel"
    echo "  main_network.sh -getconfig tsa-channel"
    echo
    echo "set config example"
    echo "  main_network.sh -setconfig channel1"
    echo "  main_network.sh -setconfig channel1"
    echo
    echo "============"
    echo "   NOTICE   "
    echo "============"
    echo "* You must start first ca, setup container to enroll fabric certs"
    echo "* You must copy crypto-config dir to other servers"
    echo
    echo

}

function removeOrg() {
    ch=$1
    removeorg=$2

    checkChannel $ch

    ./scripts/network/07_2_remove_orgs.sh $ch $removeorg
}

function addOrg() {
    ch=$1
    addorg=$2
    adddat=$4

    checkChannel $ch
    if [ -n "$3" ]; then
        anchor_port=$3
    else
        echo "You must input anchor peer port"
        exit 1

    fi
    echo "test1"
echo $anchor_port $adddat
    ./scripts/network/07_1_add_orgs.sh $ch $addorg $anchor_port $adddat
}
function getConfig() {
    ch=$1
    checkChannel $ch

    ./scripts/network/08_1_get_chanel_config.sh $ch
}

function setConfig() {
    ch=$1
    checkChannel $ch

    ./scripts/network/08_2_set_chanel_config.sh $ch
}

function upgradeChaincode() {
    ch=$1
    cc=$2
    ccver_list=()
    checkChannel $ch
    checkChaincode $ch $cc

    ccver=$(checkChaincode $ch $cc)
    ccver=${ccver#*Version:}
    ccver=${ccver%%,*}

    lists=$(echo $ccver | tr "." "\n")
    i=0

    for list in $lists; do
        ccver_list[$i]=$(echo $list)
        i=$i+1
    done

    if [ -n "$3" ]; then
        newccver=$3
    else
        echo "You must input chaincode version"
        exit 1

    fi

    ./scripts/network/04_3_upgrade_chaincode.sh $ch $cc $newccver

    echo "Succcess upgrade in $ch chaincode : $cc, version: $newccver"

}

function installChaincode() {
    ch=$1
    cc=$2
    ccver_list=()
    checkChannel $ch
    checkChaincode $ch $cc

    ccver=$(checkChaincode $ch $cc)
    ccver=${ccver#*Version:}
    ccver=${ccver%%,*}

    lists=$(echo $ccver | tr "." "\n")
    i=0

    for list in $lists; do
        ccver_list[$i]=$(echo $list)
        i=$i+1
    done

    if [ -n "$3" ]; then
        newccver=$3
    else
        echo "You must input chaincode version"
        exit 1

    fi

    ./scripts/network/04_2_install_chaincode.sh $ch $cc $ccver $newccver

    echo "Succcess install in $ch chaincode : $cc, version: $newccver"

}

function checkChannel() {

    flag=false
    if [ $# -ne 1 ]; then
        echo "You must input channel name"
        exit 1

    fi

    docker exec -it cli peer channel list >./scripts/channel.txt
    sed -i '1,2d' ./scripts/channel.txt
    tr -d '\r' <./scripts/channel.txt >./scripts/channel_list.txt

    rm -rf ./scripts/channel.txt
    channel_list=$(cat ./scripts/channel_list.txt)

    for channel in $channel_list; do

        if [[ "$channel" == "$1" ]]; then
            flag=true
            break
        else
            echo
        fi
    done

    if [ "$flag" == "false" ]; then
        echo "Please check channel name it isn't exist channel : $1"
        exit 1
    fi

}

function checkChaincode() {

    if [ $# -ne 2 ]; then
        echo "You must input chaincode name"
        exit 1

    fi

    docker exec -it cli peer chaincode list --instantiated -C $1 >./scripts/chaincodes.txt
    sed -i '1,1d' ./scripts/chaincodes.txt
    tr -d '\r' <./scripts/chaincodes.txt >./scripts/chaincode.txt

    rm -rf ./scripts/chaincodes.txt
    chaincode=$(cat ./scripts/chaincode.txt)

    if [[ "$chaincode" == *"$2"* ]]; then
        echo
    else
        echo "Please check chaincode name it isn't exist chaincode in channel $1 : $2"
        exit 1
    fi
    echo $chaincode
}

function startDocker() {

    if [ "$2" == "ca" ]; then
        container_name="ca.orgmain.com ca.orgorderer.com setup"
    elif [ "$2" == "o" ]; then
        container_name="orderer0.orgorderer.com orderer1.orgorderer.com"
    elif [ "$2" == "p" ]; then
        container_name="peer0.orgmain.com"
    elif [ "$2" == "co" ]; then
        container_name="couchdb0.orgmain.com"
    else
        container_name=$2
    fi

    if [ $# -ne 2 ]; then
        echo "start main_network"
        ./scripts/network/01_start_docker_container_server1.sh
        sleep 20
        ./scripts/network/02_channel_create_and_anchor.sh
        ./scripts/network/03_channel_join.sh
#        ./scripts/network/04_1_install_chaincode_firsttime.sh
        exit 1

    elif [ $# -ne 3 ]; then
        docker-compose -f ./compose-files/docker-compose.yaml up -d $container_name


        exit 1
    fi

}

function stopDocker() {

    if [ "$2" == "ca" ]; then
        container_name="ca.orgmain.com ca.orgorderer.com setup"
    elif [ "$2" == "o" ]; then
        container_name="orderer0.orgmain.com orderer1.orgorderer.com"
    elif [ "$2" == "p" ]; then
        container_name="peer0.orgmain.com"
    elif [ "$2" == "co" ]; then
        container_name="couchdb0.orgmain.com"

    else
        container_name=$2
    fi

    if [ $# -ne 2 ]; then

        while true; do
            read -p "Do you wish to reset network? y/n  " yn
            case $yn in
            [Yy]*)
                ./scripts/network/99_stop_docker_container_server.sh
                break
                ;;
            [Nn]*) exit ;;
            *) echo "Please answer yes or no." ;;
            esac
        done

        exit 1
    elif [ $# -ne 3 ]; then
        docker stop $container_name
        exit 1
    fi
}
export PROD_USER=$(id -g)

if [ "$1" == "-up" ]; then
    startDocker $1 $2
elif [ "$1" == "-down" ]; then
    stopDocker $1 $2
elif [ "$1" == "-install" ]; then
    installChaincode $2 $3 $4
elif [ "$1" == "-upgrade" ]; then
    upgradeChaincode $2 $3 $4
elif [ "$1" == "-getconfig" ]; then
    getConfig $2
elif [ "$1" == "-setconfig" ]; then
    setConfig $2
elif [ "$1" == "-addorg" ]; then
    addOrg $2 $3 $4 $5
elif [ "$1" == "-removeorg" ]; then
    removeOrg $2 $3
elif [ "$1" == "-h" ]; then
    printHelp
else
    printHelp
    exit 1
fi
