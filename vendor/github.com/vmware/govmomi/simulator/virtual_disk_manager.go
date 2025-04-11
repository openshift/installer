/*
Copyright (c) 2017-2024 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package simulator

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/vmware/govmomi/internal"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type VirtualDiskManager struct {
	mo.VirtualDiskManager
}

func (m *VirtualDiskManager) MO() mo.VirtualDiskManager {
	return m.VirtualDiskManager
}

func VirtualDiskBackingFileName(name string) string {
	return strings.Replace(name, ".vmdk", "-flat.vmdk", 1)
}

func vdmNames(name string) []string {
	return []string{
		VirtualDiskBackingFileName(name),
		name,
	}
}

func vdmCreateVirtualDisk(op types.VirtualDeviceConfigSpecFileOperation, req *types.CreateVirtualDisk_Task) types.BaseMethodFault {
	fm := Map.FileManager()

	file, fault := fm.resolve(req.Datacenter, req.Name)
	if fault != nil {
		return fault
	}

	shouldReplace := op == types.VirtualDeviceConfigSpecFileOperationReplace
	shouldExist := op == ""
	for _, name := range vdmNames(file) {
		_, err := os.Stat(name)
		if err == nil {
			if shouldExist {
				return nil
			}
			if shouldReplace {
				if err = os.Truncate(file, 0); err != nil {
					return fm.fault(name, err, new(types.CannotCreateFile))
				}
				return nil
			}
			return fm.fault(name, nil, new(types.FileAlreadyExists))
		} else if shouldExist {
			return fm.fault(name, nil, new(types.FileNotFound))
		}

		f, err := os.Create(name)
		if err != nil {
			return fm.fault(name, err, new(types.CannotCreateFile))
		}

		if req.Spec != nil {
			var spec interface{} = req.Spec
			fileBackedSpec, ok := spec.(*types.FileBackedVirtualDiskSpec)

			if !ok {
				return fm.fault(name, nil, new(types.FileFault))
			}

			if _, err = f.WriteString(strconv.FormatInt(fileBackedSpec.CapacityKb, 10)); err != nil {
				return fm.fault(name, err, new(types.FileFault))
			}
		}

		_ = f.Close()
	}

	return nil
}

func vdmExtendVirtualDisk(req *types.ExtendVirtualDisk_Task) types.BaseMethodFault {
	fm := Map.FileManager()

	file, fault := fm.resolve(req.Datacenter, req.Name)
	if fault != nil {
		return fault
	}

	for _, name := range vdmNames(file) {
		_, err := os.Stat(name)
		if err == nil {
			content, err := os.ReadFile(name)
			if err != nil {
				return fm.fault(name, err, new(types.FileFault))
			}

			capacity, err := strconv.Atoi(string(content))
			if err != nil {
				return fm.fault(name, err, new(types.FileFault))
			}

			if int64(capacity) > req.NewCapacityKb {
				// cannot shrink disk
				return fm.fault(name, nil, new(types.FileFault))
			}
			return nil
		} else {
			return fm.fault(name, nil, new(types.FileNotFound))
		}
	}

	return nil
}

func (m *VirtualDiskManager) CreateVirtualDiskTask(ctx *Context, req *types.CreateVirtualDisk_Task) soap.HasFault {
	task := CreateTask(m, "createVirtualDisk", func(*Task) (types.AnyType, types.BaseMethodFault) {
		if err := vdmCreateVirtualDisk(types.VirtualDeviceConfigSpecFileOperationCreate, req); err != nil {
			return "", err
		}
		return req.Name, nil
	})

	return &methods.CreateVirtualDisk_TaskBody{
		Res: &types.CreateVirtualDisk_TaskResponse{
			Returnval: task.Run(ctx),
		},
	}
}

func (m *VirtualDiskManager) ExtendVirtualDiskTask(ctx *Context, req *types.ExtendVirtualDisk_Task) soap.HasFault {
	task := CreateTask(m, "extendVirtualDisk", func(*Task) (types.AnyType, types.BaseMethodFault) {
		if err := vdmExtendVirtualDisk(req); err != nil {
			return "", err
		}
		return req.Name, nil
	})

	return &methods.ExtendVirtualDisk_TaskBody{
		Res: &types.ExtendVirtualDisk_TaskResponse{
			Returnval: task.Run(ctx),
		},
	}
}

func (m *VirtualDiskManager) DeleteVirtualDiskTask(ctx *Context, req *types.DeleteVirtualDisk_Task) soap.HasFault {
	task := CreateTask(m, "deleteVirtualDisk", func(*Task) (types.AnyType, types.BaseMethodFault) {
		fm := Map.FileManager()

		for _, name := range vdmNames(req.Name) {
			err := fm.deleteDatastoreFile(&types.DeleteDatastoreFile_Task{
				Name:       name,
				Datacenter: req.Datacenter,
			})

			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	return &methods.DeleteVirtualDisk_TaskBody{
		Res: &types.DeleteVirtualDisk_TaskResponse{
			Returnval: task.Run(ctx),
		},
	}
}

func (m *VirtualDiskManager) MoveVirtualDiskTask(ctx *Context, req *types.MoveVirtualDisk_Task) soap.HasFault {
	task := CreateTask(m, "moveVirtualDisk", func(*Task) (types.AnyType, types.BaseMethodFault) {
		fm := ctx.Map.FileManager()

		dest := vdmNames(req.DestName)

		for i, name := range vdmNames(req.SourceName) {
			err := fm.moveDatastoreFile(&types.MoveDatastoreFile_Task{
				SourceName:            name,
				SourceDatacenter:      req.SourceDatacenter,
				DestinationName:       dest[i],
				DestinationDatacenter: req.DestDatacenter,
				Force:                 req.Force,
			})

			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	return &methods.MoveVirtualDisk_TaskBody{
		Res: &types.MoveVirtualDisk_TaskResponse{
			Returnval: task.Run(ctx),
		},
	}
}

func (m *VirtualDiskManager) CopyVirtualDiskTask(ctx *Context, req *types.CopyVirtualDisk_Task) soap.HasFault {
	task := CreateTask(m, "copyVirtualDisk", func(*Task) (types.AnyType, types.BaseMethodFault) {
		if req.DestSpec != nil {
			// TODO: apply to destination vmdk.Descriptor
		}

		fm := ctx.Map.FileManager()

		dest := vdmNames(req.DestName)

		for i, name := range vdmNames(req.SourceName) {
			err := fm.copyDatastoreFile(&types.CopyDatastoreFile_Task{
				SourceName:            name,
				SourceDatacenter:      req.SourceDatacenter,
				DestinationName:       dest[i],
				DestinationDatacenter: req.DestDatacenter,
				Force:                 req.Force,
			})

			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	return &methods.CopyVirtualDisk_TaskBody{
		Res: &types.CopyVirtualDisk_TaskResponse{
			Returnval: task.Run(ctx),
		},
	}
}

func virtualDiskUUID(dc *types.ManagedObjectReference, file string) string {
	if dc != nil {
		file = dc.String() + file
	}
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(file)).String()
}

func (m *VirtualDiskManager) QueryVirtualDiskUuid(ctx *Context, req *types.QueryVirtualDiskUuid) soap.HasFault {
	body := new(methods.QueryVirtualDiskUuidBody)

	fm := ctx.Map.FileManager()

	file, fault := fm.resolve(req.Datacenter, req.Name)
	if fault != nil {
		body.Fault_ = Fault("", fault)
		return body
	}

	_, err := os.Stat(file)
	if err != nil {
		fault = fm.fault(req.Name, err, new(types.CannotAccessFile))
		body.Fault_ = Fault(fmt.Sprintf("File %s was not found", req.Name), fault)
		return body
	}

	body.Res = &types.QueryVirtualDiskUuidResponse{
		Returnval: virtualDiskUUID(req.Datacenter, req.Name),
	}

	return body
}

func (m *VirtualDiskManager) SetVirtualDiskUuid(_ *Context, req *types.SetVirtualDiskUuid) soap.HasFault {
	body := new(methods.SetVirtualDiskUuidBody)
	// TODO: validate uuid format and persist
	body.Res = new(types.SetVirtualDiskUuidResponse)
	return body
}

func (m *VirtualDiskManager) QueryVirtualDiskInfoTask(ctx *Context, req *internal.QueryVirtualDiskInfoTaskRequest) soap.HasFault {
	task := CreateTask(m, "queryVirtualDiskInfo", func(*Task) (types.AnyType, types.BaseMethodFault) {
		var res []internal.VirtualDiskInfo

		fm := ctx.Map.FileManager()

		_, fault := fm.resolve(req.Datacenter, req.Name)
		if fault != nil {
			return nil, fault
		}

		res = append(res, internal.VirtualDiskInfo{Name: req.Name, DiskType: "thin"})

		if req.IncludeParents {
			// TODO
		}

		return res, nil
	})

	return &internal.QueryVirtualDiskInfoTaskBody{
		Res: &internal.QueryVirtualDiskInfo_TaskResponse{
			Returnval: task.Run(ctx),
		},
	}
}
