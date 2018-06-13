// Copyright 2018 Portieris Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fakenotary

import (
	"sync"

	notaryclient "github.com/theupdateframework/notary/client"
	"github.com/theupdateframework/notary/client/changelist"
	"github.com/theupdateframework/notary/tuf/data"
	"github.com/theupdateframework/notary/tuf/signed"
)

// FakeRepository .
type FakeRepository struct {
	InitializeStub        func(rootKeyIDs []string, serverManagedRoles ...data.RoleName) error
	initializeMutex       sync.RWMutex
	initializeArgsForCall []struct {
		rootKeyIDs         []string
		serverManagedRoles []data.RoleName
	}
	initializeReturns struct {
		err error
	}

	InitializeWithCertificateStub        func(rootKeyIDs []string, rootCerts []data.PublicKey, serverManagedRoles ...data.RoleName) error
	initializeWithCertificateMutex       sync.RWMutex
	initializeWithCertificateArgsForCall []struct {
		rootKeyIDs         []string
		rootCerts          []data.PublicKey
		serverManagedRoles []data.RoleName
	}
	initializeWithCertificateReturns struct {
		err error
	}

	PublishStub        func() error
	publishMutex       sync.RWMutex
	publishArgsForCall []struct {
	}
	publishReturns struct {
		err error
	}

	AddTargetStub        func(target *notaryclient.Target, roles ...data.RoleName) error
	addTargetMutex       sync.RWMutex
	addTargetArgsForCall []struct {
		target *notaryclient.Target
		roles  []data.RoleName
	}
	addTargetReturns struct {
		err error
	}

	RemoveTargetStub        func(targetName string, roles ...data.RoleName) error
	removeTargetMutex       sync.RWMutex
	removeTargetArgsForCall []struct {
		targetName string
		roles      []data.RoleName
	}
	removeTargetReturns struct {
		err error
	}

	ListTargetsStub        func(roles ...data.RoleName) ([]*notaryclient.TargetWithRole, error)
	listTargetsMutex       sync.RWMutex
	listTargetsArgsForCall []struct {
		roles []data.RoleName
	}
	listTargetsReturns struct {
		targetWithRole []*notaryclient.TargetWithRole
		err            error
	}

	GetTargetByNameStub        func(name string, roles ...data.RoleName) (*notaryclient.TargetWithRole, error)
	getTargetByNameMutex       sync.RWMutex
	getTargetByNameArgsForCall []struct {
		name string
	}
	getTargetByNameReturns struct {
		targetWithRole *notaryclient.TargetWithRole
		err            error
	}

	GetAllTargetMetadataByNameStub        func(name string) ([]notaryclient.TargetSignedStruct, error)
	getAllTargetMetadataByNameMutex       sync.RWMutex
	getAllTargetMetadataByNameArgsForCall []struct {
		name string
	}
	getAllTargetMetadataByNameReturns struct {
		targetSignedStruct []notaryclient.TargetSignedStruct
		err                error
	}

	GetChangelistStub        func() (changelist.Changelist, error)
	getChangelistMutex       sync.RWMutex
	getChangelistArgsForCall []struct {
	}
	getChangelistReturns struct {
		list changelist.Changelist
		err  error
	}

	ListRolesStub        func() ([]notaryclient.RoleWithSignatures, error)
	listRolesMutex       sync.RWMutex
	listRolesArgsForCall []struct {
	}
	listRolesReturns struct {
		roleWithSignatures []notaryclient.RoleWithSignatures
		err                error
	}

	GetDelegationRolesStub        func() ([]data.Role, error)
	getDelegationRolesMutex       sync.RWMutex
	getDelegationRolesArgsForCall []struct {
	}
	getDelegationRolesReturns struct {
		role []data.Role
		err  error
	}

	AddDelegationStub        func(name data.RoleName, delegationKeys []data.PublicKey, paths []string) error
	addDelegationMutex       sync.RWMutex
	addDelegationArgsForCall []struct {
		name           data.RoleName
		delegationKeys []data.PublicKey
		paths          []string
	}
	addDelegationReturns struct {
		err error
	}

	AddDelegationRoleAndKeysStub        func(name data.RoleName, delegationKeys []data.PublicKey) error
	addDelegationRoleAndKeysMutex       sync.RWMutex
	addDelegationRoleAndKeysArgsForCall []struct {
		name           data.RoleName
		delegationKeys []data.PublicKey
	}
	addDelegationRoleAndKeysReturns struct {
		err error
	}

	AddDelegationPathsStub        func(name data.RoleName, paths []string) error
	addDelegationPathsMutex       sync.RWMutex
	addDelegationPathsArgsForCall []struct {
		name  data.RoleName
		paths []string
	}
	addDelegationPathsReturns struct {
		err error
	}

	RemoveDelegationKeysAndPathsStub        func(name data.RoleName, keyIDs []string, paths []string) error
	removeDelegationKeysAndPathsMutex       sync.RWMutex
	removeDelegationKeysAndPathsArgsForCall []struct {
		name   data.RoleName
		keyIDs []string
		paths  []string
	}
	removeDelegationKeysAndPathsReturns struct {
		err error
	}

	RemoveDelegationRoleStub        func(name data.RoleName) error
	removeDelegationRoleMutex       sync.RWMutex
	removeDelegationRoleArgsForCall []struct {
		name data.RoleName
	}
	removeDelegationRoleReturns struct {
		err error
	}

	RemoveDelegationPathsStub        func(name data.RoleName, paths []string) error
	removeDelegationPathsMutex       sync.RWMutex
	removeDelegationPathsArgsForCall []struct {
		name  data.RoleName
		paths []string
	}
	removeDelegationPathsReturns struct {
		err error
	}

	RemoveDelegationKeysStub        func(name data.RoleName, keyIDs []string) error
	removeDelegationKeysMutex       sync.RWMutex
	removeDelegationKeysArgsForCall []struct {
		name   data.RoleName
		keyIDs []string
	}
	removeDelegationKeysReturns struct {
		err error
	}

	ClearDelegationPathsStub        func(name data.RoleName) error
	clearDelegationPathsMutex       sync.RWMutex
	clearDelegationPathsArgsForCall []struct {
		name data.RoleName
	}
	clearDelegationPathsReturns struct {
		err error
	}

	WitnessStub        func(roles ...data.RoleName) ([]data.RoleName, error)
	witnessMutex       sync.RWMutex
	witnessArgsForCall []struct {
		roles []data.RoleName
	}
	witnessReturns struct {
		roleName []data.RoleName
		err      error
	}

	RotateKeyStub        func(role data.RoleName, serverManagesKey bool, keyList []string) error
	rotateKeyMutex       sync.RWMutex
	rotateKeyArgsForCall []struct {
		role             data.RoleName
		serverManagesKey bool
		keyList          []string
	}
	rotateKeyReturns struct {
		err error
	}

	GetCryptoServiceStub        func() signed.CryptoService
	getCryptoServiceMutex       sync.RWMutex
	getCryptoServiceArgsForCall []struct {
	}
	getCryptoServiceReturns struct {
		cryptoService signed.CryptoService
	}

	GetGUNStub        func() data.GUN
	getGUNMutex       sync.RWMutex
	getGUNArgsForCall []struct {
	}
	getGUNReturns struct {
		gun data.GUN
	}
}

