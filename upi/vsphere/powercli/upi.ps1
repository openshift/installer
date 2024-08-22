#!/usr/bin/pwsh

$MYINV = $MyInvocation
$SCRIPTDIR = split-path $MYINV.MyCommand.Path

Write-Output "SCRIPT DIR: $($SCRIPTDIR)"

. .\variables.ps1
. ${SCRIPTDIR}\upi-functions.ps1

$ErrorActionPreference = "Stop"

# since we do not have ca for vsphere certs, we'll just set insecure
# we will also set default vi server mode to multiple for multi vcenter support
Set-PowerCLIConfiguration -InvalidCertificateAction:Ignore -DefaultVIServerMode Multiple -ParticipateInCEIP:$false -Confirm:$false | Out-Null
$Env:GOVC_INSECURE = 1

$viservers = @{}

# Connect to vCenter
if ($null -eq $vcenters) {
    $viservers[$vcenter] = Connect-VIServer -Server $vcenter -Credential (Import-Clixml $vcentercredpath)
}
else {
    $vcenters = $vcenters | ConvertFrom-Json
    foreach ($vc in $vcenters) {
        Write-Output "Logging into $($vc.server)"
        $viservers[$($vc.server)] = Connect-VIServer -Server $vc.server -User $vc.user -Password $($vc.password)
    }
}
$cliContext = Get-PowerCLIContext

if ($downloadInstaller) {
    Write-Output "Downloading the most recent $($version) installer"

    $releaseApiUri = "https://api.github.com/repos/openshift/okd/releases"
    $progressPreference = 'silentlyContinue'
    $webrequest = Invoke-WebRequest -uri $releaseApiUri
    $progressPreference = 'Continue'
    $releases = ConvertFrom-Json $webrequest.Content -AsHashtable
    $publishedDate = (Get-Date).AddDays(-365)
    $currentRelease = $null

    foreach($r in $releases) {
        if($r['name'] -like "*$($version)*") {
            if ($publishedDate -lt $r['published_at'] ) {
                $publishedDate = $r['published_at']
                $currentRelease = $r
            }
        }
    }

    foreach($asset in $currentRelease['assets']) {
        if($asset['name'] -like "openshift-install-linux*") {
            $installerUrl = $asset['browser_download_url']
        }
    }

    # If openshift-install doesn't exist on the path, download it and extract
    if (-Not (Test-Path -Path "openshift-install")) {

        $progressPreference = 'silentlyContinue'
        Invoke-WebRequest -uri $installerUrl -OutFile "installer.tar.gz"
        tar -xvf "installer.tar.gz"
        $progressPreference = 'Continue'
    }
}

if ($uploadTemplateOva) {
    Write-Output "Checking for RHCOS OVA"

    # If the OVA doesn't exist on the path, determine the url from openshift-install and download it.
    if (-Not (Test-Path -Path "template-$($Version).ova")) {
        Write-Output "Downloading RHCOS OVA"
        Start-Process -Wait -Path ./openshift-install -ArgumentList @("coreos", "print-stream-json") -RedirectStandardOutput coreos.json

        $coreosData = Get-Content -Path ./coreos.json | ConvertFrom-Json -AsHashtable
        $ovaUri = $coreosData.architectures.x86_64.artifacts.vmware.formats.ova.disk.location
        $progressPreference = 'silentlyContinue'
        Invoke-WebRequest -uri $ovaUri -OutFile "template-$($Version).ova"
        $progressPreference = 'Continue'
    }
}

$sshKey = [string](Get-Content -Path $sshkeypath -Raw:$true) -Replace '\n',''

if ($createInstallConfig) {
    # Without having to add additional powershell modules yaml is difficult to deal
    # with. There is a supplied install-config.json which is converted to a powershell
    # object
    $config = ConvertFrom-Json -InputObject $installconfig

    # Set the install-config.json from upi-variables
    $config.metadata.name = $clustername
    $config.baseDomain = $basedomain
    $config.sshKey = $sshKey
    $config.platform.vsphere.vcenter = $vcenter
    $config.platform.vsphere.username = $username
    $config.platform.vsphere.password = $password
    $config.platform.vsphere.datacenter = $datacenter
    $config.platform.vsphere.defaultDatastore = $datastore
    $config.platform.vsphere.cluster = $cluster
    $config.platform.vsphere.network = $portgroup
    # $config.platform.vsphere.apiVIP = $apivip
    # $config.platform.vsphere.ingressVIP = $ingressvip

    $config.pullSecret = $pullsecret -replace "`n", "" -replace " ", ""

    # Write out the install-config.yaml (really json)
    $config | ConvertTo-Json -Depth 8 | Out-File -FilePath install-config.yaml -Force:$true
}

