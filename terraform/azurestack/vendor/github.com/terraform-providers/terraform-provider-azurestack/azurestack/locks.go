package azurestack

// handle the case of using the same name for different kinds of resources
func azureStackLockByName(name string, resourceType string) {
	updatedName := resourceType + "." + name
	armMutexKV.Lock(updatedName)
}

func azureStackLockMultipleByName(names *[]string, resourceType string) {
	for _, name := range *names {
		azureStackLockByName(name, resourceType)
	}
}

func azureStackUnlockByName(name string, resourceType string) {
	updatedName := resourceType + "." + name
	armMutexKV.Unlock(updatedName)
}

func azureStackUnlockMultipleByName(names *[]string, resourceType string) {
	for _, name := range *names {
		azureStackUnlockByName(name, resourceType)
	}
}
