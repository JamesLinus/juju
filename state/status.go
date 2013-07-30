// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"fmt"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/txn"

	"launchpad.net/juju-core/errors"
	"launchpad.net/juju-core/state/api/params"
)

// statusDoc represents a entity status in Mongodb.  The implicit
// _id field is explicitly set to the global key of the associated
// entity in the document's creation transaction, but omitted to allow
// direct use of the document in both create and update transactions.
type statusDoc struct {
	Status     params.Status
	StatusInfo string
}

// validate returns an error if the statusDoc is not valid.
func (doc statusDoc) validate() error {
	if !doc.Status.Valid() {
		return fmt.Errorf("cannot set invalid status %q", doc.Status)
	}
	switch doc.Status {
	case params.StatusPending:
		return fmt.Errorf("cannot set status %q", doc.Status)
	case params.StatusError:
		if doc.StatusInfo == "" {
			return fmt.Errorf("cannot set status %q without info", doc.Status)
		}
	}
	return nil
}

// getStatus retrieves the status document associated with the given
// globalKey and copies it to outStatusDoc, which needs to be created
// by the caller before.
func getStatus(st *State, globalKey string) (statusDoc, error) {
	var doc statusDoc
	err := st.statuses.FindId(globalKey).One(&doc)
	if err == mgo.ErrNotFound {
		return statusDoc{}, errors.NotFoundf("status")
	}
	if err != nil {
		return statusDoc{}, fmt.Errorf("cannot get status %q: %v", globalKey, err)
	}
	return doc, nil
}

// createStatusOp returns the operation needed to create the given
// status document associated with the given globalKey.
func createStatusOp(st *State, globalKey string, doc statusDoc) txn.Op {
	return txn.Op{
		C:      st.statuses.Name,
		Id:     globalKey,
		Assert: txn.DocMissing,
		Insert: doc,
	}
}

// updateStatusOp returns the operations needed to update the given
// status document associated with the given globalKey.
func updateStatusOp(st *State, globalKey string, doc statusDoc) txn.Op {
	return txn.Op{
		C:      st.statuses.Name,
		Id:     globalKey,
		Assert: txn.DocExists,
		Update: D{{"$set", doc}},
	}
}

// removeStatusOp returns the operation needed to remove the status
// document associated with the given globalKey.
func removeStatusOp(st *State, globalKey string) txn.Op {
	return txn.Op{
		C:      st.statuses.Name,
		Id:     globalKey,
		Remove: true,
	}
}