if ($generateIgnitions) {
    # openshift-install create manifests
    start-process -Wait -FilePath ./openshift-install -argumentlist @("create", "manifests")

    # Remove master machines and the worker machinesets
    rm -f openshift/99_openshift-cluster-api_master-machines-*.yaml openshift/99_openshift-cluster-api_worker-machineset-*.yaml

    # openshift-install create ignition-configs
    start-process -Wait -FilePath ./openshift-install -argumentlist @("create", "ignition-configs")
}

# Check failure domains.  If not set, create a default single failure domain from settings
if ($null -eq $failure_domains) {
    Write-Output "Generating Failure Domain..."
    $failure_domains = @"
[
  {
        "server": "$($vcenter)",
        "datacenter": "$($datacenter)",
        "cluster": "$($cluster)",
        "datastore": "$($datastore)",
        "network": "$($portgroup)"
  }
]
"@
}
$fds = $failure_domains | ConvertFrom-Json

# Convert the installer metadata to a powershell object
$metadata = Get-Content -Path ./metadata.json | ConvertFrom-Json

# Since we are using MachineSets for the workers make sure we set the
# template name to what is expected to be generated by the installer.
if ($null -eq $vm_template) {
    $vm_template = "$( $metadata.infraID )-rhcos"
}

# Create tag for all resources we create
foreach ($viserver in $viservers.Keys) {
    $tagCategory = Get-TagCategory -Server $viserver -Name "openshift-$($metadata.infraID)" -ErrorAction continue 2>$null
    if (-Not $?) {
        Write-Output "Creating Tag Category openshift-$($metadata.infraID)"
        $tagCategory = New-TagCategory -Server $viserver -Name "openshift-$($metadata.infraID)" -EntityType "urn:vim25:VirtualMachine","urn:vim25:ResourcePool","urn:vim25:Folder","urn:vim25:Datastore","urn:vim25:StoragePod"
    }
    $tag = Get-Tag -Server $viserver -Category $tagCategory -Name "$($metadata.infraID)" -ErrorAction continue 2>$null
    if (-Not $?) {
        Write-Output "Creating Tag $($metadata.infraID)"
        $tag = New-Tag -Server $viserver -Category $tagCategory -Name "$($metadata.infraID)"
    }
}

$jobs = @()
$templateInProgress = @()

