#!/usr/bin/env bash

function printHelp() {

    echo
    echo
    echo "Usage: "
    echo "  hana_network.sh -up <after start ca, setup container and copy crypto-config to other server, start fabric network> "
    echo "  hana_network.sh -up [ca <ca, setup>] [o <orderer>] [p <peer>] [co <couchdb>] [ex <explorer>] [exdb <explorer db>]"
    echo "  hana_network.sh -up [container_name <other containers>]"
    echo "  hana_network.sh -down <all containers stop and network reset>"
    echo "  hana_network.sh -down [ca <ca, setup>] [o <orderer>] [p <peer>] [co <couchdb>] [ex <explorer>] [exdb <explorer db>]"
    echo "  hana_network.sh -down [container_name <other containers>]"
    echo
    # echo "  hana_network.sh -install [channel_name] [chaincode_name] [version (optional)]"
    echo "  hana_network.sh -install [channel_name] [chaincode_name] [version] <install chaincode to specific version>"
    echo
    echo "start example: "
    echo "  hana_network.sh -up "
    echo "  hana_network.sh -up ca <ca, setup container start>"
    echo "  hana_network.sh -up o <orderer container start>"
    echo "  hana_network.sh -up p <peer container start>"
    echo "  hana_network.sh -up cli <cli container start>"
    echo
    echo "shutdown example"
    echo "  hana_network.sh -down "
    echo "  hana_network.sh -down ca <ca, setup container stop>"
    echo "  hana_network.sh -down o <orderer container stop>"
    echo "  hana_network.sh -down p <peer container stop>"
    echo "  hana_network.sh -down cli <cli container stop>"
    echo
    echo "install example"
    # echo "  hana_network.sh -install cert-channel hana-cert-cc <upgrade to the next version ex) 1.0.00 > 2.0.0>"
    echo "  hana_network.sh -install cert-channel hana-cert-cc 3.2.11 <install to specific version>"
    # echo "  hana_network.sh -install tsa-channel hana-tsa-cc <upgrade to the next version ex) 1.0.00 > 2.0.0>"
    echo "  hana_network.sh -install tsa-channel hana-tsa-cc 3.2.11 <install to specific version>"
    echo
    echo "upgrade example"
    # echo "  hana_network.sh -install cert-channel hana-cert-cc <upgrade to the next version ex) 1.0.00 > 2.0.0>"
    echo "  hana_network.sh -upgrade cert-channel hana-cert-cc 3.2.11 <upgrade to specific version>"
    # echo "  hana_network.sh -install tsa-channel hana-tsa-cc <upgrade to the next version ex) 1.0.00 > 2.0.0>"
    echo "  hana_network.sh -upgrade tsa-channel hana-tsa-cc 3.2.11 <upgrade to specific version>"
    echo
    echo "============"
    echo "   NOTICE   "
    echo "============"
    echo "* You must start first ca, setup container to enroll fabric certs"
    echo "* You must copy crypto-config dir to other servers"
    echo
    echo

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

        # #### 핑거버전 ####
        # newccver=$(expr ${ccver_list[2]} + 1)
        # newccver=$(seq -f "%02g" $newccver $newccver)
        # newccver=${ccver_list[0]}.${ccver_list[1]}.$newccver

        #  #### 농협정보 개발버전 ####
        # newccver=`expr ${ccver_list[1]} + 1`

        # ccvernhdev=${ccver_list[0]}.$ccvernhdev.${ccver_list[2]}
        #  #### 농협정보 운영버전 ####
        # newccver=`expr ${ccver_list[0]} + 1`
        # ccvernhprod=$ccvernhprod.${ccver_list[1]}.${ccver_list[2]}

    fi

    ./scripts/network/04_install_chaincode_with_package.sh $ch $cc $ccver $newccver

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
        container_name="ca.orghana.com setup"
    elif [ "$2" == "p" ]; then
        container_name="peer0.orghana.com"
    elif [ "$2" == "co" ]; then
        container_name="couchdb0.orghana.com"
    else
        container_name=$2
    fi

    if [ $# -ne 2 ]; then
        echo "start hana_network"
        ./scripts/network/01_start_docker_container_server1.sh
        exit 1

    elif [ $# -ne 3 ]; then
        docker-compose -f ./compose-files/docker-compose.yaml up -d $container_name

        if [ "$2" == "ca" ]; then
         ./bin/configtxgen -printOrg hana > hana.json
         fi

        exit 1
    fi

}

function stopDocker() {

    if [ "$2" == "ca" ]; then
        container_name="ca.orghana.com setup"
    elif [ "$2" == "p" ]; then
        container_name="peer0.orghana.com"
    elif [ "$2" == "co" ]; then
        container_name="couchdb0.orghana.com"
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
elif [ "$1" == "-h" ]; then
    printHelp
else
    printHelp
    exit 1
fi
