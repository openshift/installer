package contentlibrary

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/datastore"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/ovfdeploy"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/provider"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/ovf"
	"github.com/vmware/govmomi/vapi/library"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/vcenter"
	"github.com/vmware/govmomi/vim25/soap"
)

// FromName accepts a Content Library name and returns a Library object.
func FromName(c *rest.Client, name string) (*library.Library, error) {
	log.Printf("[DEBUG] contentlibrary.FromName: Retrieving content library %s by name", name)
	clm := library.NewManager(c)
	ctx := context.TODO()
	lib, err := clm.GetLibraryByName(ctx, name)
	if err != nil {
		return nil, provider.Error(name, "FromName", err)
	}
	if lib == nil {
		return nil, provider.Error(name, "FromName", fmt.Errorf("Unable to find content library (%s)", name))
	}
	log.Printf("[DEBUG] contentlibrary.FromName: Successfully retrieved content library %s", name)
	return lib, nil
}

// FromID accepts a Content Library ID and returns a Library object.
func FromID(c *rest.Client, id string) (*library.Library, error) {
	log.Printf("[DEBUG] contentlibrary.FromID: Retrieving content library %s by ID", id)
	clm := library.NewManager(c)
	ctx := context.TODO()
	lib, err := clm.GetLibraryByID(ctx, id)
	if err != nil {
		return nil, provider.Error(id, "FromID", err)
	}
	if lib == nil {
		return nil, fmt.Errorf("Unable to find content library (%s)", id)
	}
	log.Printf("[DEBUG] contentlibrary.FromID: Successfully retrieved content library %s", id)
	return lib, nil
}

// CreateLibrary creates a Content Library.
func CreateLibrary(d *schema.ResourceData, restclient *rest.Client, backings []library.StorageBackings) (string, error) {
	name := d.Get("name").(string)
	log.Printf("[DEBUG] contentlibrary.CreateLibrary: Creating content library %s", name)
	clm := library.NewManager(restclient)
	ctx := context.TODO()
	lib := library.Library{
		Description: d.Get("description").(string),
		Name:        name,
		Storage:     backings,
		Type:        "LOCAL",
	}
	if len(d.Get("publication").([]interface{})) > 0 {
		publication := d.Get("publication").([]interface{})[0].(map[string]interface{})
		lib.Publication = &library.Publication{
			Published:            structure.BoolPtr(publication["published"].(bool)),
			AuthenticationMethod: publication["authentication_method"].(string),
			UserName:             publication["username"].(string),
			Password:             publication["password"].(string),
		}
	}
	if len(d.Get("subscription").([]interface{})) > 0 {
		subscription := d.Get("subscription").([]interface{})[0].(map[string]interface{})
		lib.Subscription = &library.Subscription{
			AutomaticSyncEnabled: structure.BoolPtr(subscription["automatic_sync"].(bool)),
			OnDemand:             structure.BoolPtr(subscription["on_demand"].(bool)),
			AuthenticationMethod: subscription["authentication_method"].(string),
			UserName:             subscription["username"].(string),
			Password:             subscription["password"].(string),
			SubscriptionURL:      subscription["subscription_url"].(string),
		}
		lib.Type = "SUBSCRIBED"
	}
	id, err := clm.CreateLibrary(ctx, lib)
	if err != nil {
		return "", provider.Error(name, "CreateLibrary", err)
	}
	log.Printf("[DEBUG] contentlibrary.CreateLibrary: Content library %s successfully created", name)
	return id, nil
}

// DeleteLibrary deletes a Content Library.
func DeleteLibrary(c *rest.Client, lib *library.Library) error {
	log.Printf("[DEBUG] contentlibrary.DeleteLibrary: Deleting library %s", lib.Name)
	clm := library.NewManager(c)
	ctx := context.TODO()
	err := clm.DeleteLibrary(ctx, lib)
	if err != nil {
		return provider.Error(lib.ID, "DeleteLibrary", err)
	}
	log.Printf("[DEBUG] contentlibrary.DeleteLibrary: Deleting library %s", lib.Name)
	return nil
}

// ItemFromName accepts a Content Library item name along with a Content Library and will return the item object.
func ItemFromName(c *rest.Client, l *library.Library, name string) (*library.Item, error) {
	log.Printf("[DEBUG] contentlibrary.ItemFromName: Retrieving library item %s.", name)
	clm := library.NewManager(c)
	ctx := context.TODO()
	fi := library.FindItem{
		LibraryID: l.ID,
		Name:      name,
	}
	items, err := clm.FindLibraryItems(ctx, fi)
	if err != nil {
		return nil, nil
	}
	if len(items) < 1 {
		return nil, fmt.Errorf("Unable to find content library item (%s)", name)
	}
	item, err := clm.GetLibraryItem(ctx, items[0])
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] contentlibrary.ItemFromName: Library item %s retrieved successfully", name)
	return item, nil
}