// Initialize .
func (f *FakeRepository) Initialize(rootKeyIDs []string, serverManagedRoles ...data.RoleName) error {
	f.initializeMutex.Lock()
	f.initializeArgsForCall = append(f.initializeArgsForCall, struct {
		rootKeyIDs         []string
		serverManagedRoles []data.RoleName
	}{rootKeyIDs, serverManagedRoles})
	f.initializeMutex.Unlock()
	if f.InitializeStub != nil {
		return f.InitializeStub(rootKeyIDs, serverManagedRoles...)
	}
	return f.initializeReturns.err
}

// InitializeReturns .
func (f *FakeRepository) InitializeReturns(err error) {
	f.initializeMutex.Lock()
	defer f.initializeMutex.Unlock()
	f.initializeReturns = struct {
		err error
	}{err}
}

// InitializeWithCertificate .
func (f *FakeRepository) InitializeWithCertificate(rootKeyIDs []string, rootCerts []data.PublicKey, serverManagedRoles ...data.RoleName) error {
	f.initializeWithCertificateMutex.Lock()
	f.initializeWithCertificateArgsForCall = append(f.initializeWithCertificateArgsForCall, struct {
		rootKeyIDs         []string
		rootCerts          []data.PublicKey
		serverManagedRoles []data.RoleName
	}{rootKeyIDs, rootCerts, serverManagedRoles})
	f.initializeWithCertificateMutex.Unlock()
	if f.InitializeWithCertificateStub != nil {
		return f.InitializeWithCertificateStub(rootKeyIDs, rootCerts, serverManagedRoles...)
	}
	return f.initializeWithCertificateReturns.err
}

