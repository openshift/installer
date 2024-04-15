#!/usr/bin/pwsh

function New-OpenShiftVM {
    param(
        [Parameter(Mandatory=$true)]
        $Datastore,
        $FailureDomain,
        [Parameter(Mandatory=$true)]
        [string]$IgnitionData,
        [switch]$LinkedClone,
        $Location,
        $MemoryMB,
        [Parameter(Mandatory=$true)]
        [string]$Name,
        $Network,
        $Networking,
        $NumCpu,
        $ReferenceSnapshot,
        $ResourcePool,
        [Parameter(Mandatory=$true)]
        $Tag,
        [Parameter(Mandatory=$true)]
        $Template,
        $VMHost
    )

    #Write-Output $IgnitionData

    # Create arg collection for New-VM
    $args = $PSBoundParameters
    $args.Remove('Template') > $null
    $args.Remove('IgnitionData') > $null
    $args.Remove('Tag') > $null
    $args.Remove('Networking') > $null
    $args.Remove('Network') > $null
    $args.Remove('MemoryMB') > $null
    $args.Remove('NumCpu') > $null
    foreach ($key in $args.Keys){
        if ($NULL -eq $($args.Item($key)) -or $($args.Item($key)) -eq "") {
            $args.Remove($key) > $null
        }
    }

    # Clone the virtual machine from the imported template
    # $vm = New-VM -VM $Template -Name $Name -Datastore $Datastore -ResourcePool $ResourcePool #-Location $Folder #-LinkedClone -ReferenceSnapshot $Snapshot
    $vm = New-VM -VM $Template @args

    # Assign tag so we can later clean up
    New-TagAssignment -Entity $vm -Tag $Tag > $null

    # Update VM specs.  New-VM does not honor the passed in parameters due to Template being used.
    if ($null -ne $MemoryMB -And $null -ne $NumCpu)
    {
        Set-VM -VM $vm -MemoryMB $MemoryMB -NumCpu $NumCpu -CoresPerSocket 4 -Confirm:$false > $null
    }
    Get-HardDisk -VM $vm | Select-Object -First 1 | Set-HardDisk -CapacityGB 120 -Confirm:$false > $null

    # Configure Network (Assuming template networking may not be correct if shared across clusters)
    $pg = Get-VirtualPortgroup -Name $Network -VMHost $(Get-VMHost -VM $vm) 2> $null
    $vm | Get-NetworkAdapter | Set-NetworkAdapter -Portgroup $pg -confirm:$false > $null

    # Assign advanced settings
    New-AdvancedSetting -Entity $vm -name "disk.enableUUID" -value "TRUE" -confirm:$false -Force > $null
    New-AdvancedSetting -Entity $vm -name "stealclock.enable" -value "TRUE" -confirm:$false -Force > $null
    New-AdvancedSetting -Entity $vm -name "guestinfo.ignition.config.data.encoding" -value "base64" -confirm:$false -Force > $null
    New-AdvancedSetting -Entity $vm -name "guestinfo.ignition.config.data" -value $IgnitionData -confirm:$false -Force > $null
    New-AdvancedSetting -Entity $vm -name "guestinfo.hostname" -value $Name -Confirm:$false -Force > $null

    # Create ip kargs
    # "guestinfo.afterburn.initrd.network-kargs" = "ip=${var.ipaddress}::${cidrhost(var.machine_cidr, 1)}:${cidrnetmask(var.machine_cidr)}:${var.vmname}:ens192:none:${join(":", var.dns_addresses)}"
    # Example: ip=<ip_address>::<gateway>:<netmask>:<hostname>:<iface>:<protocol>:<dns_address>
    if ($null -ne $Networking)
    {
        $kargs = "ip=$($Networking.ipAddress)::$($Networking.gateway):$($Networking.netmask):$($Networking.hostname):ens192:none:$($Networking.dns)"
        New-AdvancedSetting -Entity $vm -name "guestinfo.afterburn.initrd.network-kargs" -value $kargs -Confirm:$false -Force > $null
    }

    return $vm
}

