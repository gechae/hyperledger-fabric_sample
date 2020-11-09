  echo "======================================="
  echo "            start ca docker            "
  echo "======================================="

docker-compose -f ./compose-files/docker-compose.yaml up -d ca.orghana.com 
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
docker-compose -f ./compose-files/docker-compose.yaml up -d peer1.orghana.com 

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

./bin/configtxgen -printOrg hana > hana.json
