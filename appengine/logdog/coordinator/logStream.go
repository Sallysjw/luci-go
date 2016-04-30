// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package coordinator

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	ds "github.com/luci/gae/service/datastore"
	"github.com/luci/luci-go/common/logdog/types"
	"github.com/luci/luci-go/common/proto/logdog/logpb"
)

// currentLogStreamSchema is the current schema version of the LogStream.
// Changes that are not backward-compatible should update this field so
// migration logic and scripts can translate appropriately.
const currentLogStreamSchema = "1"

// LogStreamState is the archival state of the log stream.
type LogStreamState int

const (
	// LSStreaming indicates that the log stream is still streaming. This implies
	// that no terminal index has been identified yet.
	LSStreaming LogStreamState = iota
	// LSArchiveTasked indicates that the log stream has had an archival task
	// generated for it and is awaiting archival.
	LSArchiveTasked
	// LSArchived indicates that the log stream has been successfully archived.
	LSArchived
)

func (s LogStreamState) String() string {
	switch s {
	case LSStreaming:
		return "STREAMING"
	case LSArchiveTasked:
		return "ARCHIVE_TASKED"
	case LSArchived:
		return "ARCHIVED"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", s)
	}
}

// Archived returns true if this LogStreamState represents a finished archival.
func (s LogStreamState) Archived() bool {
	return s >= LSArchived
}

// LogStream is the primary datastore model containing information and state of
// an individual log stream.
//
// This structure contains the standard queryable fields, and is the source of
// truth for log stream state. Writes to LogStream should be done via Put, which
// will ensure that the LogStream's related query objects are kept in sync.
//
// This structure has additional datastore fields imposed by the
// PropertyLoadSaver. These fields enable querying against some of the complex
// data types:
//	  - _C breaks the the Prefix and Name fields into positionally-queryable
//	    entries. It is used to build globbing queries.
//
//	    It is composed of entries detailing (T is the path [P]refix or [N]ame):
//	    - TF:n:value ("T" has a component, "value", at index "n").
//	    - TR:n:value ("T" has a component, "value", at reverse-index "n").
//	    - TC:count ("T" has "count" total elements).
//
//	    For example, the path "foo/bar/+/baz" would break into:
//	    ["PF:0:foo", "PF:1:bar", "PR:0:bar", "PR:1:foo", "PC:2", "NF:0:baz",
//	     "NR:0:baz", "NC:1"].
//
//	  - _Tags is a string slice containing:
//	    - KEY=[VALUE] key/value tags.
//	    - KEY key presence tags.
//
//	  - _Terminated is true if the LogStream has been terminated.
//	  - _Archived is true if the LogStream has been archived.
//
// Most of the values in QueryBase are static. Those that change can only be
// changed through service endpoint methods.
type LogStream struct {
	// HashID is the LogStream ID. It is generated from the stream's Prefix/Name
	// fields.
	HashID string `gae:"$id"`

	// Schema is the datastore schema version for this object. This can be used
	// to facilitate schema migrations.
	//
	// The current schema is currentLogStreamSchema.
	Schema string

	// Prefix is this log stream's prefix value. Log streams with the same prefix
	// are logically grouped.
	//
	// This value should not be changed once populated, as it will invalidate the
	// HashID.
	Prefix string
	// Name is the unique name of this log stream within the Prefix scope.
	//
	// This value should not be changed once populated, as it will invalidate the
	// HashID.
	Name string

	// State is the log stream's current state.
	State LogStreamState

	// Purged, if true, indicates that this log stream has been marked as purged.
	// Non-administrative queries and requests for this stream will operate as
	// if this entry doesn't exist.
	Purged bool

	// Secret is the Butler secret value for this stream.
	//
	// This value may only be returned to LogDog services; it is not user-visible.
	Secret []byte `gae:",noindex"`

	// Created is the time when this stream was created.
	Created time.Time
	// TerminatedTime is the Coordinator's record of when this log stream was
	// terminated.
	TerminatedTime time.Time `gae:",noindex"`
	// ArchivedTime is the Coordinator's record of when this log stream was
	// archived.
	ArchivedTime time.Time `gae:",noindex"`

	// ProtoVersion is the version string of the protobuf, as reported by the
	// Collector (and ultimately self-identified by the Butler).
	ProtoVersion string
	// Descriptor is the binary protobuf data LogStreamDescriptor.
	Descriptor []byte `gae:",noindex"`
	// ContentType is the MIME-style content type string for this stream.
	ContentType string
	// StreamType is the data type of the stream.
	StreamType logpb.StreamType
	// Timestamp is the Descriptor's recorded client-side timestamp.
	Timestamp time.Time

	// Tags is a set of arbitrary key/value tags associated with this stream. Tags
	// can be queried against.
	//
	// The serialization/deserialization is handled manually in order to enable
	// key/value queries.
	Tags TagMap `gae:"-"`

	// Source is the set of source strings sent by the Butler.
	Source []string

	// TerminalIndex is the index of the last log entry in the stream.
	//
	// If this is <0, the log stream is either still streaming or has been
	// archived with no log entries.
	TerminalIndex int64 `gae:",noindex"`

	// ArchiveLogEntryCount is the number of LogEntry records that were archived
	// for this log stream.
	//
	// This is valid only if the log stream is Archived.
	ArchiveLogEntryCount int64 `gae:",noindex"`
	// ArchivalKey is the archival key for this log stream. This is used to
	// differentiate the real archival request from those that were dispatched,
	// but that ultimately failed to update state.
	ArchivalKey []byte `gae:",noindex"`

	// ArchiveIndexURL is the Google Storage URL where the log stream's index is
	// archived.
	ArchiveIndexURL string `gae:",noindex"`
	// ArchiveIndexSize is the size, in bytes, of the archived Index. It will be
	// zero if the file is not archived.
	ArchiveIndexSize int64 `gae:",noindex"`
	// ArchiveStreamURL is the Google Storage URL where the log stream's raw
	// stream data is archived. If this is not empty, the log stream is considered
	// archived.
	ArchiveStreamURL string `gae:",noindex"`
	// ArchiveStreamSize is the size, in bytes, of the archived stream. It will be
	// zero if the file is not archived.
	ArchiveStreamSize int64 `gae:",noindex"`
	// ArchiveDataURL is the Google Storage URL where the log stream's assembled
	// data is archived. If this is not empty, the log stream is considered
	// archived.
	ArchiveDataURL string `gae:",noindex"`
	// ArchiveDataSize is the size, in bytes, of the archived data. It will be
	// zero if the file is not archived.
	ArchiveDataSize int64 `gae:",noindex"`

	// extra causes datastore to ignore unrecognized fields and strip them in
	// future writes.
	extra ds.PropertyMap `gae:"-,extra"`

	// noDSValidate is a testing parameter to instruct the LogStream not to
	// validate before reading/writing to datastore. It can be controlled by
	// calling SetDSValidate().
	noDSValidate bool
}