# Check each failure domain for ova template
foreach ($fd in $fds)
{
    Write-Output "Getting viserver for $($fd.server)"
    $viserver = $viservers[$fd.server]
    $datastoreInfo = Get-Datastore -Server $viserver -Name $fd.datastore -Location $fd.datacenter

    # Load tags for this FD.
    $tagCategory = Get-TagCategory -Server $viserver -Name "openshift-$($metadata.infraID)" -ErrorAction continue 2>$null
    $tag = Get-Tag -Server $viserver -Category $tagCategory -Name "$($metadata.infraID)" -ErrorAction continue 2>$null

    # If the folder already exists
    Write-Output "Checking for folder in failure domain $($fd.datacenter)/$($fd.cluster)"
    $folder = Get-Folder -Server $viserver -Name $clustername -Location $fd.datacenter -ErrorAction continue 2>$null

    # Otherwise create the folder within the datacenter as defined in the upi-variables
    if (-Not $?) {
        Write-Output "Creating folder $($clustername) in datacenter $($fd.datacenter) of vCenter $($fd.server)"
        (get-view (Get-Datacenter -Server $viserver -Name $fd.datacenter).ExtensionData.vmfolder).CreateFolder($clustername)
        $folder = Get-Folder -Server $viserver -Name $clustername -Location $fd.datacenter
        New-TagAssignment -Server $viserver -Entity $folder -Tag $tag > $null
    }

    # Create resource pool for all future VMs
    Write-Output "Checking for resource pool in failure domain $($fd.datacenter)/$($fd.cluster)"
    $rp = Get-ResourcePool -Server $viserver -Name $($metadata.infraID) -Location $(Get-Cluster -Server $viserver -Name $($fd.cluster)) -ErrorAction continue 2>$null

    if (-Not $?) {
        Write-Output "Creating resource pool $($metadata.infraID) in datacenter $($fd.datacenter)"
        $rp = New-ResourcePool -Server $viserver -Name $($metadata.infraID) -Location $(Get-Cluster -Server $viserver -Name $($fd.cluster))
        New-TagAssignment -Server $viserver -Entity $rp -Tag $tag > $null
    }

    # If the rhcos virtual machine already exists
    Write-Output "Checking for vm template in failure domain $($fd.datacenter)/$($fd.cluster)"
    $template = Get-VM -Server $viserver -Name $vm_template -Location $fd.datacenter -ErrorAction continue 2>$null

    # Otherwise import the ova to a random host on the vSphere cluster
    if (-Not $? -And -Not $templateInProgress.Contains($fd.datacenter))
    {
        $templateInProgress += $fd.datacenter
        $vmhost = Get-Random -InputObject (Get-VMHost -Location (Get-Cluster $fd.cluster))
        $ovfConfig = Get-OvfConfiguration -Server $viserver -Ovf "template-$($Version).ova"
        $ovfConfig.NetworkMapping.VM_Network.Value = $fd.network
        Write-Output "OVF: $($ovfConfig)"
        $jobs += Start-ThreadJob -n "upload-template-$($fd.cluster)" -ScriptBlock {
            param($Version,$vm_template,$ovfConfig,$vmhost,$datastoreInfo,$folder,$tag,$scriptdir,$cliContext,$viserver)
            . .\variables.ps1
            . ${scriptdir}\upi-functions.ps1
            Write-Output "Version: $($Version)"
            Write-Output "VM Template: $($vm_template)"
            Write-Output "OVF Config: $($ovfConfig)"
            Write-Output "VM Host: $($vmhost)"
            Use-PowerCLIContext -PowerCLIContext $cliContext
            $template = Import-Vapp -Server $viserver -Source "template-$($Version).ova" -Name $vm_template -OvfConfiguration $ovfConfig -VMHost $vmhost -Datastore $datastoreInfo -InventoryLocation $folder -Force:$true

            $templateVIObj = Get-View -VIObject $template.Name
            # Need to look into upgrading hardware.  For me it keeps throwing exception.
            <# try {
            $templateVIObj.UpgradeVM($hardwareVersion)
        }
        catch {
            Write-Output "Something happened setting VM hardware version"
            Write-Output $_
        } #>

            New-TagAssignment -Server $viserver -Entity $template -Tag $tag
            Set-VM -Server $viserver -VM $template -MemoryGB 16 -NumCpu 4 -CoresPerSocket 4 -Confirm:$false > $null
            #Get-HardDisk -VM $template | Select-Object -First 1 | Set-HardDisk -CapacityGB 120 -Confirm:$false > $null
            updateDisk -VM $template -CapacityGB 120
            New-AdvancedSetting -Server $viserver -Entity $template -name "disk.EnableUUID" -value 'TRUE' -confirm:$false -Force > $null
            New-AdvancedSetting -Server $viserver -Entity $template -name "guestinfo.ignition.config.data.encoding" -value "base64" -confirm:$false -Force > $null
            #$snapshot = New-Snapshot -VM $template -Name "linked-clone" -Description "linked-clone" -Memory -Quiesce
        } -ArgumentList @($Version,$vm_template,$ovfConfig,$vmhost,$datastoreInfo,$folder,$tag,$SCRIPTDIR,$cliContext,$viserver)
    }
}

# If jobs were started, lets wait till they are done
if ($jobs.count -gt 0)
{
    Wait-Job -Job $jobs
    foreach ($job in $jobs) {
        Receive-Job -Job $job
    }
}

Write-Output "Creating LB"

