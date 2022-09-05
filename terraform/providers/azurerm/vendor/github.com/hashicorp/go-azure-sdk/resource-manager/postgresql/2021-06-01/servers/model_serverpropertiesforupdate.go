package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerPropertiesForUpdate struct {
	AdministratorLoginPassword *string              `json:"administratorLoginPassword,omitempty"`
	Backup                     *Backup              `json:"backup,omitempty"`
	CreateMode                 *CreateModeForUpdate `json:"createMode,omitempty"`
	HighAvailability           *HighAvailability    `json:"highAvailability,omitempty"`
	MaintenanceWindow          *MaintenanceWindow   `json:"maintenanceWindow,omitempty"`
	Storage                    *Storage             `json:"storage,omitempty"`
}
