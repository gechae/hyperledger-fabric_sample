  echo "======================================="
  echo "            start ca docker            "
  echo "======================================="

docker-compose -f ./compose-files/docker-compose.yaml up -d ca.orgminj.bminjabc.ex.co.kr
sleep 2

  echo "======================================="
  echo "          start setup docker           "
  echo "======================================="

 docker-compose -f ./compose-files/docker-compose.yaml up -d setup 
sleep 3


echo "======================================="
echo "           start peer docker           "
echo "======================================="

sleep 3
docker-compose -f ./compose-files/docker-compose.yaml up -d peer1.orgminj.bminjabc.ex.co.kr orderer3.orgorderer.ex.co.kr

echo "======================================="
echo "           start other docker          "
echo "======================================="

docker-compose -f ./compose-files/docker-compose.yaml up -d cli

sleep 2

echo "                                       "
echo "                                       "
echo "                                       "
echo "                                       "

docker ps

echo "                                       "
echo "                                       "
echo "                                       "
echo "                                       "

./bin/configtxgen -printOrg minj > minj.json
