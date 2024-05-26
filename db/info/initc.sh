TName=""
Entry=""
Entries=""
NewFields=""
NewValues=""
UpdateFields=""
Output=""

template='
-- name: CreateEntry :one
INSERT INTO TName (NewFields)
VALUES (NewValues)
RETURNING *;

-- name: GetEntry :one
SELECT * FROM TName
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM TName
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateEntry :one
UPDATE TName
  set UpdateFields
WHERE id = $1
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM TName WHERE id = $1;
'

function blurb() {
  echo "Use $0 to generate a SqlC template for a table
  Usage:
    $0 Single Plural -table TName -fields 'f1, f2, ...' -values '\$1, \$2, ...' -update 's1=\$1, s2=\$2...' -output outfile

  where
    Single .. is the singular name of the function suffix is Golang (e.g. Entry in GetEntry)
    Plural .. is the plural name of the function suffix name in Golang (e.g. Entries in ListEntries)
    TName ... is the Table Name used in the database
    fN ...... is the field name of Nth field in a create statement
    sN ...... is the field name of Nth field in an update statement

  and
    -TName  .. the database table name used in SQL for all statement types CREATE, READ (Get and List), UPDATE, DELETE
    -fields .. specifies the fields in a CREATE sql statement
    -values .. specifies the values in a CREATE sql statement
    -update .. the set statement in the UPDATE sql statement 
  Example usage:
    $0 Entry Entries -table entries -fields 'account_id, amount' -values '\$1, \$2' -update 'account_id=\$2, amount=\$3'

  Note: the -update needs to count the settings from 2, 3 ... and not 1: the first arg will be used as the ID.
  "  

  if [ "$*" = "" ]; then
     exit 0
  else
    echo "Error:\n$*"
    exit 1
  fi
}

function invokeTemplate() {
echo "$template" | sed "s/Entry/$Entry/g
 s/Entries/$Entries/g
 s/NewFields/$NewFields/g
 s/NewValues/$NewValues/g
 s/UpdateFields/$UpdateFields/g
 s/TName/$TName/g" 
}

Entry="$1"
Entries="$2"
shift 2
if [ "$Entry" = "" -o  "$Entries" = "" ]; then
  blurb 
  exit 0
fi  

while [ $# != 0 ]; do
  arg="$1"
  shift
  case "$arg" in
    -tname | -table) TName="$1"; shift;;
    -field*) NewFields="$1"; shift;;
    -value*) NewValues="$1"; shift;;
    -update*) UpdateFields="$1"; shift;;
    -output) Output="$1"; shift;;
    *) blurb "Don't understand $arg"; exit 1;;
  esac 
done

if [ "$TName" = "" ]; then
  blurb 
  exit 0
fi
if [ "$Output" = "" ]; then
  invokeTemplate
else
  invokeTemplate | tee $Output
  echo "Output was saved in $Output"
fi
