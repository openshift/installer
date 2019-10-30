#Example: ./gitpath.sh upi/azure/arm/azuredeploy.json
# Returns: https://raw.githubusercontent.com/glennswest/installer/master/upi/azure/arm/azuredeploy.json
export thepath=`git config --get remote.origin.url`
export userproject=`cut -d ":" -f 2  <<< $thepath`
export gitusername=`cut -d "/" -f 1 <<< $userproject`
export gitproject=`cut -d "/" -f 2 <<< $userproject | cut -d "." -f 1`
export gitbranch=`git branch | grep \* | cut -d ' ' -f2`
export gitrawurl="https://raw.githubusercontent.com/"${gitusername}/${gitproject}/${gitbranch}/$1
echo $gitrawurl

