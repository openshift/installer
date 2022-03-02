package vsphere

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/folder"
)

// resourceVSphereFolderMigrateState is the main state migration function for
// the vsphere_folder resource.
func resourceVSphereFolderMigrateState(version int, os *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	// Guard against a nil state.
	if os == nil {
		return nil, nil
	}

	// Guard against empty state, can't do anything with it
	if os.Empty() {
		return os, nil
	}

	var migrateFunc func(*terraform.InstanceState, interface{}) error
	switch version {
	case 0:
		log.Printf("[DEBUG] Migrating vsphere_folder state: old v%d state: %#v", version, os)
		migrateFunc = resourceVSphereFolderMigrateStateV1
	default:
		// Migration is complete
		log.Printf("[DEBUG] Migrating vsphere_folder state: completed v%d state: %#v", version, os)
		return os, nil
	}
	if err := migrateFunc(os, meta); err != nil {
		return nil, err
	}
	version++
	log.Printf("[DEBUG] Migrating vsphere_folder state: new v%d state: %#v", version, os)
	return resourceVSphereFolderMigrateState(version, os, meta)
}

// resourceVSphereFolderMigrateStateV1 migrates the state of the vsphere_folder
// from version 0 to version 1.
func resourceVSphereFolderMigrateStateV1(s *terraform.InstanceState, meta interface{}) error {
	// Our path for migration here is pretty much the same as our import path, so
	// we just leverage that functionality.
	//
	// We just need the path and the datacenter to proceed. We don't have an
	// analog in for existing_path in the new resource, so we just drop that on
	// the floor.
	dcp := folder.NormalizePath(s.Attributes["datacenter"])
	p := folder.NormalizePath(s.Attributes["path"])

	// Discover our datacenter first. This field can be empty, so we have to
	// search for it as we normally would.
	client := meta.(*VSphereClient).vimClient
	dc, err := getDatacenter(client, dcp)
	if err != nil {
		return err
	}

	// The old resource only supported VM folders, so this part is easy enough,
	// we can derive our full path by combining the VM path particle and our
	// relative path.
	fp := folder.RootPathParticleVM.PathFromDatacenter(dc, p)
	folder, err := folder.FromAbsolutePath(client, fp)
	if err != nil {
		return err
	}

	// We got our folder!
	//
	// Read will handle everything except for the ID, so just wipe attributes,
	// update the ID, and let read take care of the rest.
	s.Attributes = make(map[string]string)
	s.ID = folder.Reference().Value

	return nil
}
