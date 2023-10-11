#!/usr/bin/pwsh

. ./variables.ps1
. ./upi-functions.ps1

$ErrorActionPreference = "Stop"

# Connect to vCenter
Connect-VIServer -Server $vcenter -Credential (Import-Clixml $vcentercredpath)

# Create Ignition
$sshKey = [string](Get-Content -Path $sshkeypath -Raw:$true) -Replace '\n',''
New-LoadBalancerIgnition -sshKey $sshKey

<#$haproxyService = (Get-Content -Path ./lb/haproxy.service -Raw) | ConvertTo-Json
$haproxyConfig = [Convert]::ToBase64String((Get-Content -Path ./lb/haproxy.tmpl -AsByteStream))

$ignition = @"
{
  "ignition": { "version": "3.0.0" },
  "passwd": {
    "users": [
      {
        "name": "core",
        "sshAuthorizedKeys": [
          "$($sshKey)"
        ]
      }
    ]
  }
  "storage": {
    "files": [{
      "path": "/etc/haproxy/haproxy.conf",
      "mode": 420,
      "contents": { "source": "data:text/plain;charset=utf-8;base64,$($haproxyConfig)" }
    }]
  }
  "systemd": {
    "units": [{
      "name": ""haproxy.service",
      "enabled": true,
      "contents": "$($haproxyService)"
    }]
  }
}
"@#>

$lb_ip_address = "192.168.14.2"
$api = "192.168.14.10","192.168.14.11","192.168.14.12"
$ingress = "192.168.14.20","192.168.14.21","192.168.14.22"

$Binding = @{ 'lb_ip_address' = $lb_ip_address; 'api' = $api; 'ingress' = $ingress }
Invoke-EpsTemplate -Path "lb/haproxy.erb.tmpl" -Binding $Binding

return
#Write-Output $ignition

# Create VM
$vmhost = Get-Random -InputObject (Get-VMHost -Location (Get-Cluster $cluster))
$datastore = Get-Datastore -Name mdcnc-ds-4 -Location datacenter-2
$Template = Get-VM -Name rhcos-414.92.202305090606-0-vmware.x86_64.ova -Datastore $datastore
New-OpenShiftVM -Ignition $ignition -Name "ngirard-Test-LB" -Template $Template -VMHost $vmhost -ResourcePool $rp -Datastore $datastore -Location $folder