// InitializeWithCertificateReturns .
func (f *FakeRepository) InitializeWithCertificateReturns(err error) {
	f.initializeWithCertificateMutex.Lock()
	defer f.initializeWithCertificateMutex.Unlock()
	f.initializeWithCertificateReturns = struct {
		err error
	}{err}
}

// Publish .
func (f *FakeRepository) Publish() error {
	f.publishMutex.Lock()
	f.publishArgsForCall = append(f.publishArgsForCall, struct {
	}{})
	f.publishMutex.Unlock()
	if f.PublishStub != nil {
		return f.PublishStub()
	}
	return f.publishReturns.err
}

// PublishReturns .
func (f *FakeRepository) PublishReturns(err error) {
	f.publishMutex.Lock()
	defer f.publishMutex.Unlock()
	f.publishReturns = struct {
		err error
	}{err}
}

// AddTarget .
func (f *FakeRepository) AddTarget(target *notaryclient.Target, roles ...data.RoleName) error {
	f.addTargetMutex.Lock()
	f.addTargetArgsForCall = append(f.addTargetArgsForCall, struct {
		target *notaryclient.Target
		roles  []data.RoleName
	}{target, roles})
	f.addTargetMutex.Unlock()
	if f.AddTargetStub != nil {
		return f.AddTargetStub(target, roles...)
	}
	return f.addTargetReturns.err
}

// AddTargetReturns .
func (f *FakeRepository) AddTargetReturns(err error) {
	f.addTargetMutex.Lock()
	defer f.addTargetMutex.Unlock()
	f.addTargetReturns = struct {
		err error
	}{err}
}

// RemoveTarget .
func (f *FakeRepository) RemoveTarget(targetName string, roles ...data.RoleName) error {
	f.removeTargetMutex.Lock()
	f.removeTargetArgsForCall = append(f.removeTargetArgsForCall, struct {
		targetName string
		roles      []data.RoleName
	}{targetName, roles})
	f.removeTargetMutex.Unlock()
	if f.RemoveTargetStub != nil {
		return f.RemoveTargetStub(targetName, roles...)
	}
	return f.removeTargetReturns.err
}

// RemoveTargetReturns .
func (f *FakeRepository) RemoveTargetReturns(err error) {
	f.removeTargetMutex.Lock()
	defer f.removeTargetMutex.Unlock()
	f.removeTargetReturns = struct {
		err error
	}{err}
}

// ListTargets .
func (f *FakeRepository) ListTargets(roles ...data.RoleName) ([]*notaryclient.TargetWithRole, error) {
	f.listTargetsMutex.Lock()
	f.listTargetsArgsForCall = append(f.listTargetsArgsForCall, struct {
		roles []data.RoleName
	}{roles})
	f.listTargetsMutex.Unlock()
	if f.ListTargetsStub != nil {
		return f.ListTargetsStub(roles...)
	}
	return f.listTargetsReturns.targetWithRole, f.listTargetsReturns.err
}

