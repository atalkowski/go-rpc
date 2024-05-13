fake=""

function quickLook() {
  docker ps -a --format '{{json .}}' | jq '.ID + "|" + .Image + "|" + .Names + "|" + .Status' | sed 's/^"//
s/"$//' | awk -F'|' '{ print(sprintf("%16.16s | %25.25s | %20.20s | %25.25s", $1, $2, $3, $4)); }' 
}

fullstatus=`docker container ls -a | grep postgres`
status=`echo "$fullstatus" | awk 'BEGIN { status="stopped"; }
        /postgres/ { status="started"; }
        /Exited/ { status="paused"; }
        END { print status; }'`
container=`echo "$fullstatus" | awk '{ print $1; }'`

quickLook
echo "Postgres Run status  ... $status"
echo "Postgres Container   ... $container"

function stoppg() {
  echo "Stopping postgres..."
  case "$status" in
    stopped) $fake echo postgres $container is not running;;
    paused)  $fake echo postgres $container is paused;;
    started) $fake docker stop $container;;
    *) $fake echo "Don't recognize postgres state $status";;
  esac
}

function startpg() {
  echo "Starting postgres $container..."
  case "$status" in
    paused)  $fake docker start $container;;
    started) $fake echo "docker container $container (postgres) is already started; use restart if you want to stop and start it";;
    stopped) $fake docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecret -d postgres:16.2-alpine3.19;;  
    *) echo "Don't know what state we are in $1";;
  esac
}

function removepg() {
  docker ps -a | grep postgres
  echo "Removing postgres container $container..."
  echo "Are you sure (y/N)?"
  read $answer
  case "$answer" in
    y | Y) $fake docker rm $container
  esac
}

while [ $# != 0 ]; do
  arg="$1"
  shift
  case "$arg" in 
  -n) fake="echo Will run: ";;
  restart) stoppg; startpg;;
  start) startpg;;
  stop) stoppg;;
  remove) removepg;;
  *) echo "Don't understand $arg"; exit 1;;
  esac
done


