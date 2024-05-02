#!/usr/bin/pwsh

. .\variables.ps1

$ErrorActionPreference = "Stop"

# since we do not have ca for vsphere certs, we'll just set insecure
Set-PowerCLIConfiguration -InvalidCertificateAction:Ignore -Confirm:$false | Out-Null
$Env:GOVC_INSECURE = 1

# Connect to vCenter
Connect-VIServer -Server $vcenter -Credential (Import-Clixml $vcentercredpath)

# Convert the installer metadata to a powershell object
$metadata = Get-Content -Path ./metadata.json | ConvertFrom-Json

# Get tag for all resources we created
$tagCategory = Get-TagCategory -Name "openshift-$($metadata.infraID)"
$tag = Get-Tag -Category $tagCategory -Name "$($metadata.infraID)"

# Clean up all VMs
$vms = Get-VM -Tag $tag
foreach ($vm in $vms) {
    if ($vm.PowerState -eq "PoweredOn") {
        Write-Output "Stopping VM $vm"
        Stop-VM -VM $vm -confirm:$false > $null
    }
    Write-Output "Removing VM $vm"
    Remove-VM -VM $vm -DeletePermanently -confirm:$false
}

# Clean up all templates
$templates = Get-TagAssignment -Tag $tag -Entity (Get-Template)
foreach ($template in $templates) {
    Write-Output "Removing template $($template.Entity)"
    Remove-Template -Template $($template.Entity) -DeletePermanently -confirm:$false
}

# Clean up all resource pools
$rps = Get-TagAssignment -Tag $tag -Entity (Get-ResourcePool)
foreach ($rp in $rps) {
    Write-Output "Removing resource pool $($rp.Entity)"
    Remove-ResourcePool -ResourcePool $($rp.Entity) -confirm:$false
}

# Clean up all folders
$folders = Get-TagAssignment -Tag $tag -Entity (Get-Folder)
foreach ($folder in $folders) {
    Write-Output "Removing folder $($folder.Entity)"
    Remove-Folder -Folder $($folder.Entity) -DeletePermanently -confirm:$false
}

# Clean up storage policy.  Must be done after all other object cleanup except tag/tagCategory
$storagePolicies = Get-SpbmStoragePolicy -Tag $tag

foreach ($policy in $storagePolicies) {

    $clusterInventory = @()
    $splitResults = @($policy.Name -split "openshift-storage-policy-")

    if ($splitResults.Count -eq 2) {
        $clusterId = $splitResults[1]
        if ($clusterId -ne "") {
            Write-Host "Checking for storage policies for "$clusterId
            $clusterInventory = @(Get-Inventory -Name "$($clusterId)*" -ErrorAction Continue)

            if ($clusterInventory.Count -eq 0) {
                Write-Host "Removing policy: $($policy.Name)"
                $policy | Remove-SpbmStoragePolicy -Confirm:$false
            }
            else {
                Write-Host "not deleting: $($clusterInventory)"
            }
        }
    }
}

# Clean up tags
Remove-Tag -Tag $tag -confirm:$false
Remove-TagCategory -Category $tagCategory -confirm:$false

Disconnect-VIServer -Server $vcenter -Force:$true -Confirm:$false