// ListTargetsReturns .
func (f *FakeRepository) ListTargetsReturns(targetWithRole []*notaryclient.TargetWithRole, err error) {
	f.listTargetsMutex.Lock()
	defer f.listTargetsMutex.Unlock()
	f.listTargetsReturns = struct {
		targetWithRole []*notaryclient.TargetWithRole
		err            error
	}{targetWithRole, err}
}

// GetTargetByName ...
func (f *FakeRepository) GetTargetByName(name string, roles ...data.RoleName) (*notaryclient.TargetWithRole, error) {
	f.getTargetByNameMutex.Lock()
	f.getTargetByNameArgsForCall = append(f.getTargetByNameArgsForCall, struct {
		name string
	}{name})
	f.getTargetByNameMutex.Unlock()
	if f.GetTargetByNameStub != nil {
		return f.GetTargetByNameStub(name)
	}
	return f.getTargetByNameReturns.targetWithRole, f.getTargetByNameReturns.err
}

// GetTargetByNameReturns .
func (f *FakeRepository) GetTargetByNameReturns(targetWithRole *notaryclient.TargetWithRole, err error) {
	f.getTargetByNameMutex.Lock()
	defer f.getTargetByNameMutex.Unlock()
	f.getTargetByNameReturns = struct {
		targetWithRole *notaryclient.TargetWithRole
		err            error
	}{targetWithRole, err}
}

// GetAllTargetMetadataByName .
func (f *FakeRepository) GetAllTargetMetadataByName(name string) ([]notaryclient.TargetSignedStruct, error) {
	f.getAllTargetMetadataByNameMutex.Lock()
	f.getAllTargetMetadataByNameArgsForCall = append(f.getAllTargetMetadataByNameArgsForCall, struct {
		name string
	}{name})
	f.getAllTargetMetadataByNameMutex.Unlock()
	if f.GetAllTargetMetadataByNameStub != nil {
		return f.GetAllTargetMetadataByNameStub(name)
	}
	return f.getAllTargetMetadataByNameReturns.targetSignedStruct, f.getAllTargetMetadataByNameReturns.err
}

// GetAllTargetMetadataByNameReturns .
func (f *FakeRepository) GetAllTargetMetadataByNameReturns(targetSignedStruct []notaryclient.TargetSignedStruct, err error) {
	f.getAllTargetMetadataByNameMutex.Lock()
	defer f.getAllTargetMetadataByNameMutex.Unlock()
	f.getAllTargetMetadataByNameReturns = struct {
		targetSignedStruct []notaryclient.TargetSignedStruct
		err                error
	}{targetSignedStruct, err}
}

// GetChangelist .
func (f *FakeRepository) GetChangelist() (changelist.Changelist, error) {
	f.getChangelistMutex.Lock()
	f.getChangelistArgsForCall = append(f.getChangelistArgsForCall, struct {
	}{})
	f.getChangelistMutex.Unlock()
	if f.GetChangelistStub != nil {
		return f.GetChangelistStub()
	}
	return f.getChangelistReturns.list, f.getChangelistReturns.err
}

// GetChangelistReturns .
func (f *FakeRepository) GetChangelistReturns(list changelist.Changelist, err error) {
	f.getChangelistMutex.Lock()
	defer f.getChangelistMutex.Unlock()
	f.getChangelistReturns = struct {
		list changelist.Changelist
		err  error
	}{list, err}
}

// ListRoles .
func (f *FakeRepository) ListRoles() ([]notaryclient.RoleWithSignatures, error) {
	f.listRolesMutex.Lock()
	f.listRolesArgsForCall = append(f.listRolesArgsForCall, struct {
	}{})
	f.listRolesMutex.Unlock()
	if f.ListRolesStub != nil {
		return f.ListRolesStub()
	}
	return f.listRolesReturns.roleWithSignatures, f.listRolesReturns.err
}

