fake=""
owner=""
currency=""
port="9090"
id="1"

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
    case `myrand 1 10` in
    1) echo "Bonfonte";;
    2) echo "Jones";;
    3) echo "Smithers";;
    4) echo "Brady";;
    5) echo "Tomson";;
    *) echo "Bunsley";;
    esac
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

function do_get() {
 #  $fake curl $verbose http://localhost:$port/accounts/$id
  showsql -psql sql "select * from accounts where ID = '$id'"
}

function do_delete() {
 #  $fake curl $verbose http://localhost:$port/accounts/$id
  showsql -psql sql "delete from accounts where ID = '$id'"
}

function do_list() {
#  $fake curl $verbose http://localhost:$port/accounts
  showsql -psql sql "select * from accounts" 
}

function do_create() {
  # Using the existing API example:
  data="{
     \"Owner\":\"$(getOwner)\",
     \"Currency\":\"$(getCurrency)\"
   }"
  echo "Using this data:"
  echo "$data"
  $fake curl $verbose -X POST "http://localhost:$port/accounts" -d "$data"
}

function blurb() {
  echo "Use $0 to run the AIP request on the simple bank API
Usage:
  $0 [options] list   ...... list accounts 
  $0 [options] create ...... create a new account for given params 
  $0 [options] get ID ...... get account for a given ID
  $0 [options] delete ID ... delete account for a given ID
  
and options are:
  -v ................... use verbose mode in the curl requests
  -name Owner .......... set account Owner name (dft random)
  -curr XXX ............ set currency to XXX (one of USD, CAN, EUR or GBP) dft = USD
  -fake ................ just show the curl request but don't execute
$*"
}

id=""
cmd="blurb"

while [ $# != 0 ]; do
  arg="$1"
  shift 
  case "$arg" in
  -v | -vv) verbose="$arg";;
  -name) name="$1"; shift;;
  -curr*) curr="$1"; shift;;
  -fake) fake="echo";;
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


