package softwaredefinedstorage

const (
	app                 = "app"
	monLabel            = "rook-ceph-mon"
	osdLabel            = "rook-ceph-osd"
	crashcollectorLabel = "rook-ceph-crashcollector"
	monId               = "mon"
	osdId               = "ceph-osd-id"
	crashcollectorId    = "node_name"
	odfNamespace        = "openshift-storage"
)

var AppLabelId = map[string]string{
	monLabel:            monId,
	osdLabel:            osdId,
	crashcollectorLabel: crashcollectorId,
}