// ListRolesReturns .
func (f *FakeRepository) ListRolesReturns(roleWithSignatures []notaryclient.RoleWithSignatures, err error) {
	f.listRolesMutex.Lock()
	defer f.listRolesMutex.Unlock()
	f.listRolesReturns = struct {
		roleWithSignatures []notaryclient.RoleWithSignatures
		err                error
	}{roleWithSignatures, err}
}

// GetDelegationRoles .
func (f *FakeRepository) GetDelegationRoles() ([]data.Role, error) {
	f.getDelegationRolesMutex.Lock()
	f.getDelegationRolesArgsForCall = append(f.getDelegationRolesArgsForCall, struct {
	}{})
	f.getDelegationRolesMutex.Unlock()
	if f.GetDelegationRolesStub != nil {
		return f.GetDelegationRolesStub()
	}
	return f.getDelegationRolesReturns.role, f.getDelegationRolesReturns.err
}

// GetDelegationRolesReturns .
func (f *FakeRepository) GetDelegationRolesReturns(role []data.Role, err error) {
	f.getDelegationRolesMutex.Lock()
	defer f.getDelegationRolesMutex.Unlock()
	f.getDelegationRolesReturns = struct {
		role []data.Role
		err  error
	}{role, err}
}

// AddDelegation .
func (f *FakeRepository) AddDelegation(name data.RoleName, delegationKeys []data.PublicKey, paths []string) error {
	f.addDelegationMutex.Lock()
	f.addDelegationArgsForCall = append(f.addDelegationArgsForCall, struct {
		name           data.RoleName
		delegationKeys []data.PublicKey
		paths          []string
	}{name, delegationKeys, paths})
	f.addDelegationMutex.Unlock()
	if f.AddDelegationStub != nil {
		return f.AddDelegationStub(name, delegationKeys, paths)
	}
	return f.addDelegationReturns.err
}

// AddDelegationReturns .
func (f *FakeRepository) AddDelegationReturns(err error) {
	f.addDelegationMutex.Lock()
	defer f.addDelegationMutex.Unlock()
	f.addDelegationReturns = struct {
		err error
	}{err}
}

// AddDelegationRoleAndKeys .
func (f *FakeRepository) AddDelegationRoleAndKeys(name data.RoleName, delegationKeys []data.PublicKey) error {
	f.addDelegationRoleAndKeysMutex.Lock()
	f.addDelegationRoleAndKeysArgsForCall = append(f.addDelegationRoleAndKeysArgsForCall, struct {
		name           data.RoleName
		delegationKeys []data.PublicKey
	}{name, delegationKeys})
	f.addDelegationRoleAndKeysMutex.Unlock()
	if f.AddDelegationRoleAndKeysStub != nil {
		return f.AddDelegationRoleAndKeysStub(name, delegationKeys)
	}
	return f.addDelegationRoleAndKeysReturns.err
}

// AddDelegationRoleAndKeysReturns .
func (f *FakeRepository) AddDelegationRoleAndKeysReturns(err error) {
	f.addDelegationRoleAndKeysMutex.Lock()
	defer f.addDelegationRoleAndKeysMutex.Unlock()
	f.addDelegationRoleAndKeysReturns = struct {
		err error
	}{err}
}

// AddDelegationPaths .
func (f *FakeRepository) AddDelegationPaths(name data.RoleName, paths []string) error {
	f.addDelegationPathsMutex.Lock()
	f.addDelegationPathsArgsForCall = append(f.addDelegationPathsArgsForCall, struct {
		name  data.RoleName
		paths []string
	}{name, paths})
	f.addDelegationPathsMutex.Unlock()
	if f.AddDelegationPathsStub != nil {
		return f.AddDelegationPathsStub(name, paths)
	}
	return f.addDelegationPathsReturns.err
}

// AddDelegationPathsReturns .
func (f *FakeRepository) AddDelegationPathsReturns(err error) {
	f.addDelegationPathsMutex.Lock()
	defer f.addDelegationPathsMutex.Unlock()
	f.addDelegationPathsReturns = struct {
		err error
	}{err}
}

