#!/usr/bin/pwsh

# Assert-NotNullOrEmpty validates that a value is not null or empty, and throws
# a descriptive error indicating which parameter is missing and what the caller
# should check to fix it.
function Assert-NotNullOrEmpty {
    param(
        [Parameter(Mandatory=$true)]
        [string]$ParameterName,
        $Value,
        [string]$Context = ""
    )
    if ($null -eq $Value -or ($Value -is [string] -and [string]::IsNullOrWhiteSpace($Value))) {
        $msg = "Parameter '$ParameterName' is null or empty."
        if ($Context -ne "") {
            $msg += " $Context"
        }
        throw $msg
    }
}

# Assert-FailureDomainFields validates that all required fields in a failure domain
# object are populated before attempting vSphere operations.
function Assert-FailureDomainFields {
    param(
        [Parameter(Mandatory=$true)]
        $FailureDomain,
        [int]$Index = -1
    )
    $label = if ($Index -ge 0) { "failure domain [$Index]" } else { "failure domain" }
    $requiredFields = @("server", "datacenter", "cluster", "datastore", "network")
    foreach ($field in $requiredFields) {
        $val = $FailureDomain.$field
        if ($null -eq $val -or [string]::IsNullOrWhiteSpace($val)) {
            throw "Required field '$field' is null or empty in $label. Check your failure_domains configuration or variables.ps1. Failure domain value: $($FailureDomain | ConvertTo-Json -Compress)"
        }
    }
}

# Get-VSphereObject wraps common vSphere Get-* cmdlets with error handling that
# reports which object could not be found and why, instead of a generic
# "Value cannot be null" error.
function Get-VSphereObject {
    param(
        [Parameter(Mandatory=$true)]
        [string]$Command,
        [Parameter(Mandatory=$true)]
        [string]$Name,
        [string]$LocationName = "",
        $Server,
        [hashtable]$ExtraParams = @{},
        [string]$Context = ""
    )

    # Validate inputs before calling the cmdlet
    Assert-NotNullOrEmpty -ParameterName "Name" -Value $Name -Context "Cannot look up $Command with a null/empty name. $Context"
    if ($null -eq $Server) {
        throw "Server is null when attempting '$Command -Name $Name'. Verify that the vCenter connection succeeded and the server value is set. $Context"
    }

    $params = @{
        Name = $Name
        Server = $Server
    }

    if ($LocationName -ne "" -and $ExtraParams.ContainsKey("Location")) {
        $params["Location"] = $ExtraParams["Location"]
        $ExtraParams.Remove("Location")
    }

    foreach ($key in $ExtraParams.Keys) {
        $params[$key] = $ExtraParams[$key]
    }

    try {
        $result = & $Command @params
        if ($null -eq $result) {
            throw "$Command returned null for Name='$Name', Server='$Server'. $Context"
        }
        return $result
    }
    catch {
        $errMsg = "$Command failed for Name='$Name'"
        if ($LocationName -ne "") {
            $errMsg += ", Location='$LocationName'"
        }
        $errMsg += ", Server='$Server'. $Context Original error: $($_.Exception.Message)"
        throw $errMsg
    }
}