var _ interface {
	ds.PropertyLoadSaver
} = (*LogStream)(nil)

// NewLogStream returns a LogStream instance with its ID field initialized based
// on the supplied path.
//
// The supplied value is a LogDog stream path or a hash of the LogDog stream
// path.
func NewLogStream(value string) (*LogStream, error) {
	path := types.StreamPath(value)
	if err := path.Validate(); err != nil {
		// If it's not a path, see if it's a SHA256 sum.
		hash, hashErr := normalizeHash(value)
		if hashErr != nil {
			return nil, fmt.Errorf("invalid path (%s) and hash (%s)", err, hashErr)
		}

		// Load this LogStream with its SHA256 hash directly. This stream will not
		// have its Prefix/Name fields populated until it's loaded from datastore.
		return LogStreamFromID(hash), nil
	}

	return LogStreamFromPath(path), nil
}

// LogStreamFromID returns an empty LogStream instance with a known hash ID.
func LogStreamFromID(hashID string) *LogStream {
	return &LogStream{
		HashID: hashID,
	}
}

// LogStreamFromPath returns an empty LogStream instance initialized from a
// known path value.
//
// The supplied path is assumed to be valid and is not checked.
func LogStreamFromPath(path types.StreamPath) *LogStream {
	// Load the prefix/name fields into the log stream.
	prefix, name := path.Split()
	ls := LogStream{
		Prefix: string(prefix),
		Name:   string(name),
	}
	ls.recalculateHashID()
	return &ls
}

// Path returns the LogDog path for this log stream.
func (s *LogStream) Path() types.StreamPath {
	return types.StreamName(s.Prefix).Join(types.StreamName(s.Name))
}