// RemoveDelegationKeysAndPaths .
func (f *FakeRepository) RemoveDelegationKeysAndPaths(name data.RoleName, keyIDs []string, paths []string) error {
	f.removeDelegationKeysAndPathsMutex.Lock()
	f.removeDelegationKeysAndPathsArgsForCall = append(f.removeDelegationKeysAndPathsArgsForCall, struct {
		name   data.RoleName
		keyIDs []string
		paths  []string
	}{name, keyIDs, paths})
	f.removeDelegationKeysAndPathsMutex.Unlock()
	if f.RemoveDelegationKeysAndPathsStub != nil {
		return f.RemoveDelegationKeysAndPathsStub(name, keyIDs, paths)
	}
	return f.removeDelegationKeysAndPathsReturns.err
}

// RemoveDelegationKeysAndPathsReturns .
func (f *FakeRepository) RemoveDelegationKeysAndPathsReturns(err error) {
	f.removeDelegationKeysAndPathsMutex.Lock()
	defer f.removeDelegationKeysAndPathsMutex.Unlock()
	f.removeDelegationKeysAndPathsReturns = struct {
		err error
	}{err}
}

// RemoveDelegationRole .
func (f *FakeRepository) RemoveDelegationRole(name data.RoleName) error {
	f.removeDelegationRoleMutex.Lock()
	f.removeDelegationRoleArgsForCall = append(f.removeDelegationRoleArgsForCall, struct {
		name data.RoleName
	}{name})
	f.removeDelegationRoleMutex.Unlock()
	if f.RemoveDelegationRoleStub != nil {
		return f.RemoveDelegationRoleStub(name)
	}
	return f.removeDelegationRoleReturns.err
}

// RemoveDelegationRoleReturns .
func (f *FakeRepository) RemoveDelegationRoleReturns(err error) {
	f.removeDelegationRoleMutex.Lock()
	defer f.removeDelegationRoleMutex.Unlock()
	f.removeDelegationRoleReturns = struct {
		err error
	}{err}
}

// RemoveDelegationPaths .
func (f *FakeRepository) RemoveDelegationPaths(name data.RoleName, paths []string) error {
	f.removeDelegationPathsMutex.Lock()
	f.removeDelegationPathsArgsForCall = append(f.removeDelegationPathsArgsForCall, struct {
		name  data.RoleName
		paths []string
	}{name, paths})
	f.removeDelegationPathsMutex.Unlock()
	if f.RemoveDelegationPathsStub != nil {
		return f.RemoveDelegationPathsStub(name, paths)
	}
	return f.removeDelegationPathsReturns.err
}

// RemoveDelegationPathsReturns .
func (f *FakeRepository) RemoveDelegationPathsReturns(err error) {
	f.removeDelegationPathsMutex.Lock()
	defer f.removeDelegationPathsMutex.Unlock()
	f.removeDelegationPathsReturns = struct {
		err error
	}{err}
}

// RemoveDelegationKeys .
func (f *FakeRepository) RemoveDelegationKeys(name data.RoleName, keyIDs []string) error {
	f.removeDelegationKeysMutex.Lock()
	f.removeDelegationKeysArgsForCall = append(f.removeDelegationKeysArgsForCall, struct {
		name   data.RoleName
		keyIDs []string
	}{name, keyIDs})
	f.removeDelegationKeysMutex.Unlock()
	if f.RemoveDelegationKeysStub != nil {
		return f.RemoveDelegationKeysStub(name, keyIDs)
	}
	return f.removeDelegationKeysReturns.err
}

// RemoveDelegationKeysReturns .
func (f *FakeRepository) RemoveDelegationKeysReturns(err error) {
	f.removeDelegationKeysMutex.Lock()
	defer f.removeDelegationKeysMutex.Unlock()
	f.removeDelegationKeysReturns = struct {
		err error
	}{err}
}