// ItemFromID accepts a Content Library item ID and will return the item object.
func ItemFromID(c *rest.Client, id string) (*library.Item, error) {
	log.Printf("[DEBUG] contentlibrary.ItemFromID: Retrieving library item %s", id)
	clm := library.NewManager(c)
	ctx := context.TODO()
	item, err := clm.GetLibraryItem(ctx, id)
	if err != nil {
		return nil, provider.Error(id, "ItemFromID", err)
	}
	log.Printf("[DEBUG] contentlibrary.ItemFromID: Library item %s retrieved successfully", id)
	return item, nil
}

// IsContentLibraryItem accepts an ID and determines if that ID is associated with an item in a Content Library.
func IsContentLibraryItem(c *rest.Client, id string) bool {
	log.Printf("[DEBUG] contentlibrary.IsContentLibrary: Checking if %s is a content library source", id)
	item, _ := ItemFromID(c, id)
	return item != nil
}

// CreateLibraryItem creates an item in a Content Library.
func CreateLibraryItem(c *rest.Client, l *library.Library, name string, desc string, t string, file string, moid string) (*string, error) {
	log.Printf("[DEBUG] contentlibrary.CreateLibraryItem: Creating content library item %s.", name)
	clm := library.NewManager(c)
	ctx := context.TODO()
	item := library.Item{
		Description: desc,
		LibraryID:   l.ID,
		Name:        name,
		Type:        t,
	}
	uploadSession := libraryUploadSession{
		ContentLibraryManager: clm,
		RestClient:            c,
		LibraryID:             l.ID,
	}
	if moid != "" {
		return uploadSession.cloneTemplate(moid, name, t)
	}

	id, err := clm.CreateLibraryItem(ctx, item)
	if err != nil {
		return nil, provider.Error(name, "CreateLibraryItem", err)
	}
	session, err := clm.CreateLibraryItemUpdateSession(ctx, library.Session{LibraryItemID: id})
	if err != nil {
		return nil, provider.Error(name, "CreateLibraryItem", err)
	}
	uploadSession.UploadSession = session
	defer func() {
		_ = clm.CompleteLibraryItemUpdateSession(ctx, session)
	}()

	isOva := false
	isLocal := true

	if strings.HasPrefix(file, "http") {
		isLocal = false
	}
	if strings.HasSuffix(file, ".ova") {
		isOva = true
	}

	ovfDescriptor, err := ovfdeploy.GetOvfDescriptor(file, isOva, isLocal, true)
	if err != nil {
		return nil, provider.Error(name, "CreateLibraryItem", err)
	}

	switch {
	case isLocal && isOva:
		return &id, uploadSession.deployLocalOva(file, ovfDescriptor)
	case isLocal && !isOva:
		return &id, uploadSession.deployLocalOvf(file, ovfDescriptor)
	case !isLocal && isOva:
		return &id, uploadSession.deployRemoteOva(file, ovfDescriptor)
	case !isLocal && !isOva:
		return &id, uploadSession.deployRemoteOvf(file)
	}

	log.Printf("[DEBUG] contentlibrary.CreateLibraryItem: Successfully created content library item %s.", name)
	return &id, nil
}

func (uploadSession *libraryUploadSession) deployRemoteOvf(file string) error {
	ctx := context.TODO()
	_, err := uploadSession.ContentLibraryManager.AddLibraryItemFileFromURI(ctx, uploadSession.UploadSession, filepath.Base(file), file)
	if err != nil {
		return err
	}
	return uploadSession.ContentLibraryManager.WaitOnLibraryItemUpdateSession(ctx, uploadSession.UploadSession, time.Second*10, func() { log.Printf("Waiting...") })
}

func (uploadSession *libraryUploadSession) deployRemoteOva(file string, ovfDescriptor string) error {
	e, err := readEnvelope(ovfDescriptor)
	if err != nil {
		return fmt.Errorf("failed to parse ovf: %s", err)
	}
	name := strings.TrimSuffix(filepath.Base(file), "ova")
	if err := uploadSession.uploadString(ovfDescriptor, name+"ovf"); err != nil {
		return err
	}
	for _, disk := range e.References {
		if err := uploadSession.uploadOvaDisksFromURL(file, disk.Href, int64(disk.Size)); err != nil {
			return err
		}
	}
	return nil
}

