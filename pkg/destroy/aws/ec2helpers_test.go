package aws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/destroy/aws/mock"
)

var nullLogger = func() logrus.FieldLogger {
	logger := logrus.StandardLogger()
	logger.SetOutput(io.Discard)
	return logger
}()

func TestDeleteEC2SecurityGroupObject(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ec2Client := mock.NewMockEC2API(mockCtrl)

	const (
		defaultGroupID    = "arn:default-group-id"
		defaultGroupName  = "default"
		failDeleteGroupID = "arn:fail-group-id"
	)

	perm := ec2.IpPermission{}

	defaultGroup := func(g *ec2.SecurityGroup) {
		g.GroupId = aws.String(defaultGroupID)
		g.GroupName = aws.String(defaultGroupName)
	}
	nonDefaultGroup := func(g *ec2.SecurityGroup) {
		g.GroupName = aws.String("group-name")
	}
	failDeleteGroup := func(g *ec2.SecurityGroup) {
		g.GroupName = aws.String("fail-group-name")
		g.GroupId = aws.String(failDeleteGroupID)
	}
	addIngressPermission := func(g *ec2.SecurityGroup) {
		g.IpPermissions = append(g.IpPermissions, &perm)
	}
	addEgressPermission := func(g *ec2.SecurityGroup) {
		g.IpPermissionsEgress = append(g.IpPermissionsEgress, &perm)
	}

	type editGroupFuncs []func(g *ec2.SecurityGroup)
	cases := []struct {
		name      string
		editFuncs editGroupFuncs
		errorMsg  string
	}{
		{
			name:      "SecurityGroup Ingress revoked",
			editFuncs: editGroupFuncs{addIngressPermission},
			errorMsg:  "",
		},
		{
			name:      "SecurityGroup Ingress revoke error",
			editFuncs: editGroupFuncs{addIngressPermission, addIngressPermission},
			errorMsg:  "revoking ingress permissions",
		},
		{
			name:      "SecurityGroup Egress revoked",
			editFuncs: editGroupFuncs{addEgressPermission},
			errorMsg:  "",
		},
		{
			name:      "SecurityGroup Egress revoke error",
			editFuncs: editGroupFuncs{addEgressPermission, addEgressPermission},
			errorMsg:  "revoking egress permissions",
		},
		{
			name:      "SecurityGroup default group is not deleted",
			editFuncs: editGroupFuncs{defaultGroup},
			errorMsg:  "",
		},
		{
			name:      "SecurityGroup null name is deleted",
			editFuncs: editGroupFuncs{},
			errorMsg:  "",
		},
		{
			name:      "SecurityGroup non-default group is deleted",
			editFuncs: editGroupFuncs{nonDefaultGroup},
			errorMsg:  "",
		},
		{
			name:      "SecurityGroup delete fails",
			editFuncs: editGroupFuncs{failDeleteGroup},
			errorMsg:  "cannot delete",
		},
	}

	ec2Client.
		EXPECT().
		RevokeSecurityGroupEgress(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, g *ec2.SecurityGroup) error {
			switch len(g.IpPermissionsEgress) {
			case 0:
				return errors.New("Should not be called")
			case 1:
				return nil
			default:
				return errors.New("Egress revoke error")
			}
		}).
		AnyTimes()
	ec2Client.
		EXPECT().
		RevokeSecurityGroupIngress(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, g *ec2.SecurityGroup) error {
			switch len(g.IpPermissions) {
			case 0:
				return errors.New("Should not be called")
			case 1:
				return nil
			default:
				return errors.New("Ingress revoke error")
			}
		}).
		AnyTimes()
	ec2Client.
		EXPECT().
		DeleteSecurityGroup(gomock.Any(), gomock.Eq(defaultGroupID)).
		Times(0) // Default security group should not be deleted
	ec2Client.
		EXPECT().
		DeleteSecurityGroup(gomock.Any(), gomock.Eq(failDeleteGroupID)).
		Return(errors.New("cannot delete"))
	ec2Client.
		EXPECT().
		DeleteSecurityGroup(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedGroup := &ec2.SecurityGroup{GroupId: aws.String("arn:group-id")}
			for _, edit := range tc.editFuncs {
				edit(editedGroup)
			}
			err := deleteEC2SecurityGroupObject(context.TODO(), ec2Client, editedGroup, nullLogger)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteRouteTableObject(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	ec2Client := mock.NewMockEC2API(mockCtrl)

	const (
		mainAssociationID = "arn:main-association-id"
		mainRouteTableID  = "arn:main-route-table-id"
		failAssociationID = "arn:fail-association-id"
		failRouteTableID  = "arn:fail-route-table-id"
	)

	association := &ec2.RouteTableAssociation{
		RouteTableAssociationId: aws.String("association-id"),
		Main:                    aws.Bool(false),
	}
	failAssociation := &ec2.RouteTableAssociation{
		RouteTableAssociationId: aws.String(failAssociationID),
		Main:                    aws.Bool(false),
	}
	mainAssociation := &ec2.RouteTableAssociation{
		RouteTableAssociationId: aws.String(mainAssociationID),
		Main:                    aws.Bool(true),
	}
	addMainAssociation := func(t *ec2.RouteTable) {
		t.Associations = append(t.Associations, mainAssociation)
		t.RouteTableId = aws.String(mainRouteTableID)
	}
	addAssociation := func(t *ec2.RouteTable) {
		t.Associations = append(t.Associations, association)
	}
	addFailAssociation := func(t *ec2.RouteTable) {
		t.Associations = append(t.Associations, failAssociation)
	}
	failRouteTable := func(t *ec2.RouteTable) {
		t.RouteTableId = aws.String(failRouteTableID)
	}

	type editTableFuncs []func(t *ec2.RouteTable)

	cases := []struct {
		name      string
		editFuncs editTableFuncs
		errorMsg  string
	}{
		{
			name:      "Empty RouteTable associations",
			editFuncs: editTableFuncs{},
			errorMsg:  "",
		},
		{
			name:      "RouteTable associations are deleted except main",
			editFuncs: editTableFuncs{addAssociation, addMainAssociation},
			errorMsg:  "",
		},
		{
			name:      "RouteTable association fails",
			editFuncs: editTableFuncs{addAssociation, addMainAssociation, addFailAssociation},
			errorMsg:  "dissociating ",
		},
		{
			name:      "RouteTable with main associaton is not deleted",
			editFuncs: editTableFuncs{addAssociation, addMainAssociation},
			errorMsg:  "",
		},
		{
			name:      "RouteTable without main association is deleted",
			editFuncs: editTableFuncs{addAssociation, addAssociation},
			errorMsg:  "",
		},
		{
			name:      "RouteTable without main association fails to delete",
			editFuncs: editTableFuncs{addAssociation, addAssociation, failRouteTable},
			errorMsg:  "some AWS error",
		},
	}

	ec2Client.
		EXPECT().
		DisassociateRouteTable(gomock.Any(), gomock.Eq(mainAssociationID)).
		Times(0) // Cannot delete Main route table association
	ec2Client.
		EXPECT().
		DisassociateRouteTable(gomock.Any(), gomock.Eq(failAssociationID)).
		Return(fmt.Errorf("dissociating %s: could not delete", failAssociationID)).
		AnyTimes()
	ec2Client.
		EXPECT().
		DisassociateRouteTable(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	ec2Client.
		EXPECT().
		DeleteRouteTable(gomock.Any(), gomock.Eq(failRouteTableID)).
		Return(errors.New("some AWS error")).
		AnyTimes()
	ec2Client.
		EXPECT().
		DeleteRouteTable(gomock.Any(), gomock.Eq(mainRouteTableID)).
		Times(0) // Cannot delete main route table
	ec2Client.
		EXPECT().
		DeleteRouteTable(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedTable := &ec2.RouteTable{RouteTableId: aws.String("arn:table-id")}
			for _, edit := range tc.editFuncs {
				edit(editedTable)
			}
			err := deleteEC2RouteTableObject(context.TODO(), ec2Client, editedTable, nullLogger)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
