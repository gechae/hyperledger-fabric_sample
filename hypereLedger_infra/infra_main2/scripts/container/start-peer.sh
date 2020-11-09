touch /log/$CORE_PEER_ID.log
chmod 666 /log/$CORE_PEER_ID.log
chown -R $PROD_USER:$PROD_USER /var/hyperledger/
peer node start >> /log/$CORE_PEER_ID.log 2>&1