func (uploadSession *libraryUploadSession) deployLocalOvf(file string, ovfDescriptor string) error {
	e, err := readEnvelope(ovfDescriptor)
	if err != nil {
		return fmt.Errorf("failed to parse ovf: %s", err)
	}
	if err := uploadSession.uploadLocalFile(file); err != nil {
		return err
	}
	dir := filepath.Dir(file)
	for i := range e.References {
		if err := uploadSession.uploadLocalFile(dir + "/" + e.References[i].Href); err != nil {
			return err
		}
	}
	return nil
}

func (uploadSession *libraryUploadSession) deployLocalOva(file string, ovfDescriptor string) error {
	e, err := readEnvelope(ovfDescriptor)
	if err != nil {
		return fmt.Errorf("failed to parse ovf: %s", err)
	}
	name := strings.TrimSuffix(filepath.Base(file), "ova")
	if err := uploadSession.uploadString(ovfDescriptor, name+"ovf"); err != nil {
		return err
	}
	return uploadSession.uploadOvaDisksFromLocal(file, e)
}

type libraryUploadSession struct {
	ContentLibraryManager *library.Manager
	RestClient            *rest.Client
	UploadSession         string
	LibraryID             string
}

func (uploadSession libraryUploadSession) cloneTemplate(moid string, name string, templateType string) (*string, error) {
	ctx := context.TODO()
	if templateType == "ovf" {
		ovfItem := vcenter.OVF{
			Spec: vcenter.CreateSpec{
				Name: name,
			},
			Source: vcenter.ResourceID{
				Value: moid,
			},
			Target: vcenter.LibraryTarget{
				LibraryID: uploadSession.LibraryID,
			},
		}
		id, err := vcenter.NewManager(uploadSession.RestClient).CreateOVF(ctx, ovfItem)
		if err != nil {
			return nil, err
		}
		return &id, nil
	}
	return nil, fmt.Errorf("Unsupported template type. Only ovf can be used when cloning from vCenter")
}

func (uploadSession libraryUploadSession) uploadString(data string, name string) error {
	stringReader := strings.NewReader(data)
	openFile := io.Reader(stringReader)
	size := int64(len([]byte(data)))
	return uploadSession.upload(name, &openFile, size)
}

func (uploadSession libraryUploadSession) uploadLocalFile(file string) error {
	openFile, size, err := openLocalFile(file)
	if err != nil {
		return err
	}

	return uploadSession.upload(filepath.Base(file), openFile, *size)
}

func openLocalFile(file string) (*io.Reader, *int64, error) {
	openFile, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	statFile, err := openFile.Stat()
	if err != nil {
		return nil, nil, err
	}
	openFileReader := io.Reader(openFile)
	size := statFile.Size()
	return &openFileReader, &size, nil
}

func (uploadSession libraryUploadSession) uploadOvaDisksFromLocal(ovaFilePath string, envelope *ovf.Envelope) error {
	ovaFile, _, err := openLocalFile(ovaFilePath)
	if err != nil {
		return err
	}

	for _, disk := range envelope.References {
		size := disk.Size
		fileName := disk.Href
		if err = uploadSession.findAndUploadDiskFromOva(*ovaFile, fileName, int64(size)); err != nil {
			return err
		}
	}
	return err
}

func (uploadSession libraryUploadSession) uploadOvaDisksFromURL(ovfFilePath string, diskName string, size int64) error {
	resp, err := http.Get(ovfFilePath)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		err = uploadSession.findAndUploadDiskFromOva(resp.Body, diskName, size)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("got status %d while getting the file from remote url %s ", resp.StatusCode, ovfFilePath)
	}
	return nil
}

func (uploadSession libraryUploadSession) findAndUploadDiskFromOva(ovaFile io.Reader, diskName string, size int64) error {
	log.Printf("[DEBUG] findAndUploadDiskFromOva: Finding %s", diskName)
	ovaReader := tar.NewReader(ovaFile)
	for {
		fileHdr, err := ovaReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if fileHdr.Name == diskName {
			log.Printf("[DEBUG] findAndUploadDiskFromOva: %s found", diskName)
			ioOvaReader := io.Reader(ovaReader)
			err = uploadSession.upload(diskName, &ioOvaReader, size)
			if err != nil {
				return fmt.Errorf("error while uploading the file %s %s", diskName, err)
			}
			return nil
		}
	}
	return fmt.Errorf("disk %s not found inside ova", diskName)
}

func readEnvelope(data string) (*ovf.Envelope, error) {
	e, err := ovf.Unmarshal(strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ovf: %s", err)
	}

	return e, nil
}

