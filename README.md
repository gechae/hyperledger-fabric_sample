# hyperledger-fabric_sample


# 1. hyperledger-fabric bin 파일 다운로드
curl -sSL https://bit.ly/2ysbOFE | bash -s -- <fabric_version> <fabric-ca_version>
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 1.4.7 1.4.7

받은 파일 위치에서 /bin 파일인 있는지 확인

# 2. infra_main에서 channel-actifacts 생성
mkdir <infra_main>/channel-actifacts
cp -rf <바이너리 파일 위치>/bin <infra_main>/channel-actifacts