# Data needed for LB VM creation
$tagCategory = Get-TagCategory -Server $fds[0].server -Name "openshift-$($metadata.infraID)" -ErrorAction continue 2>$null
$tag = Get-Tag -Server $fds[0].server -Category $tagCategory -Name "$($metadata.infraID)" -ErrorAction continue 2>$null
$rp = Get-ResourcePool -Name $($metadata.infraID) -Location $(Get-Cluster -Server $fds[0].server -Name $($fds[0].cluster))
$datastoreInfo = Get-Datastore -Name $fds[0].datastore -Server $fds[0].server -Location $fds[0].datacenter
$folder = Get-Folder -Server $fds[0].server -Name $clustername -Location $fds[0].datacenter
$template = Get-VM -Server $fds[0].server -Name $vm_template -Location $fds[0].datacenter

# Create LB for Cluster
$ignition = [Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes((New-LoadBalancerIgnition $sshKey)))
$network = New-VMNetworkConfig -Server $fds[0].server -Hostname "$($metadata.infraID)-lb" -IPAddress $lb_ip_address -Netmask $netmask -Gateway $gateway -DNS $dns -Network $failure_domains[0].network
$vm = New-OpenShiftVM -IgnitionData $ignition -Name "$($metadata.infraID)-lb" -Template $template -Server $fds[0].server -ResourcePool $rp -Datastore $datastoreInfo -Location $folder -Tag $tag -Networking $network -Network $($fds[0].network) -SecureBoot $secureboot -StoragePolicy $storagepolicy -MemoryMB 8192 -NumCpu 4
$vm | Start-VM

# Take the $virtualmachines defined in upi-variables and convert to a powershell object
if ($null -eq $virtualmachines)
{
    $virtualmachines = New-VMConfigs
}
$vmHash = ConvertFrom-Json -InputObject $virtualmachines -AsHashtable

Write-Progress -id 222 -Activity "Creating virtual machines" -PercentComplete 0

New-OpenshiftVMs "bootstrap"
New-OpenshiftVMs "master"
New-OpenshiftVMs "worker"
Write-Progress -id 222 -Activity "Completed virtual machines" -PercentComplete 100 -Completed

## This is nice to have to clear screen when doing things manually.  Maybe i'll
# make this configurable.
# Clear-Host

# Instead of restarting openshift-install to wait for bootstrap, monitor
# the bootstrap configmap in the kube-system namespace

# Extract the Client Certificate Data from auth/kubeconfig
$match = Select-String "client-certificate-data: (.*)" -Path ./auth/kubeconfig
[Byte[]]$bytes = [Convert]::FromBase64String($match.Matches.Groups[1].Value)
$clientCertData = [System.Text.Encoding]::ASCII.GetString($bytes)

# Extract the Client Key Data from auth/kubeconfig
$match = Select-String "client-key-data: (.*)" -Path ./auth/kubeconfig
$bytes = [Convert]::FromBase64String($match.Matches.Groups[1].Value)
$clientKeyData = [System.Text.Encoding]::ASCII.GetString($bytes)

# Create a X509Certificate2 object for Invoke-WebRequest
$cert = [System.Security.Cryptography.X509Certificates.X509Certificate2]::CreateFromPem($clientCertData, $clientKeyData)

# Extract the kubernetes endpoint uri
$match = Select-String "server: (.*)" -Path ./auth/kubeconfig
$kubeurl = $match.Matches.Groups[1].Value

