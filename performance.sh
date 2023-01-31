maxThreads=5 
iterations=$1
images=$2
TIMEFORMAT=%R

doThread() {
  local times=()
  for j in `seq $iterations`
  do
    x=$( { time go run ./cmd/api/main.go -amount $images -threads $i; } 2>&1)
    times+=($x)
  done
  avg=$(echo ${times[@]} | jq -s add/length)
  printf $i" thread(s): %.3f\n" $avg
}

run() {
  for i in `seq $maxThreads`
  do
    doThread $i &
  done

  wait
}

run | sort -k1