// Load implements ds.PropertyLoadSaver.
func (s *LogStream) Load(pmap ds.PropertyMap) error {
	// Handle custom properties. Consume them before using the default
	// PropertyLoadSaver.
	for k, v := range pmap {
		if !strings.HasPrefix(k, "_") {
			continue
		}

		switch k {
		case "_Tags":
			// Load the tag map. Ignore errors.
			tm, _ := tagMapFromProperties(v)
			s.Tags = tm
		}
		delete(pmap, k)
	}

	if err := ds.GetPLS(s).Load(pmap); err != nil {
		return err
	}

	// Migrate schema (if needed), then validate.
	if err := s.migrateSchema(); err != nil {
		return err
	}

	// Validate the log stream. Don't enforce HashID correctness, since
	// datastore hasn't populated that field yet.
	if !s.noDSValidate {
		if err := s.validateImpl(false); err != nil {
			return err
		}
	}
	return nil
}

// Save implements ds.PropertyLoadSaver.
func (s *LogStream) Save(withMeta bool) (ds.PropertyMap, error) {
	if !s.noDSValidate {
		if err := s.validateImpl(true); err != nil {
			return nil, err
		}
	}
	s.Schema = currentLogStreamSchema

	// Save default struct fields.
	pmap, err := ds.GetPLS(s).Save(withMeta)
	if err != nil {
		return nil, err
	}

	// Encode _Tags.
	pmap["_Tags"], err = s.Tags.toProperties()
	if err != nil {
		return nil, fmt.Errorf("failed to encode tags: %v", err)
	}

	// Generate our path components, "_C".
	pmap["_C"] = generatePathComponents(s.Prefix, s.Name)

	// Add our derived statuses.
	pmap["_Terminated"] = []ds.Property{ds.MkProperty(s.Terminated())}
	pmap["_Archived"] = []ds.Property{ds.MkProperty(s.Archived())}

	return pmap, nil
}

// recalculateHashID calculates the log stream's hash ID from its Prefix/Name
// fields, which must be populated else this function will panic.
//
// The value is loaded into its HashID field.
func (s *LogStream) recalculateHashID() {
	s.HashID = s.getHashID()
}

// recalculateHashID calculates the log stream's hash ID from its Prefix/Name
// fields, which must be populated else this function will panic.
func (s *LogStream) getHashID() string {
	hash := sha256.Sum256([]byte(s.Path()))
	return hex.EncodeToString(hash[:])
}

// Validate evaluates the state and data contents of the LogStream and returns
// an error if it is invalid.
func (s *LogStream) Validate() error {
	return s.validateImpl(true)
}

func (s *LogStream) validateImpl(enforceHashID bool) error {
	if enforceHashID {
		// Make sure our Prefix and Name match the Hash ID.
		if hid := s.getHashID(); hid != s.HashID {
			return fmt.Errorf("hash IDs don't match (%q != %q)", hid, s.HashID)
		}
	}

	if err := types.StreamName(s.Prefix).Validate(); err != nil {
		return fmt.Errorf("invalid prefix: %s", err)
	}
	if err := types.StreamName(s.Name).Validate(); err != nil {
		return fmt.Errorf("invalid name: %s", err)
	}
	if err := types.PrefixSecret(s.Secret).Validate(); err != nil {
		return fmt.Errorf("invalid prefix secret: %s", err)
	}
	if s.ContentType == "" {
		return errors.New("empty content type")
	}
	if s.Created.IsZero() {
		return errors.New("created time is not set")
	}

	if s.Terminated() && s.TerminatedTime.IsZero() {
		return errors.New("log stream is terminated, but missing terminated time")
	}
	if s.Archived() && s.ArchivedTime.IsZero() {
		return errors.New("log stream is archived, but missing archived time")
	}

	switch s.StreamType {
	case logpb.StreamType_TEXT, logpb.StreamType_BINARY, logpb.StreamType_DATAGRAM:
		break

	default:
		return fmt.Errorf("unsupported stream type: %v", s.StreamType)
	}

	for k, v := range s.Tags {
		if err := types.ValidateTag(k, v); err != nil {
			return fmt.Errorf("invalid tag [%s]: %s", k, err)
		}
	}

	// Ensure that our Descriptor can be unmarshalled.
	if _, err := s.DescriptorValue(); err != nil {
		return fmt.Errorf("could not unmarshal descriptor: %v", err)
	}
	return nil
}