// ClearDelegationPaths .
func (f *FakeRepository) ClearDelegationPaths(name data.RoleName) error {
	f.clearDelegationPathsMutex.Lock()
	f.clearDelegationPathsArgsForCall = append(f.clearDelegationPathsArgsForCall, struct {
		name data.RoleName
	}{name})
	f.clearDelegationPathsMutex.Unlock()
	if f.ClearDelegationPathsStub != nil {
		return f.ClearDelegationPathsStub(name)
	}
	return f.clearDelegationPathsReturns.err
}

// ClearDelegationPathsReturns .
func (f *FakeRepository) ClearDelegationPathsReturns(err error) {
	f.clearDelegationPathsMutex.Lock()
	defer f.clearDelegationPathsMutex.Unlock()
	f.clearDelegationPathsReturns = struct {
		err error
	}{err}
}

// Witness .
func (f *FakeRepository) Witness(roles ...data.RoleName) ([]data.RoleName, error) {
	f.witnessMutex.Lock()
	f.witnessArgsForCall = append(f.witnessArgsForCall, struct {
		roles []data.RoleName
	}{roles})
	f.witnessMutex.Unlock()
	if f.WitnessStub != nil {
		return f.WitnessStub(roles...)
	}
	return f.witnessReturns.roleName, f.witnessReturns.err
}

// WitnessReturns .
func (f *FakeRepository) WitnessReturns(roleName []data.RoleName, err error) {
	f.witnessMutex.Lock()
	defer f.witnessMutex.Unlock()
	f.witnessReturns = struct {
		roleName []data.RoleName
		err      error
	}{roleName, err}
}

// RotateKey .
func (f *FakeRepository) RotateKey(role data.RoleName, serverManagesKey bool, keyList []string) error {
	f.rotateKeyMutex.Lock()
	f.rotateKeyArgsForCall = append(f.rotateKeyArgsForCall, struct {
		role             data.RoleName
		serverManagesKey bool
		keyList          []string
	}{role, serverManagesKey, keyList})
	f.rotateKeyMutex.Unlock()
	if f.RotateKeyStub != nil {
		return f.RotateKeyStub(role, serverManagesKey, keyList)
	}
	return f.rotateKeyReturns.err
}

// RotateKeyReturns .
func (f *FakeRepository) RotateKeyReturns(err error) {
	f.rotateKeyMutex.Lock()
	defer f.rotateKeyMutex.Unlock()
	f.rotateKeyReturns = struct {
		err error
	}{err}
}

// GetCryptoService .
func (f *FakeRepository) GetCryptoService() signed.CryptoService {
	f.getCryptoServiceMutex.Lock()
	f.getCryptoServiceArgsForCall = append(f.getCryptoServiceArgsForCall, struct {
	}{})
	f.getCryptoServiceMutex.Unlock()
	if f.GetCryptoServiceStub != nil {
		return f.GetCryptoServiceStub()
	}
	return f.getCryptoServiceReturns.cryptoService
}

// GetCryptoServiceReturns .
func (f *FakeRepository) GetCryptoServiceReturns(cryptoService signed.CryptoService) {
	f.getCryptoServiceMutex.Lock()
	defer f.getCryptoServiceMutex.Unlock()
	f.getCryptoServiceReturns = struct {
		cryptoService signed.CryptoService
	}{cryptoService}
}

// SetLegacyVersions .
func (f *FakeRepository) SetLegacyVersions(version int) {
}

// GetGUN .
func (f *FakeRepository) GetGUN() data.GUN {
	f.getGUNMutex.Lock()
	f.getGUNArgsForCall = append(f.getGUNArgsForCall, struct {
	}{})
	f.getGUNMutex.Unlock()
	if f.GetGUNStub != nil {
		return f.GetGUNStub()
	}
	return f.getGUNReturns.gun
}

// GetGUNReturns .
func (f *FakeRepository) GetGUNReturns(gun data.GUN) {
	f.getGUNMutex.Lock()
	defer f.getGUNMutex.Unlock()
	f.getGUNReturns = struct {
		gun data.GUN
	}{gun}
}
