package flex

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"

	"github.com/timw255/flex-go/util"
)

// KinveyEntityModule ...
type KinveyEntityModule struct {
	environmentID     string
	useBSONObjectID   bool
	objectIDGenerator *objectIDGenerator
}

func newKinveyEntityModule(environmentID string, useBSONObjectID bool) KinveyEntityModule {
	return KinveyEntityModule{
		environmentID:   environmentID,
		useBSONObjectID: useBSONObjectID,
	}
}

func (m *KinveyEntityModule) createEntityID() *string {
	if m.objectIDGenerator == nil {
		idg := newObjectIDGenerator()
		m.objectIDGenerator = &idg
	}
	id := m.objectIDGenerator.NewObjectID()
	if m.useBSONObjectID {
		return id
	}
	return id
}

// NewKinveyEntity ...
func (m KinveyEntityModule) NewKinveyEntity(id string) KinveyEntity {
	entity := KinveyEntity{}

	if id != "" {
		entity.ID = String(id)
	} else {
		entity.ID = m.createEntityID()
	}

	t := time.Now()

	entity.ACL = &AccessControlList{
		Creator: String(m.environmentID),
	}

	entity.KMD = &KinveyMetadata{
		EntityCreatedTime: String(t.String()),
		LastModifiedTime:  String(t.String()),
	}

	return entity
}

// IsKinveyEntity ...
func (m KinveyEntityModule) IsKinveyEntity(testObject interface{}) bool {
	return true
}

// KinveyEntity ...
type KinveyEntity struct {
	ID  *string            `json:"_id,omitempty"`
	ACL *AccessControlList `json:"_acl,omitempty"`
	KMD *KinveyMetadata    `json:"_kmd,omitempty"`
}

// GetID ...
func (e *KinveyEntity) GetID() *string {
	return e.ID
}

// KinveyMetadata ...
type KinveyMetadata struct {
	EntityCreatedTime *string `json:"ect,omitempty"`
	LastModifiedTime  *string `json:"lmt,omitempty"`
}

// AccessControlList ...
type AccessControlList struct {
	Creator          *string   `json:"creator,omitempty"`
	Readers          *[]string `json:"r,omitempty"`
	Writers          *[]string `json:"w,omitempty"`
	Groups           *Groups   `json:"groups,omitempty"`
	Roles            *Roles    `json:"roles,omitempty"`
	GloballyReadable *bool     `json:"gr,omitempty"`
	GloballyWritable *bool     `json:"gw,omitempty"`
}

// Groups ...
type Groups struct {
	Readers *[]string `json:"r,omitempty"`
	Writers *[]string `json:"w,omitempty"`
}

// Roles ...
type Roles struct {
	Readers  *[]string `json:"r,omitempty"`
	Writers  *[]string `json:"w,omitempty"`
	Updaters *[]string `json:"u,omitempty"`
	Deleters *[]string `json:"d,omitempty"`
}

// GetCreator ...
func (acl *AccessControlList) GetCreator() *string {
	return acl.Creator
}

// GetReaders ...
func (acl *AccessControlList) GetReaders() *[]string {
	return acl.Readers
}

// GetWriters ...
func (acl *AccessControlList) GetWriters() *[]string {
	return acl.Writers
}

// GetReaderGroups ...
func (acl *AccessControlList) GetReaderGroups() *[]string {
	if acl.Groups == nil {
		return nil
	}
	return acl.Groups.Readers
}

// GetWriterGroups ...
func (acl *AccessControlList) GetWriterGroups() *[]string {
	if acl.Groups == nil {
		return nil
	}
	return acl.Groups.Writers
}

// GetReaderRoles ...
func (acl *AccessControlList) GetReaderRoles() *[]string {
	if acl.Roles == nil {
		return nil
	}
	return acl.Roles.Readers
}

// GetUpdateRoles ...
func (acl *AccessControlList) GetUpdateRoles() *[]string {
	if acl.Roles == nil {
		return nil
	}
	return acl.Roles.Updaters
}

// GetDeleteRoles ...
func (acl *AccessControlList) GetDeleteRoles() *[]string {
	if acl.Roles == nil {
		return nil
	}
	return acl.Roles.Deleters
}

