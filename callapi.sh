fake=""
owner=""
currency=""
port="9090"
id=""
mode="api"
page_id="1"
page_size="10"
limit="10"
offset="0"

function dft() {
  echo "$1"
}

function myrand {
  min=`dft $1 1`
  max=`dft $2 6`
  number=$(expr $min + $RANDOM$RANDOM % $max)
  echo $number
}

function getOwner() {
  if [ "$owner" = "" ]; then
    alpha="adgjlpiybc"
    case `myrand 1 4` in
    2) alpha="tksmdeobpz";;
    3) alpha="asdfghjklm";;
    4) alpha="mnbvcxzjgd";;
    esac
    myrand 100000 999999 | tr '[0-9]' "[$alpha]"
  else
    echo "$owner"
  fi
}

function getCurrency() {
  if [ "$currency" = "" ]; then
    case `myrand 1 10` in
    1) echo "USD";;
    2) echo "CAN";;
    3) echo "GBP";;
    4) echo "EUR";;
    *) echo "USD";;
    esac
  else
    echo "$currency"
  fi
}

function mycurl() {
  if [ "$fake" = "" ]; then
    curl $@ | jq
  else
    echo "curl $@ | jq"
  fi
}


function do_get() {
  case "$mode" in
  api) mycurl $verbose "http://localhost:$port/accounts/$id";;
  *) $fake showsql -psql sql "select * from accounts where ID = '$id'";;
  esac
}

function do_delete() {
  case "$mode" in
  api) mycurl $verbose -X DELETE "http://localhost:$port/accounts/$id";;
  *) $fake showsql -psql sql "delete from accounts where ID = '$id'";;
  esac
}

function do_list() {
  case "$mode" in
  api) mycurl $verbose "http://localhost:$port/accounts?page_size=$page_size&page_id=$page_id";;
  *) $fake showsql -psql sql "select * from accounts ORDER BY ID LIMIT $limit OFFSET $offset";;
  esac
}

function do_create() {
  case "$mode" in
  api)
    # Using the existing API example:
    data="{
      \"Owner\":\"$(getOwner)\",
      \"Currency\":\"$(getCurrency)\"
    }"
    # echo "Using this data:"
    # echo "$data"
    if [ "$fake" = "" ]; then
      curl -X POST "http://localhost:$port/accounts" -d "$data" | jq
    else
      echo "curl -X POST "http://localhost:$port/accounts" -d "$data" | jq"
    fi
    ;;
  *) $fake showsql -psql sql "
      INSERT INTO accounts( owner, currency, balance)
      VALUES('$(getOwner)', '$(getCurrency)', 0)";;
  esac
  
}

function blurb() {
  echo "Use $0 to run the AIP request on the simple bank API
Usage:
  $0 [options] list   ...... list accounts 
  $0 [options] create ...... create a new account for given params 
  $0 [options] get ID ...... get account for a given ID
  $0 [options] delete ID ... delete account for a given ID
  
and options are:
  -api ................. use the api (default)
  -psql ................ use direct SQL call rather than the API call
  -v ................... use verbose mode in the curl requests
  -name Owner .......... set account Owner name (dft random)
  -curr XXX ............ set currency to XXX (one of USD, CAN, EUR or GBP) dft = USD
  -fake ................ just show the curl request but don't execute

and for the list options:
- (API only)  
  -page_id N ........... select page N (N >= 1 and default 1)
  -page_size K ......... select page size K (>= 1 and default 10) 
- (SQL only)
  -limit limit ......... for account list limit (default 10)
  -offset offset ....... for account list offset (default 0)
$*"
}

cmd="blurb"

while [ $# != 0 ]; do
  arg="$1"
  shift 
  case "$arg" in
  -n) fake="echo";;
  -v | -vv) verbose="$arg";;
  -psql | -sql) mode="psql";;
  -name) name="$1"; shift;;
  -curr*) curr="$1"; shift;;
  -fake) fake="echo";;
  -page_id) page_id="$1"; shift;;
  -page_size) page_size="$1"; shift;;
  -limit) limit="$1"; shift;;
  -offset) offset="$1"; shift;;
  list | create) cmd="do_$arg";;
  delete | get) id="$1"; shift
     if [ "$id" = "" ]; then
       blurb "Missing ID for $arg request"
       exit 1
     fi
     cmd="do_$arg";;
  *) blurb "Don't understand $arg";;
  esac
done 

$cmd


