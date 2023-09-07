cd test-network
if [ $1 = "down" ]; then
./network.sh down
else
./network.sh down
./network.sh up createChannel -ca
./network.sh deployCC -ccn basic -ccp ../asset-transfer-kwon/chaincode-go -ccl go
cd ../asset-transfer-kwon/rest-api-go
go build main.go
echo "Test"
echo "test Moon"
./main
fi