// DescriptorValue returns the unmarshalled Descriptor field protobuf.
func (s *LogStream) DescriptorValue() (*logpb.LogStreamDescriptor, error) {
	pb := logpb.LogStreamDescriptor{}
	if err := proto.Unmarshal(s.Descriptor, &pb); err != nil {
		return nil, err
	}
	return &pb, nil
}

// Terminated returns true if this stream has been terminated.
func (s *LogStream) Terminated() bool {
	if s.Archived() {
		return true
	}
	return s.TerminalIndex >= 0
}

// Archived returns true if this stream has been archived.
func (s *LogStream) Archived() bool {
	return s.State.Archived()
}

// ArchiveComplete returns true if this stream has been archived and all of its
// log entries were present.
func (s *LogStream) ArchiveComplete() bool {
	return (s.Archived() && s.ArchiveLogEntryCount == (s.TerminalIndex+1))
}

// LoadDescriptor loads the fields in the log stream descriptor into this
// LogStream entry. These fields are:
//   - Prefix
//   - Name
//   - ContentType
//   - StreamType
//   - Descriptor
//   - Timestamp
//   - Tags
func (s *LogStream) LoadDescriptor(desc *logpb.LogStreamDescriptor) error {
	// If the descriptor's Prefix/Name don't match ours, refuse to load it.
	if desc.Prefix != s.Prefix {
		return fmt.Errorf("prefixes don't match (%q != %q)", desc.Prefix, s.Prefix)
	}
	if desc.Name != s.Name {
		return fmt.Errorf("names don't match (%q != %q)", desc.Name, s.Name)
	}

	if err := desc.Validate(true); err != nil {
		return fmt.Errorf("invalid descriptor: %v", err)
	}

	pb, err := proto.Marshal(desc)
	if err != nil {
		return fmt.Errorf("failed to marshal descriptor: %v", err)
	}

	s.Prefix = desc.Prefix
	s.Name = desc.Name
	s.ContentType = desc.ContentType
	s.StreamType = desc.StreamType
	s.Descriptor = pb

	// We know that the timestamp is valid b/c it's checked in ValidateDescriptor.
	if ts := desc.Timestamp; ts != nil {
		s.Timestamp = ds.RoundTime(ts.Time().UTC())
	}

	// Note: tag content was validated via ValidateDescriptor.
	s.Tags = TagMap(desc.Tags)
	return nil
}

// DescriptorProto unmarshals a LogStreamDescriptor from the stream's Descriptor
// field. It will return an error if the unmarshalling fails.
func (s *LogStream) DescriptorProto() (*logpb.LogStreamDescriptor, error) {
	desc := logpb.LogStreamDescriptor{}
	if err := proto.Unmarshal(s.Descriptor, &desc); err != nil {
		return nil, err
	}
	return &desc, nil
}

// SetDSValidate controls whether this LogStream is validated prior to being
// read from or written to datastore.
//
// This is a testing parameter, and should NOT be used in production code.
func (s *LogStream) SetDSValidate(v bool) {
	s.noDSValidate = !v
}