func (uploadSession libraryUploadSession) upload(name string, file *io.Reader, size int64) error {
	ctx := context.TODO()

	info := library.UpdateFile{
		Name:       name,
		SourceType: "PUSH",
		Size:       size,
	}

	update, err := uploadSession.ContentLibraryManager.AddLibraryItemFile(ctx, uploadSession.UploadSession, info)
	if err != nil {
		return err
	}

	p := soap.DefaultUpload
	p.ContentLength = size
	u, err := url.Parse(update.UploadEndpoint.URI)
	if err != nil {
		return err
	}
	return uploadSession.RestClient.Upload(ctx, *file, u, &p)
}

// DeleteLibraryItem deletes an item from a Content Library.
func DeleteLibraryItem(c *rest.Client, item *library.Item) error {
	log.Printf("[DEBUG] contentlibrary.DeleteLibraryItem: Deleting content library item %s.", item.Name)
	clm := library.NewManager(c)
	ctx := context.TODO()
	err := clm.DeleteLibraryItem(ctx, item)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] contentlibrary.DeleteLibraryItem: Successfully deleted content library item %s.", item.Name)
	return nil
}

// ExpandStorageBackings takes ResourceData, and returns a list of StorageBackings.
func ExpandStorageBackings(c *govmomi.Client, d *schema.ResourceData) ([]library.StorageBackings, error) {
	log.Printf("[DEBUG] contentlibrary.ExpandStorageBackings: Expanding OVF storage backing.")
	sb := []library.StorageBackings{}
	for _, dsID := range d.Get("storage_backing").(*schema.Set).List() {
		ds, err := datastore.FromID(c, dsID.(string))
		if err != nil {
			return nil, provider.Error(d.Id(), "ExpandStorageBackings", err)
		}
		sb = append(sb, library.StorageBackings{
			DatastoreID: ds.Reference().Value,
			Type:        "DATASTORE",
		})
	}
	log.Printf("[DEBUG] contentlibrary.ExpandStorageBackings: Successfully expanded OVF storage backing.")
	return sb, nil
}

// FlattenPublication takes a Publication sub resource and sets it in ResourceData.
func FlattenPublication(d *schema.ResourceData, publication *library.Publication) error {
	if publication == nil {
		return nil
	}
	log.Printf("[DEBUG] contentlibrary.FlattenPublication: Flattening publication.")
	flatPublication := map[string]interface{}{}
	flatPublication["authentication_method"] = publication.AuthenticationMethod
	flatPublication["username"] = publication.UserName
	flatPublication["password"] = d.Get("publication.0.password").(string)
	flatPublication["publish_url"] = publication.PublishURL
	flatPublication["published"] = publication.Published
	return d.Set("publication", []interface{}{flatPublication})
}

// FlattenSubscription takes a Subscription sub resource and sets it in ResourceData.
func FlattenSubscription(d *schema.ResourceData, subscription *library.Subscription) error {
	if subscription == nil {
		return nil
	}
	log.Printf("[DEBUG] contentlibrary.FlattenSubscription: Flattening subscription.")
	flatSubscription := map[string]interface{}{}
	flatSubscription["authentication_method"] = subscription.AuthenticationMethod
	flatSubscription["username"] = subscription.UserName
	flatSubscription["password"] = d.Get("subscription.0.password").(string)
	flatSubscription["subscription_url"] = subscription.SubscriptionURL
	flatSubscription["automatic_sync"] = subscription.AutomaticSyncEnabled
	flatSubscription["on_demand"] = subscription.OnDemand
	return d.Set("subscription", []interface{}{flatSubscription})
}

// FlattenStorageBackings takes a list of StorageBackings, and returns a list of datastore IDs.
func FlattenStorageBackings(d *schema.ResourceData, sb []library.StorageBackings) error {
	log.Printf("[DEBUG] contentlibrary.FlattenStorageBackings: Flattening OVF storage backing.")
	sbl := []string{}
	for _, backing := range sb {
		if backing.Type == "DATASTORE" {
			sbl = append(sbl, backing.DatastoreID)
		}
	}
	log.Printf("[DEBUG] contentlibrary.FlattenStorageBackings: Successfully flattened OVF storage backing.")
	return d.Set("storage_backing", sbl)
}

// MapNetworkDevices maps NICs defined in the OVF to networks..
func MapNetworkDevices(d *schema.ResourceData) []vcenter.NetworkMapping {
	nm := []vcenter.NetworkMapping{}
	nics := d.Get("network_interface").([]interface{})
	for _, di := range nics {
		dm := di.(map[string]interface{})["ovf_mapping"].(string)
		dd := di.(map[string]interface{})["network_id"].(string)
		nm = append(nm, vcenter.NetworkMapping{Key: dm, Value: dd})
	}
	return nm
}