// AddReader ...
func (acl *AccessControlList) AddReader(userID string) *AccessControlList {
	if acl.Readers == nil {
		r := make([]string, 0)
		acl.Readers = &r
	}

	if !util.Contains(*acl.Readers, userID) {
		*acl.Readers = append(*acl.Readers, userID)
	}

	return acl
}

// AddWriter ...
func (acl *AccessControlList) AddWriter(userID string) *AccessControlList {
	if acl.Writers == nil {
		r := make([]string, 0)
		acl.Writers = &r
	}

	if !util.Contains(*acl.Writers, userID) {
		*acl.Writers = append(*acl.Writers, userID)
	}

	return acl
}

// AddReaderGroup ...
func (acl *AccessControlList) AddReaderGroup(groupID string) *AccessControlList {
	if acl.Groups == nil {
		acl.Groups = &Groups{}
	}

	if acl.Groups.Readers == nil {
		r := make([]string, 0)
		acl.Groups.Readers = &r
	}

	if !util.Contains(*acl.Groups.Readers, groupID) {
		*acl.Groups.Readers = append(*acl.Groups.Readers, groupID)
	}

	return acl
}

// AddWriterGroup ...
func (acl *AccessControlList) AddWriterGroup(groupID string) *AccessControlList {
	if acl.Groups == nil {
		acl.Groups = &Groups{}
	}

	if acl.Groups.Writers == nil {
		r := make([]string, 0)
		acl.Groups.Writers = &r
	}

	if !util.Contains(*acl.Groups.Writers, groupID) {
		*acl.Groups.Writers = append(*acl.Groups.Writers, groupID)
	}

	return acl
}

// AddReaderRole ...
func (acl *AccessControlList) AddReaderRole(roleID string) *AccessControlList {
	if acl.Roles == nil {
		acl.Roles = &Roles{}
	}

	if acl.Roles.Readers == nil {
		r := make([]string, 0)
		acl.Roles.Readers = &r
	}

	if !util.Contains(*acl.Roles.Readers, roleID) {
		*acl.Roles.Readers = append(*acl.Roles.Readers, roleID)
	}

	return acl
}

// AddUpdateRole ...
func (acl *AccessControlList) AddUpdateRole(roleID string) *AccessControlList {
	if acl.Roles == nil {
		acl.Roles = &Roles{}
	}

	if acl.Roles.Updaters == nil {
		r := make([]string, 0)
		acl.Roles.Updaters = &r
	}

	if !util.Contains(*acl.Roles.Updaters, roleID) {
		*acl.Roles.Updaters = append(*acl.Roles.Updaters, roleID)
	}

	return acl
}

// AddDeleteRole ...
func (acl *AccessControlList) AddDeleteRole(roleID string) *AccessControlList {
	if acl.Roles == nil {
		acl.Roles = &Roles{}
	}

	if acl.Roles.Deleters == nil {
		r := make([]string, 0)
		acl.Roles.Deleters = &r
	}

	if !util.Contains(*acl.Roles.Deleters, roleID) {
		*acl.Roles.Deleters = append(*acl.Roles.Deleters, roleID)
	}

	return acl
}

// RemoveReader ...
func (acl *AccessControlList) RemoveReader(userID string) *AccessControlList {
	if acl.Readers == nil {
		return acl
	}

	*acl.Readers = util.Delete(*acl.Readers, userID)

	return acl
}

// RemoveWriter ...
func (acl *AccessControlList) RemoveWriter(userID string) *AccessControlList {
	if acl.Writers == nil {
		return acl
	}

	*acl.Writers = util.Delete(*acl.Writers, userID)

	return acl
}

// RemoveReaderGroup ...
func (acl *AccessControlList) RemoveReaderGroup(groupID string) *AccessControlList {
	if acl.Groups == nil {
		return acl
	}

	if acl.Groups.Readers == nil {
		return acl
	}

	*acl.Groups.Readers = util.Delete(*acl.Groups.Readers, groupID)

	return acl
}

// RemoveWriterGroup ...
func (acl *AccessControlList) RemoveWriterGroup(groupID string) *AccessControlList {
	if acl.Groups == nil {
		return acl
	}

	if acl.Groups.Writers == nil {
		return acl
	}

	*acl.Groups.Writers = util.Delete(*acl.Groups.Writers, groupID)

	return acl
}

