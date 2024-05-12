
function titled() {
   echo "#$*" | sed 's/./=/g' | cut -c2-
   echo "#$*" | cut -c2-
   echo "#$*" | sed 's/./=/g' | cut -c2-
}

titled "$2"
echo "$1"
echo "Run this now - are you sure (y/N)?"
read answer
case "$answer" in
  y | Y) eval "$1";;
esac
