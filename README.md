# instapay-go-server


Go Server for InstaPay

go-ethereum: geth --datadir . --networkid 3333 --rpc --rpcaddr 141.223.121.139 --rpcport 8555 --ws --wsaddr 141.223.121.139 --wsport 8881 --wsorigins="*" --port 30303 --rpccorsdomain "*" --rpcapi "db,eth,net,web3,personal,admin,miner,debug,txpool" --wsapi "db,eth,net,web3,personal,admin,miner,debug,txpool" --nodiscover console


인스타페이 오프체인 서버

: 사용자 간의 합의를 중계해주는 역할
: Ethereum으로 부터 결제 채널 관련 정보를 받으면 enclave에 셋업.