// RemoveReaderRole ...
func (acl *AccessControlList) RemoveReaderRole(roleID string) *AccessControlList {
	if acl.Roles == nil {
		return acl
	}

	if acl.Roles.Readers == nil {
		return acl
	}

	*acl.Roles.Readers = util.Delete(*acl.Roles.Readers, roleID)

	return acl
}

// RemoveUpdateRole ...
func (acl *AccessControlList) RemoveUpdateRole(roleID string) *AccessControlList {
	if acl.Roles == nil {
		return acl
	}

	if acl.Roles.Updaters == nil {
		return acl
	}

	*acl.Roles.Updaters = util.Delete(*acl.Roles.Updaters, roleID)

	return acl
}

// RemoveDeleteRole ...
func (acl *AccessControlList) RemoveDeleteRole(roleID string) *AccessControlList {
	if acl.Roles == nil {
		return acl
	}

	if acl.Roles.Deleters == nil {
		return acl
	}

	*acl.Roles.Deleters = util.Delete(*acl.Roles.Deleters, roleID)

	return acl
}

// GetGloballyReadable ...
func (acl *AccessControlList) GetGloballyReadable() bool {
	if acl.GloballyReadable == nil {
		return false
	}
	return *acl.GloballyReadable
}

// GetGloballyWritable ...
func (acl *AccessControlList) GetGloballyWritable() bool {
	if acl.GloballyWritable == nil {
		return false
	}
	return *acl.GloballyWritable
}

// SetGloballyReadable ...
func (acl *AccessControlList) SetGloballyReadable(gr bool) *AccessControlList {
	acl.GloballyReadable = Bool(gr)
	return acl
}

// SetGloballyWritable ...
func (acl *AccessControlList) SetGloballyWritable(gw bool) *AccessControlList {
	acl.GloballyWritable = Bool(gw)
	return acl
}

type objectIDGenerator struct {
	objectIDCounter uint32
	machineID       []byte
	processID       int
}

func newObjectIDGenerator() objectIDGenerator {
	g := objectIDGenerator{
		objectIDCounter: readRandomUint32(),
		machineID:       readMachineID(),
		processID:       os.Getpid(),
	}

	return g
}

// ObjectIDFromHex ...
func (g *objectIDGenerator) ObjectIDFromHex(s string) *string {
	d, err := hex.DecodeString(s)

	if err != nil || len(d) != 12 {
		panic(fmt.Sprintf("invalid input to ObjectIdHex: %q", s))
	}

	id := string(d)

	return &id
}

// IsObjectIDHex ...
func (g *objectIDGenerator) IsObjectIDHex(s string) bool {
	if len(s) != 24 {
		return false
	}

	_, err := hex.DecodeString(s)

	return err == nil
}

func readRandomUint32() uint32 {
	var b [4]byte

	_, err := io.ReadFull(rand.Reader, b[:])

	if err != nil {
		panic(fmt.Errorf("cannot read random object id: %v", err))
	}

	return uint32((uint32(b[0]) << 0) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24))
}

func readMachineID() []byte {
	var sum [3]byte

	id := sum[:]

	hostname, err1 := os.Hostname()

	if err1 != nil {
		_, err2 := io.ReadFull(rand.Reader, id)
		if err2 != nil {
			panic(fmt.Errorf("cannot get hostname: %v; %v", err1, err2))
		}
		return id
	}

	hw := md5.New()

	hw.Write([]byte(hostname))

	copy(id, hw.Sum(nil))

	return id
}

// NewObjectID ...
func (g *objectIDGenerator) NewObjectID() *string {
	var b [12]byte

	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))

	b[4] = g.machineID[0]
	b[5] = g.machineID[1]
	b[6] = g.machineID[2]

	b[7] = byte(g.processID >> 8)
	b[8] = byte(g.processID)

	i := atomic.AddUint32(&g.objectIDCounter, 1)

	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)

	id := string(b[:])

	return &id
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// Int64 is a helper routine that allocates a new int64 value
// to store v and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }
