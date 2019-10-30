export values=`cat ~/.azure/osServicePrincipal.json`
for s in $(echo $values | jq -r "to_entries|map(\"\(.key)=\(.value|tostring)\")|.[]" ); do
    export $s
done
az login --service-principal -u $clientId -p "$clientSecret"  --tenant $tenantId
