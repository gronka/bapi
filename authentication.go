package main

import (
	"github.com/gocql/gocql"
)

func (gibs *Gibs) isEmployee() bool {
	if gibs.UserUuid == AdminUuidBytes {
		return true
	}
	return false
}

func (gibs *Gibs) TacRequesteeCanViewTacs(requesteeUuid gocql.UUID) bool {
	if gibs.isEmployee() {
		return true
	}
	if gibs.UserUuid == requesteeUuid {
		return true
	}
	return false
}

func authWriteRequestFilledTime(gibs Gibs, ownerUuid gocql.UUID) bool {
	if gibs.isEmployee() {
		return true
	}
	return true
}

func authPlannerOwner(gibs Gibs, ownerUuid gocql.UUID) bool {
	if gibs.isEmployee() {
		return true
	}
	return gibs.simpleAuth(ownerUuid)
}

func authPlannerOwnerOrRequestee(gibs Gibs, ownerUuid, requesteeUuid gocql.UUID) bool {
	if gibs.isEmployee() {
		return true
	}
	return gibs.simpleAuthForTwo(ownerUuid, requesteeUuid)
}

func (gibs *Gibs) simpleAuth(userUuid gocql.UUID) bool {
	if gibs.isEmployee() || gibs.UserUuid == userUuid {
		return true
	}
	return false
}

func (gibs *Gibs) simpleAuthForTwo(userUuid01, userUuid02 gocql.UUID) bool {
	if gibs.isEmployee() || gibs.UserUuid == userUuid01 || gibs.UserUuid == userUuid02 {
		return true
	}
	return false
}

func (gibs *Gibs) authUserCanEditUser(userUuidToUpdate gocql.UUID) bool {
	if gibs.isEmployee() {
		return true
	}
	if gibs.UserUuid == userUuidToUpdate {
		return true
	}
	return false
}

func (gibs *Gibs) authUserCanAdminEvent(eventUuid gocql.UUID) bool {
	if gibs.isEmployee() {
		return true
	}
	if gibs.isEventAdmin(eventUuid) {
		return true
	}
	return false
}
