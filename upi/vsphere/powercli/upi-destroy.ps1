#!/usr/bin/pwsh

. .\variables.ps1

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

foreach ($vc in $vcenters) {
    $vcenter = $($vc.server)
    Write-Output "Processing vCenter $($vcenter)"
    # Convert the installer metadata to a powershell object
    $metadata = Get-Content -Path ./metadata.json | ConvertFrom-Json

    # Get tag for all resources we created
    $tagCategory = Get-TagCategory -Server $vcenter -Name "openshift-$($metadata.infraID)"
    $tag = Get-Tag -Server $vcenter -Category $tagCategory -Name "$($metadata.infraID)"

    # Clean up all VMs
    $vms = Get-VM -Tag $tag
    foreach ($vm in $vms) {
        if ($vm.PowerState -eq "PoweredOn") {
            Write-Output "Stopping VM $vm"
            Stop-VM -Server $vcenter -VM $vm -confirm:$false > $null
        }
        Write-Output "Removing VM $vm"
        Remove-VM -Server $vcenter -VM $vm -DeletePermanently -confirm:$false
    }

    # Clean up all templates
    $templates = Get-TagAssignment -Server $vcenter -Tag $tag -Entity (Get-Template -Server $vcenter)
    foreach ($template in $templates) {
        Write-Output "Removing template $($template.Entity)"
        Remove-Template -Server $vcenter -Template $($template.Entity) -DeletePermanently -confirm:$false
    }

    # Clean up all resource pools
    $rps = Get-TagAssignment -Server $vcenter -Tag $tag -Entity (Get-ResourcePool -Server $vcenter)
    foreach ($rp in $rps) {
        Write-Output "Removing resource pool $($rp.Entity)"
        Remove-ResourcePool -Server $vcenter -ResourcePool $($rp.Entity) -confirm:$false
    }

    # Clean up all folders
    $folders = Get-TagAssignment -Server $vcenter -Tag $tag -Entity (Get-Folder -Server $vcenter)
    foreach ($folder in $folders) {
        Write-Output "Removing folder $($folder.Entity)"
        Remove-Folder -Server $vcenter -Folder $($folder.Entity) -DeletePermanently -confirm:$false
    }

    # Clean up CNS volumes
    # TODO: Add cleanup for CNS here.  This will be done in later PR.

    # Clean up storage policy.  Must be done after all other object cleanup except tag/tagCategory
    $storagePolicies = Get-SpbmStoragePolicy -Server $vcenter -Tag $tag

    foreach ($policy in $storagePolicies) {

        $clusterInventory = @()
        $splitResults = @($policy.Name -split "openshift-storage-policy-")

        if ($splitResults.Count -eq 2) {
            $clusterId = $splitResults[1]
            if ($clusterId -ne "") {
                Write-Host "Checking for storage policies for "$clusterId
                $clusterInventory = @(Get-Inventory -Server $vcenter -Name "$($clusterId)*" -ErrorAction Continue)

                if ($clusterInventory.Count -eq 0) {
                    # Remove policy from all configurations that still may exist
                    $entityConfig = Get-SpbmEntityConfiguration -StoragePolicy $policy
                    if ($null -ne $entityConfig -and $entityConfig -ne $null) {
                        Write-Host "Unsetting storage policy for "$entityConfig
                        Set-SpbmEntityConfiguration $entityConfig -StoragePolicy $null
                    }

                    # Now we can delete the policy
                    Write-Host "Removing policy: $($policy.Name)"
                    $policy | Remove-SpbmStoragePolicy -Server $vcenter -Confirm:$false -ErrorAction Continue
                }
                else {
                    Write-Host "not deleting: $($clusterInventory)"
                }
            }
        }
    }

    # Clean up tags
    Remove-Tag -Server $vcenter -Tag $tag -confirm:$false
    Remove-TagCategory -Server $vcenter -Category $tagCategory -confirm:$false

    Disconnect-VIServer -Server $vcenter -Force:$true -Confirm:$false
}
