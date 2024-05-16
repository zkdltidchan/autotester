# for loop to run pubsubtest/cmd/broker/main.go
BASE_DIR=/Users/zkdltid/Desktop/go-test2/pubsubtest

for i in {1..3}
do
  echo "Running broker $i"
  # open a new terminal and run the broker
  osascript -e 'tell app "Terminal"
    do script "cd '$BASE_DIR' && go run cmd/broker/main.go -port 808'$i'"
  end tell'
done
