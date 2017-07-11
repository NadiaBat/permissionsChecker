package rbac

import (
	"github.com/NadiaBat/permissionsChecker/storage"
	"sync"
)

type Permission struct {
	UserId     int
	ActionName string
	HasAccess  bool
}

type Permissions []Permission

type Checker struct {
	wg          sync.WaitGroup
	permissions Permissions
}

func BulkCheck(userId int, actions []string) Permissions {
	checker := new(Checker)

	checker.permissions = make(Permissions, len(actions))
	checker.wg = sync.WaitGroup{}

	for _, action := range actions {
		checker.wg.Add(1)

		permission := Permission{UserId: userId, ActionName: action}
		go func(permissions *Permission) {
			permission.HasAccess = checkAccess(permission.UserId, permission.ActionName)
			checker.permissions = append(checker.permissions, permission)
			checker.wg.Done()
		}(&permission)
	}

	checker.wg.Wait()

	return checker.permissions
}

func checkAccess(userId int, actionName string) bool {
	// implement checking logic
	allAssignments := storage.GetAllAssignments(true)
	return userId < 2000 && len(actionName) > 4 && len(allAssignments) > 1
}