function New-VMConfigs {
    $virtualMachines = @"
{
    "virtualmachines": {}
}
"@ | ConvertFrom-Json -Depth 2
    $fds = ConvertFrom-Json $failure_domains

    # Generate Bootstrap
    $vm = createNode -FailureDomain $fds[0] -Type "bootstrap" -VCenter $vcenter -IPAddress $bootstrap_ip_address
    add-member -Name "bootstrap" -value $vm -MemberType NoteProperty -InputObject $virtualMachines.virtualmachines

    # Generate Control Plane
    for (($i =0); $i -lt $control_plane_count; $i++) {
        $vm = createNode -FailureDomain $fds[$i % $fds.Length] -Type "master" -VCenter $vcenter -IPAddress $control_plane_ip_addresses[$i]
        add-member -Name $control_plane_hostnames[$i] -value $vm -MemberType NoteProperty -InputObject $virtualMachines.virtualmachines
    }

    # Generate Compute
    for (($i =0); $i -lt $compute_count; $i++) {
        $vm = createNode -FailureDomain $fds[$i % $fds.Length] -Type "worker" -VCenter $vcenter -IPAddress $compute_ip_addresses[$i]
        add-member -Name $compute_hostnames[$i] -value $vm -MemberType NoteProperty -InputObject $virtualMachines.virtualmachines
    }

    return $virtualMachines | ConvertTo-Json
}

function createNode {
    param (
        $FailureDomain,
        $IPAddress,
        $Type,
        $VCenter
    )

    $vmConfig = @"
{
    "server": "$($VCenter)",
    "datacenter": "$($FailureDomain. datacenter)",
    "cluster": "$($FailureDomain.cluster)",
    "network": "$($FailureDomain.network)",
    "datastore": "$($FailureDomain.datastore)",
    "type": "$($Type)",
    "ip": "$($IPAddress)"
}
"@
    return ConvertFrom-Json -InputObject $vmConfig
}

function New-LoadBalancerIgnition {
    param (
        [string]$sshKey
    )

    $haproxyService = (Get-Content -Path ./lb/haproxy.service -Raw) | ConvertTo-Json

    $api = $control_plane_ip_addresses + $bootstrap_ip_address
    if ($compute_count -gt 0)
    {
        $ingress = $compute_ip_addresses
    } else {
        $ingress = $control_plane_ip_addresses
    }

    $Binding = @{ 'lb_ip_address' = $lb_ip_address; 'api' = $api; 'ingress' = $ingress }
    $haproxyConfig = Invoke-EpsTemplate -Path "lb/haproxy.erb.tmpl" -Binding $Binding

    $haproxyConfig = [Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes($haproxyConfig))

    $lbIgnition = @"
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
  },
  "storage": {
    "files": [{
      "path": "/etc/haproxy/haproxy.conf",
      "mode": 420,
      "contents": { "source": "data:text/plain;charset=utf-8;base64,$($haproxyConfig)" }
    }]
  },
  "systemd": {
    "units": [{
      "name": "haproxy.service",
      "enabled": true,
      "contents": $($haproxyService)
    }]
  }
}
"@
    return $lbIgnition
}

function New-VMNetworkConfig {
    param(
        $DNS,
        $Gateway,
        $Hostname,
        $IPAddress,
        $Netmask
    )
    $network = $null

    $network = @"
{
  "ipAddress": "$($IPAddress)",
  "netmask": "$($Netmask)",
  "dns": "$($DNS)",
  "hostname": "$($Hostname)",
  "gateway": "$($Gateway)"
}
"@
    return ConvertFrom-Json -InputObject $network
}