function New-OpenShiftVM {
    param(
        [int]$CoresPerSocket = 1, # Default is 1 due to not knowing how many may come in via NumCpu variable
        [Parameter(Mandatory=$true)]
        $Datastore,
        [Parameter(Mandatory=$true)]
        [string]$IgnitionData,
        [switch]$LinkedClone,
        $Location,
        [int]$MemoryMB = 8192,
        [Parameter(Mandatory=$true)]
        [string]$Name,
        $Network,
        $Networking,
        [int]$NumCpu = 4,
        $ReferenceSnapshot,
        $ResourcePool,
        $SecureBoot,
        [Parameter(Mandatory=$true)]
        $Server,
        $StoragePolicy,
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
    $args.Remove('CoresPerSocket') > $null
    $args.Remove('SecureBoot') > $null
    foreach ($key in $args.Keys){
        if ($NULL -eq $($args.Item($key)) -or $($args.Item($key)) -eq "") {
            $args.Remove($key) > $null
        }
    }

    # If storage policy is set, lets pull the mo ref
    if ($NULL -ne $StoragePolicy -and $StoragePolicy -ne "")
    {
        try {
            $storagePolicyRef = Get-SpbmStoragePolicy -Server $Server -Id $StoragePolicy
        }
        catch {
            throw "Get-SpbmStoragePolicy failed for VM '$Name': Could not find storage policy with Id '$StoragePolicy' on server '$Server'. Verify the storage policy ID is correct. Original error: $($_.Exception.Message)"
        }
        $args["StoragePolicy"] = $storagePolicyRef
    }

    # Clone the virtual machine from the imported template
    # $vm = New-VM -VM $Template -Name $Name -Datastore $Datastore -ResourcePool $ResourcePool #-Location $Folder #-LinkedClone -ReferenceSnapshot $Snapshot
    try {
        $vm = New-VM -Server $Server -VM $Template @args
    }
    catch {
        throw "New-VM failed for VM '$Name': Could not clone from template '$($Template.Name)' on server '$Server'. Verify the template, datastore, resource pool, and folder are all valid. Original error: $($_.Exception.Message)"
    }

    if ($null -eq $vm) {
        throw "New-VM returned null for VM '$Name'. The clone operation from template '$($Template.Name)' produced no result. Check vCenter tasks for errors."
    }

    # Assign tag so we can later clean up
    New-TagAssignment -Server $Server -Entity $vm -Tag $Tag > $null

    # Update VM specs.  New-VM does not honor the passed in parameters due to Template being used.
    if ($null -ne $MemoryMB -And $null -ne $NumCpu)
    {
        Set-VM -Server $Server -VM $vm -MemoryMB $MemoryMB -NumCpu $NumCpu -CoresPerSocket $CoresPerSocket -Confirm:$false > $null
    }
    #Get-HardDisk -VM $vm | Select-Object -First 1 | Set-HardDisk -CapacityGB 120 -Confirm:$false > $null
    updateDisk -VM $vm -CapacityGB 120

    # Configure Network (Assuming template networking may not be correct if shared across clusters)
    if ($null -eq $Network -or $Network -eq "") {
        throw "Network (port group) name is null or empty for VM '$Name'. Set the 'network' field in your failure domain or variables.ps1."
    }
    try {
        $vmHostObj = Get-VMHost -VM $vm
        $pg = Get-VirtualPortgroup -Server $Server -Name $Network -VMHost $vmHostObj 2> $null
        if ($null -eq $pg) {
            throw "Port group '$Network' not found on host '$($vmHostObj.Name)'. Available port groups: $((Get-VirtualPortgroup -Server $Server -VMHost $vmHostObj | Select-Object -ExpandProperty Name) -join ', ')"
        }
    }
    catch {
        throw "Failed to configure network for VM '$Name': Could not find port group '$Network'. Verify the network/portgroup name matches an existing port group on the ESXi host. Original error: $($_.Exception.Message)"
    }
    $vm | Get-NetworkAdapter -Server $Server | Set-NetworkAdapter -Server $Server -Portgroup $pg -confirm:$false > $null

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

    # Enable secure boot if needed
    if ($true -eq $SecureBoot)
    {
        Set-SecureBoot -VM $vm
    }

    return $vm
}

# This function was created to work around issue in vSphere 8.0 where vCenter crashed
# when Set-HardDisk is called.
function updateDisk {
    param (
        $CapacityGB,
        $VM
    )

    $newDiskSizeKB = $CapacityGB * 1024 * 1024
    $newDiskSizeBytes = $newDiskSizeKB * 1024

    $vmMo = get-view -id $VM.ExtensionData.MoRef

    $devices = $vmMo.Config.Hardware.Device

    $spec = New-Object VMware.Vim.VirtualMachineConfigSpec
    $spec.DeviceChange = New-Object VMware.Vim.VirtualDeviceConfigSpec[] (1)
    $spec.DeviceChange[0] = New-Object VMware.Vim.VirtualDeviceConfigSpec
    $spec.DeviceChange[0].Operation = 'edit'

    foreach($d in $devices) {
        if ($d.DeviceInfo.Label.Contains("Hard disk")) {
            $spec.DeviceChange[0].Device = $d
        }
    }

    $spec.DeviceChange[0].Device.CapacityInBytes = $newDiskSizeBytes
    $spec.DeviceChange[0].Device.CapacityInKB = $newDiskSizeKB

    $vmMo.ReconfigVM_Task($spec) > $null
}

function New-VMConfigs {
    $virtualMachines = @"
{
    "virtualmachines": {}
}
"@ | ConvertFrom-Json -Depth 2
    $fds = ConvertFrom-Json $failure_domains

    # Generate Bootstrap
    $vm = createNode -FailureDomain $fds[0] -Type "bootstrap" -IPAddress $bootstrap_ip_address
    add-member -Name "bootstrap" -value $vm -MemberType NoteProperty -InputObject $virtualMachines.virtualmachines

    # Generate Control Plane
    for (($i =0); $i -lt $control_plane_count; $i++) {
        $vm = createNode -FailureDomain $fds[$i % $fds.Length] -Type "master" -IPAddress $control_plane_ip_addresses[$i]
        add-member -Name $control_plane_hostnames[$i] -value $vm -MemberType NoteProperty -InputObject $virtualMachines.virtualmachines
    }

    # Generate Compute
    for (($i =0); $i -lt $compute_count; $i++) {
        $vm = createNode -FailureDomain $fds[$i % $fds.Length] -Type "worker" -IPAddress $compute_ip_addresses[$i]
        add-member -Name $compute_hostnames[$i] -value $vm -MemberType NoteProperty -InputObject $virtualMachines.virtualmachines
    }

    return $virtualMachines | ConvertTo-Json
}

function createNode {
    param (
        $FailureDomain,
        $IPAddress,
        $Type
    )

    $vmConfig = @"
{
    "server": "$($FailureDomain.server)",
    "datacenter": "$($FailureDomain.datacenter)",
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
            $nodeContext = "VM '$name': server='$($node.server)', datacenter='$($node.datacenter)', cluster='$($node.cluster)', datastore='$($node.datastore)', network='$($node.network)'. Check your virtualmachines/failure_domains configuration."
            Write-Output "Creating $($name)"

            # Validate node fields before vSphere lookups
            Assert-NotNullOrEmpty -ParameterName "node.server" -Value $node.server -Context "Node '$key' has no server set. $nodeContext"
            Assert-NotNullOrEmpty -ParameterName "node.cluster" -Value $node.cluster -Context "Node '$key' has no cluster set. $nodeContext"
            Assert-NotNullOrEmpty -ParameterName "node.datastore" -Value $node.datastore -Context "Node '$key' has no datastore set. $nodeContext"
            Assert-NotNullOrEmpty -ParameterName "node.datacenter" -Value $node.datacenter -Context "Node '$key' has no datacenter set. $nodeContext"

            try {
                $clusterObj = Get-Cluster -Name $($node.cluster) -Server $node.server
            }
            catch {
                throw "Get-Cluster failed for VM '$name': Could not find cluster '$($node.cluster)' on server '$($node.server)'. $nodeContext Original error: $($_.Exception.Message)"
            }

            try {
                $rp = Get-ResourcePool -Name $($metadata.infraID) -Location $clusterObj -Server $node.server
            }
            catch {
                throw "Get-ResourcePool failed for VM '$name': Could not find resource pool '$($metadata.infraID)' in cluster '$($node.cluster)' on server '$($node.server)'. Ensure the resource pool was created during failure domain setup. $nodeContext Original error: $($_.Exception.Message)"
            }

            try {
                $datastoreInfo = Get-Datastore -Name $node.datastore -Location $node.datacenter -Server $node.server
            }
            catch {
                throw "Get-Datastore failed for VM '$name': Could not find datastore '$($node.datastore)' in datacenter '$($node.datacenter)' on server '$($node.server)'. $nodeContext Original error: $($_.Exception.Message)"
            }

            # Pull network config for each node
            if ($node.type -eq "master") {
                $numCPU = $control_plane_num_cpus
                $memory = $control_plane_memory
                $coresPerSocket = $control_plane_cores_per_socket
            } elseif ($node.type -eq "worker") {
                $numCPU = $compute_num_cpus
                $memory = $compute_memory
                $coresPerSocket = $compute_cores_per_socket
            } else {
                # should only be bootstrap
                $numCPU = $control_plane_num_cpus
                $memory = $control_plane_memory
                $coresPerSocket = $control_plane_cores_per_socket
            }

            # Since coresPerSocket is not required for configs, we need to make sure its not zero (default).  We'll make it match NumCPU.
            if ($NULL -eq $coresPerSocket -or $coresPerSocket -lt 1) {
                $coresPerSocket = $numCPU
            }

            $ip = $node.ip
            $network = New-VMNetworkConfig -Server $node.server -Hostname $name -IPAddress $ip -Netmask $netmask -Gateway $gateway -DNS $dns

            # Get the content of the ignition file per machine type (bootstrap, master, worker)
            $ignPath = "./$($node.type).ign"
            if (-Not (Test-Path -Path $ignPath)) {
                throw "Ignition file '$ignPath' not found for VM '$name'. Ensure ignition configs were generated (generateIgnitions step). $nodeContext"
            }
            $bytes = Get-Content -Path $ignPath -AsByteStream
            $ignition = [Convert]::ToBase64String($bytes)

            # Get correct tag
            $tagCategory = Get-TagCategory -Server $node.server -Name "openshift-$($metadata.infraID)" -ErrorAction continue 2>$null
            $tag = Get-Tag -Server $node.server -Category $tagCategory -Name "$($metadata.infraID)" -ErrorAction continue 2>$null

            # Get correct template / folder
            try {
                $dc = Get-Datacenter -Server $node.server -Name $node.datacenter
            }
            catch {
                throw "Get-Datacenter failed for VM '$name': Could not find datacenter '$($node.datacenter)' on server '$($node.server)'. $nodeContext Original error: $($_.Exception.Message)"
            }

            try {
                $folder = Get-Folder -Server $node.server -Name $metadata.infraID -Location $dc
            }
            catch {
                throw "Get-Folder failed for VM '$name': Could not find folder '$($metadata.infraID)' in datacenter '$($node.datacenter)' on server '$($node.server)'. Ensure the folder was created during failure domain setup. $nodeContext Original error: $($_.Exception.Message)"
            }

            try {
                $template = Get-VM -Server $node.server -Name $vm_template -Location $dc
            }
            catch {
                throw "Get-VM failed for VM '$name': Could not find VM template '$vm_template' in datacenter '$($node.datacenter)' on server '$($node.server)'. Ensure the OVA template was uploaded. $nodeContext Original error: $($_.Exception.Message)"
            }

            # Clone the virtual machine from the imported template
            try {
                $vm = New-OpenShiftVM -Server $node.server -Template $template -Name $name -ResourcePool $rp -Datastore $datastoreInfo -Location $folder -IgnitionData $ignition -Tag $tag -Networking $network -Network $node.network -SecureBoot $secureboot -StoragePolicy $storagepolicy -NumCPU $numCPU -MemoryMB $memory -CoresPerSocket $coresPerSocket
            }
            catch {
                throw "New-OpenShiftVM failed for VM '$name'. $nodeContext Original error: $($_.Exception.Message)"
            }

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

# This function is used to set secure boot.
function Set-SecureBoot {
    param(
        $VM
    )

    $spec = New-Object VMware.Vim.VirtualMachineConfigSpec
    $spec.Firmware = [VMware.Vim.GuestOsDescriptorFirmwareType]::efi

    $boot = New-Object VMware.Vim.VirtualMachineBootOptions
    $boot.EfiSecureBootEnabled = $true

    $spec.BootOptions = $boot

    $VM.ExtensionData.ReconfigVM($spec)
}
