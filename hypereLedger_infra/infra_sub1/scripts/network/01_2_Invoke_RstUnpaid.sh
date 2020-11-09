. script/conf
 
JSONSTR='[{\"work_dates\":\"'20200901'\",\"tolof_cd\":\"A001\",\"work_no\":\"B001\",\"vhcl_pros_no\":\"FF1ZZ2\",\"Ecard_mnrc_clss_cd\":\"1\",\"hoinst_cd\":\"22\",\"vhno\":\"'FF1ZZ1'\",\"viol_user_nm\":\"tester\",\"mnrc_yn\":\"Y\"}]'
echo $JSONSTR
echo "========================== Invoke setRstUnpaid ======================="
echo
docker exec $ENV_STR_PEER0 cli peer chaincode invoke -o $ORDERER_ADDRESS0 --tls true --cafile $CA_FILE -C $CHANNEL_NAME -n $CHAINCODE_NAME -c '{"Args":["setRstUnpaid","'${JSONSTR}'"]}'
echo
echo "========================== END Invoke setRstUnpaid ======================="
