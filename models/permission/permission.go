package permission

import "sync"

type Action struct {
	Name string
	//description string
	//data string // find, how it is used (from auth_item table)
}

type Assignment struct {
	UserId    int
	Action    Action
	HasAccess bool
}

func BulkCheck(wg *sync.WaitGroup, assignments []*Assignment) {
	for _, assignment := range assignments {
		go setAccess(assignment, wg)
	}
}

func setAccess(assignment *Assignment, wg *sync.WaitGroup) {
	defer wg.Done()

	assignment.HasAccess = assignment.UserId > 3456 && len(assignment.Action.Name) < 10
	// implement access checking
	// find in db (see in permissions_manager)
}