// normalizeHash takes a SHA256 hexadecimal string as input. It validates that
// it is a valid SHA256 hash and, if so, returns a normalized version that can
// be used as a log stream key.
func normalizeHash(v string) (string, error) {
	if decodeSize := hex.DecodedLen(len(v)); decodeSize != sha256.Size {
		return "", fmt.Errorf("invalid SHA256 hash size (%d != %d)", decodeSize, sha256.Size)
	}
	b, err := hex.DecodeString(v)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// generatePathComponents generates the "_C" property path components for path
// glob querying.
//
// See the comment on LogStream for more infromation.
func generatePathComponents(prefix, name string) []ds.Property {
	ps, ns := types.StreamName(prefix).Segments(), types.StreamName(name).Segments()

	// Allocate our components array. For each component, there are two entries
	// (forward and reverse), as well as one count entry per component type.
	c := make([]ds.Property, 0, (len(ps)+len(ns)+1)*2)

	gen := func(b string, segs []string) {
		// Generate count component (PC:4).
		c = append(c, ds.MkProperty(fmt.Sprintf("%sC:%d", b, len(segs))))

		// Generate forward and reverse components.
		for i, s := range segs {
			c = append(c,
				// Forward (PF:0:foo).
				ds.MkProperty(fmt.Sprintf("%sF:%d:%s", b, i, s)),
				// Reverse (PR:3:foo)
				ds.MkProperty(fmt.Sprintf("%sR:%d:%s", b, len(segs)-i-1, s)),
			)
		}
	}
	gen("P", ps)
	gen("N", ns)
	return c
}

func addComponentFilter(q *ds.Query, full, f, base string, value types.StreamName) (*ds.Query, error) {
	segments := value.Segments()
	if len(segments) == 0 {
		// All-component query; impose no constraints.
		return q, nil
	}

	// Profile the string. If it doesn't have any glob characters in it,
	// fully constrain the base field.
	hasGlob := false
	for i, seg := range segments {
		switch seg {
		case "*", "**":
			hasGlob = true

		default:
			// Regular segment. Assert that it is valid.
			if err := types.StreamName(seg).Validate(); err != nil {
				return nil, fmt.Errorf("invalid %s component at index %d (%s): %s", full, i, seg, err)
			}
		}
	}
	if !hasGlob {
		// Direct field (full) query.
		return q.Eq(full, string(value)), nil
	}

	// Add specific field constraints for each non-glob segment.
	greedy := false
	rstack := []string(nil)
	for i, seg := range segments {
		switch seg {
		case "*":
			// Skip asserting this segment.
			if greedy {
				// Add a placeholder (e.g., .../**/a/*/b, placeholder ensures "a" gets
				// position -2 instead of -1.
				//
				// Note that "" can never be a segment, because we validate each
				// non-glob path segment and "" is not a valid stream name component.
				rstack = append(rstack, "")
			}
			continue

		case "**":
			if greedy {
				return nil, fmt.Errorf("cannot have more than one greedy glob")
			}

			// Mark that we're greedy, and skip asserting this segment.
			greedy = true
			continue

		default:
			if greedy {
				// Add this to our reverse stack. We'll query the reverse field from
				// this stack after we know how many elements are ultimately in it.
				rstack = append(rstack, seg)
			} else {
				q = q.Eq(f, fmt.Sprintf("%sF:%d:%s", base, i, seg))
			}
		}
	}

	// Add the reverse stack to pin the elements at the end of the path name
	// (e.g., a/b/**/c/d, stack will be {c, d}, need to map to {r.0=d, r.1=c}.
	for i, seg := range rstack {
		if seg == "" {
			// Placeholder, skip.
			continue
		}
		q = q.Eq(f, fmt.Sprintf("%sR:%d:%s", base, len(rstack)-i-1, seg))
	}

	// If we're not greedy, fix this size of this component.
	if !greedy {
		q = q.Eq(f, fmt.Sprintf("%sC:%d", base, len(segments)))
	}
	return q, nil
}

// AddLogStreamPathFilter constructs a compiled LogStreamPathQuery. It will
// return an error if the supllied query string describes an invalid query.
func AddLogStreamPathFilter(q *ds.Query, path string) (*ds.Query, error) {
	prefix, name := types.StreamPath(path).Split()

	err := error(nil)
	q, err = addComponentFilter(q, "Prefix", "_C", "P", prefix)
	if err != nil {
		return nil, err
	}
	q, err = addComponentFilter(q, "Name", "_C", "N", name)
	if err != nil {
		return nil, err
	}
	return q, nil
}

// AddLogStreamTerminatedFilter returns a derived query that asserts that a log
// stream has been terminated.
func AddLogStreamTerminatedFilter(q *ds.Query, v bool) *ds.Query {
	return q.Eq("_Terminated", v)
}

// AddLogStreamArchivedFilter returns a derived query that asserts that a log
// stream has been archived.
func AddLogStreamArchivedFilter(q *ds.Query, v bool) *ds.Query {
	return q.Eq("_Archived", v)
}

// AddLogStreamPurgedFilter returns a derived query that asserts that a log
// stream has been archived.
func AddLogStreamPurgedFilter(q *ds.Query, v bool) *ds.Query {
	return q.Eq("Purged", v)
}

// AddOlderFilter adds a filter to queries that restricts them to results that
// were created before the supplied time.
func AddOlderFilter(q *ds.Query, t time.Time) *ds.Query {
	return q.Lt("Created", t.UTC()).Order("-Created")
}

// AddNewerFilter adds a filter to queries that restricts them to results that
// were created after the supplied time.
func AddNewerFilter(q *ds.Query, t time.Time) *ds.Query {
	return q.Gt("Created", t.UTC()).Order("-Created")
}