function New-OpenshiftVMs {
    param(
        $NodeType
    )

    Write-Output "Creating $($NodeType) VMs"

    $jobs = @()
    $vmStep = (100 / $vmHash.virtualmachines.Count)
    $vmCount = 1
    foreach ($key in $vmHash.virtualmachines.Keys) {
        $node = $vmHash.virtualmachines[$key]

        if ($NodeType -ne $node.type) {
            continue
        }

        $jobs += Start-ThreadJob -n "create-vm-$($metadata.infraID)-$($key)" -ScriptBlock {
            param($key,$node,$vm_template,$metadata,$tag,$scriptdir,$cliContext)
            . .\variables.ps1
            . ${scriptdir}\upi-functions.ps1
            Use-PowerCLIContext -PowerCLIContext $cliContext

            $name = "$($metadata.infraID)-$($key)"
            Write-Output "Creating $($name)"

            $rp = Get-ResourcePool -Name $($metadata.infraID) -Location $(Get-Cluster -Name $($node.cluster)) -Server $vcenter
            ##$datastore = Get-Datastore -Name $node.datastore -Server $node.server
            $datastoreInfo = Get-Datastore -Name $node.datastore -Location $node.datacenter

            # Pull network config for each node
            if ($node.type -eq "master") {
                $numCPU = $control_plane_num_cpus
                $memory = $control_plane_memory
            } elseif ($node.type -eq "worker") {
                $numCPU = $compute_num_cpus
                $memory = $compute_memory
            } else {
                # should only be bootstrap
                $numCPU = $control_plane_num_cpus
                $memory = $control_plane_memory
            }
            $ip = $node.ip
            $network = New-VMNetworkConfig -Hostname $name -IPAddress $ip -Netmask $netmask -Gateway $gateway -DNS $dns

            # Get the content of the ignition file per machine type (bootstrap, master, worker)
            $bytes = Get-Content -Path "./$($node.type).ign" -AsByteStream
            $ignition = [Convert]::ToBase64String($bytes)

            # Get correct template / folder
            $folder = Get-Folder -Name $metadata.infraID -Location $node.datacenter
            $template = Get-VM -Name $vm_template -Location $($node.datacenter)

            # Clone the virtual machine from the imported template
            #$vm = New-OpenShiftVM -Template $template -Name $name -ResourcePool $rp -Datastore $datastoreInfo -Location $folder -LinkedClone -ReferenceSnapshot $snapshot -IgnitionData $ignition -Tag $tag -Networking $network -NumCPU $numCPU -MemoryMB $memory
            $vm = New-OpenShiftVM -Template $template -Name $name -ResourcePool $rp -Datastore $datastoreInfo -Location $folder -IgnitionData $ignition -Tag $tag -Networking $network -Network $node.network -NumCPU $numCPU -MemoryMB $memory

            # Assign tag so we can later clean up
            # New-TagAssignment -Entity $vm -Tag $tag
            # New-AdvancedSetting -Entity $vm -name "guestinfo.ignition.config.data" -value $ignition -confirm:$false -Force > $null
            # New-AdvancedSetting -Entity $vm -name "guestinfo.hostname" -value $name -Confirm:$false -Force > $null

            if ($node.type -eq "master" -And $delayVMStart) {
                # To give bootstrap some time to start, lets wait 2 minutes
                Start-ThreadJob -ThrottleLimit 5 -InputObject $vm {
                    Start-Sleep -Seconds 90
                    $input | Start-VM
                }
            } elseif ($node.type -eq "worker" -And $delayVMStart) {
                # Workers are not needed right away, gotta wait till masters
                # have started machine-server.  wait 7 minutes to start.
                Start-ThreadJob -ThrottleLimit 5 -InputObject $vm {
                    Start-Sleep -Seconds 600
                    $input | Start-VM
                }
            }
            else {
                $vm | Start-VM
            }
        } -ArgumentList @($key,$node,$vm_template,$metadata,$tag,$SCRIPTDIR,$cliContext)
        Write-Progress -id 222 -Activity "Creating virtual machines" -PercentComplete ($vmStep * $vmCount)
        $vmCount++
    }
    Wait-Job -Job $jobs
    foreach ($job in $jobs) {
        Receive-Job -Job $job
    }
}