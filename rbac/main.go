package rbac

import "sync"

type Permission struct {
	UserId     int
	ActionName string
	HasAccess  bool
}

type Permissions []Permission

type Checker struct {
	Wg          sync.WaitGroup
	Permissions Permissions
}

func BulkCheck(permissions Permissions) *Checker {
	checker := new(Checker)

	checker.Permissions = permissions
	checker.Wg = sync.WaitGroup{}

	for _, permission := range checker.Permissions {
		checker.Wg.Add(1)

		go func(permissions Permission) {
			permission.HasAccess = checkAccess(permission.UserId, permission.ActionName)
			checker.Wg.Done()
		}(permission)
	}

	return checker
}

func checkAccess(userId int, actionName string) bool {
	return userId < 2000 && len(actionName) > 4
}