if ($waitForComplete)
{
    $apiTimeout = (20*60)
    $apiCount = 1
    $apiSleep = 30
    Write-Progress -Id 444 -Status "1% Complete" -Activity "API" -PercentComplete 1
    :api while ($true) {
        Start-Sleep -Seconds $apiSleep
        try {
            $webrequest = Invoke-WebRequest -Uri "$($kubeurl)/version" -SkipCertificateCheck
            $version = (ConvertFrom-Json $webrequest.Content).gitVersion

            if ($version -ne "" ) {
                Write-Debug "API Version: $($version)"
                Write-Progress -Id 444 -Status "Completed" -Activity "API" -PercentComplete 100
                break api
            }
        }
        catch {}

        $percentage = ((($apiCount*$apiSleep)/$apiTimeout)*100)
        if ($percentage -le 100) {
            Write-Progress -Id 444 -Status "$percentage% Complete" -Activity "API" -PercentComplete $percentage
        }
        $apiCount++
    }

    $bootstrapTimeout = (30*60)
    $bootstrapCount = 1
    $bootstrapSleep = 30
    Write-Progress -Id 333 -Status "1% Complete" -Activity "Bootstrap" -PercentComplete 1
    :bootstrap while ($true)
    {
        Start-Sleep -Seconds $bootstrapSleep

        try
        {
            $webrequest = Invoke-WebRequest -Certificate $cert -Uri "$( $kubeurl )/api/v1/namespaces/kube-system/configmaps/bootstrap" -SkipCertificateCheck

            $bootstrapStatus = (ConvertFrom-Json $webrequest.Content).data.status

            if ($bootstrapStatus -eq "complete")
            {
                Get-VM "$( $metadata.infraID )-bootstrap" | Stop-VM -Confirm:$false | Remove-VM -DeletePermanently -Confirm:$false
                Write-Progress -Id 333 -Status "Completed" -Activity "Bootstrap" -PercentComplete 100
                break bootstrap
            }
        }
        catch
        {
        }

        $percentage = ((($bootstrapCount*$bootstrapSleep)/$bootstrapTimeout)*100)
        if ($percentage -le 100)
        {
            Write-Progress -Id 333 -Status "$percentage% Complete" -Activity "Bootstrap" -PercentComplete $percentage
        }
        else
        {
            Write-Output "Warning: Bootstrap taking longer than usual." -NoNewLine -ForegroundColor Yellow
        }

        $bootstrapCount++
    }

    # Now that bootstrap is complete, we should be getting worker node CSRs that need to be approved before being
    # able to finish installation.  Lets monitor for CSRs, approve them and verify the number of worker nodes have
    # now appeared and are Ready before moving on.

    # [ngirard@fedora ibm7-installs]$ oc get csr | grep Pending
    #csr-2hgbd                                        2m52s   kubernetes.io/kube-apiserver-client-kubelet   system:serviceaccount:openshift-machine-config-operator:node-bootstrapper         <none>              Pending
    #csr-lwmgf                                        2m19s   kubernetes.io/kube-apiserver-client-kubelet   system:serviceaccount:openshift-machine-config-operator:node-bootstrapper         <none>              Pending
    #csr-scvk6                                        2m30s   kubernetes.io/kube-apiserver-client-kubelet   system:serviceaccount:openshift-machine-config-operator:node-bootstrapper         <none>              Pending

    # apis/certificates.k8s.io/v1/certificatesigningrequests
    $csrTimeout = (600/5)
    $csrCount = 1
    $csrSleep = 5
    Write-Progress -Id 222 -Status "1% Complete" -Activity "Worker Ready" -PercentComplete 0
    :csrLoop while ($true)
    {
        Start-Sleep -Seconds $csrSleep

        try
        {
            $webrequest = Invoke-WebRequest -Certificate $cert -Uri "$( $kubeurl )/apis/certificates.k8s.io/v1/certificatesigningrequests" -SkipCertificateCheck

            $csrs = (ConvertFrom-Json $webrequest.Content).items

            foreach ($csr in $csrs)
            {
                # Check if no status (Pending) and if its a type we are looking for (kubernetes.io/kubelet-serving) (kubernetes.io/kube-apiserver-client-kubelet)
                $bootstrapper = "system:serviceaccount:openshift-machine-config-operator:node-bootstrapper"
                $nodeUser = "system:node:"

                $csrUser = $csr.spec.username
                if ($csr.status.conditions -eq $null -And ($csrUser -eq $bootstrapper -Or $csrUser.Contains($nodeUser)))
                {
                    $conditions = New-Object System.Collections.ArrayList
                    $condition = @{
                        type = "Approved"
                        status = "True"
                        reason = "PowershellApproved"
                        message = "This CSR was approved by script in PowerShell."
                    }
                    $conditions.Add($condition) > $null
                    $csr.status | add-member -Name "conditions" -value $conditions  -MemberType NoteProperty
                    Write-Output "Accepting CSR: $( $csr.metadata.name )"
                    $csrResponse = Invoke-RestMethod -Method "Put" -Certificate $cert -Uri "$( $kubeurl )/apis/certificates.k8s.io/v1/certificatesigningrequests/$( $csr.metadata.name )/approval" -SkipCertificateCheck -Body (ConvertTo-Json $csr -Depth 6)
                }
            }
        }
        catch
        {
            #Write-Output $_
        }

        # Check number of worker nodes with NotReady/Ready.  NotReady will be 1 pt where Ready will be 2.
        $currentComputePoints = 0

        try
        {
            $webrequest = Invoke-WebRequest -Certificate $cert -Uri "$( $kubeurl )/api/v1/nodes" -SkipCertificateCheck

            $nodes = (ConvertFrom-Json $webrequest.Content).items

            foreach ($node in $nodes)
            {
                if ($node.metadata.labels.psobject.properties.name -Contains "node-role.kubernetes.io/worker")
                {
                    #Write-Output "Checking node $($node.metadata.name)"
                    foreach ($condition in $node.status.conditions)
                    {
                        if ($condition.type -eq "Ready")
                        {
                            #Write-Output "Is node ready? $($condition.status)"
                            if ($condition.status -eq "True")
                            {
                                $currentComputePoints = $currentComputePoints + 2
                            }
                            else
                            {
                                $currentComputePoints++
                            }
                        }
                    }
                }
            }
        }
        catch
        {
            #Write-Output $_
        }

        $maxComputePoints = $compute_count * 2
        $percentage = ((($currentComputePoints)/$maxComputePoints)*100)
        if ($percentage -eq 100)
        {
            Write-Progress -Id 222 -Status "Completed" -Activity "Worker Ready" -PercentComplete 100
            break csrLoop
        }
        elseif ($percentage -le 100)
        {
            Write-Progress -Id 222 -Status "$percentage% Complete" -Activity "Worker Ready" -PercentComplete $percentage
        }

        if ($csrCount -ge $csrTimeout)
        {
            Write-Output "Warning: Bootstrap taking longer than usual." -NoNewLine -ForegroundColor Yellow
            break csrLoop
        }

        $csrCount++
    }

    $progressMsg = ""
    Write-Progress -Id 111 -Status "1% Complete" -Activity "Install" -PercentComplete 1
    :installcomplete while ($true)
    {
        Start-Sleep -Seconds 30
        try
        {
            $webrequest = Invoke-WebRequest -Certificate $cert -Uri "$( $kubeurl )/apis/config.openshift.io/v1/clusterversions" -SkipCertificateCheck

            $clusterversions = ConvertFrom-Json $webrequest.Content -AsHashtable

            # just like the installer check the status conditions of the clusterversions config
            foreach ($condition in $clusterversions['items'][0]['status']['conditions'])
            {
                switch ($condition['type'])
                {
                    "Progressing" {
                        if ($condition['status'] -eq "True")
                        {

                            $matchper = ($condition['message'] | Select-String "^Working.*\(([0-9]{1,3})\%.*\)")
                            $matchmsg = ($condition['message'] | Select-String -AllMatches -Pattern "^(Working.*)\:.*")

                            # During install, the status of CVO will / may go degraded due to operators going
                            # degraded from taking a while to install.  It seems this is the new norm as control
                            # plane takes a while to roll out and certain operators go degraded until the control
                            # plane is stable.
                            if ($matchmsg.Matches.Groups -ne $null)
                            {
                                $progressMsg = $matchmsg.Matches.Groups[1].Value
                                $progressPercent = $matchper.Matches.Groups[1].Value

                                Write-Progress -Id 111 -Status "$progressPercent% Complete - $( $progressMsg )" -Activity "Install" -PercentComplete $progressPercent
                            }
                            continue
                        }
                    }
                    "Available" {
                        if ($condition['status'] -eq "True")
                        {
                            Write-Progress -Id 111 -Activity "Install" -Status "Completed" -PercentComplete 100
                            break installcomplete
                        }
                        continue
                    }
                    Default {
                        continue
                    }
                }
            }
        }
        catch
        {
            Write-Output "Unable to check operators"
            Write-Output $_
        }
    }
}

Get-Job | Remove-Job

foreach ($key in $viservers.Keys) {
    Write-Output "Disconnecting from $($key)"
    Disconnect-VIServer -Server $key -Force:$true -Confirm:$false
}

Write-Output "Install Complete!"
