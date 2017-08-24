package main

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"sync"
)

type Permission struct {
	UserId     int
	ActionName string
	HasAccess  bool
}

type checkingParams struct {
	userId       int
	region       int
	project      int
	isCommercial bool
}

type Permissions []*Permission

type Checker struct {
	sync.WaitGroup
	permissions Permissions
}

func BulkCheck(userId int, actions []string, additionalParams map[string]string) (Permissions, error) {
	checker := &Checker{
		permissions: make(Permissions, len(actions)),
	}

	params, err := getCheckingParams(userId, additionalParams)
	if err != nil {
		return nil, errors.Wrap(err, "Не удалось выполнить проверку.")
	}

	// @todo probably, should not have async checking
	// only for several users (unlikely case)
	var errs []error
	for _, action := range actions {
		checker.Add(1)

		permission := &Permission{UserId: userId, ActionName: action}
		go func(permission *Permission, errs *[]error) {
			var checkingErr error
			permission.HasAccess, checkingErr = checkAccess(userId, permission.ActionName, params)

			if checkingErr != nil {
				checkingErr = errors.Wrapf(
					checkingErr,
					"Can`t execute checking for userId=%d, actionName=%s",
					permission.UserId,
					permission.ActionName,
				)

				*errs = append(*errs, checkingErr)
			}
			checker.permissions = append(checker.permissions, permission)

			checker.Done()
		}(permission, &errs)
	}

	checker.Wait()
	if len(errs) > 0 {
		return checker.permissions, errs[0]
	}

	return checker.permissions, nil
}

func getCheckingParams(userId int, additionalParams map[string]string) (*checkingParams, error) {
	params := checkingParams{userId: userId, region: 0, project: 0}
	var err error
	for name, value := range additionalParams {
		switch name {
		case "region":
			params.region, err = strconv.Atoi(value)
		case "project":
			params.project, err = strconv.Atoi(value)
		default:
			continue
		}
	}

	if err != nil {
		return nil, err
	}

	return &params, nil
}

func checkAccess(userId int, actionName string, params *checkingParams) (bool, error) {
	userAssignments, err := getUserAssignments(userId)
	if err != nil {
		return false, errors.Wrap(err, "Can`t get user assignments.")
	}

	return checkAccessRecursive(userId, actionName, params, userAssignments), nil
}

func getUserAssignments(userId int) (map[string]Assignment, error) {
	allAssignments := GetAllAssignments()
	userAssignments, ok := allAssignments[userId]
	if !ok {
		return nil, errors.New("User assignments doesn`t exists.")
	}

	return userAssignments.Items, nil

}

func checkAccessRecursive(
	userId int, itemName string, params *checkingParams, assignments map[string]Assignment,
) bool {
	permissionItem, err := getPermissionItem(itemName)
	if err != nil {
		return false
	}

	if !executeRule(permissionItem.Rule, params, permissionItem.Data) {
		return false
	}

	itemAssignment, ok := assignments[itemName]
	if ok {
		if executeRule(itemAssignment.Rule, params, itemAssignment.Data) {
			return true
		}
	}

	parents, err := getParents(itemName)
	if err != nil {
		return false
	}

	for _, parentItem := range parents {
		if checkAccessRecursive(userId, parentItem, params, assignments) {
			return true
		}
	}

	return false
}

func getPermissionItem(name string) (*PermissionItem, error) {
	allPermissionItems := GetAllPermissionItems()
	permissionItem, ok := allPermissionItems[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("There is no permission item %s", name))
	}

	return &permissionItem, nil
}

func getParents(childName string) (ItemParents, error) {
	allParents := GetAllParents()
	itemParents, ok := allParents[childName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("There is no parents for item %s", childName))
	}

	return itemParents, nil
}

func executeRule(rule string, params *checkingParams, data string) bool {
	if len(rule) == 0 {
		return true
	}

	// @TODO there is only 1 rule (isCommercial = 1)
	if rule != "News_Permissions_Rules::inArray" || len(data) == 0 {
		return false
	}

	return params.isCommercial
}

/*
Варианты параметров для правила inArray
Написать тесты
'a:1:{s:9:"paramsKey";s:4:"name";}',
'a:2:{s:9:"paramsKey";s:12:"isCommercial";s:4:"data";a:1:{i:0;i:0;}}',
'a:2:{s:9:"paramsKey";s:12:"isCommercial";s:4:"data";a:1:{i:0;i:1;}}',
'a:2:{s:9:"paramsKey";s:3:"pid";s:4:"data";a:1:{i:0;s:6:"369550";}}',
'a:2:{s:9:"paramsKey";s:3:"pid";s:4:"data";a:1:{i:0;s:8:"14338667";}}',
'a:2:{s:9:"paramsKey";s:3:"pid";s:4:"data";a:1:{i:0;s:8:"14338727";}}',
'a:2:{s:9:"paramsKey";s:3:"pid";s:4:"data";a:1:{i:0;s:8:"14338747";}}',
'a:2:{s:9:"paramsKey";s:3:"pid";s:4:"data";a:1:{i:0;s:9:"145919821";}}',
'a:2:{s:9:"paramsKey";s:3:"pid";s:4:"data";a:1:{i:0;s:9:"152645602";}}',
'a:2:{s:9:"paramsKey";s:3:"pid";s:4:"data";a:1:{i:0;s:9:"200132743";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";N;}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:0:{}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:10:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:15:"flags[yarabota]";i:3;s:4:"type";i:4;s:9:"mainPhoto";i:5;s:6:"images";i:6;s:10:"sourceName";i:7;s:5:"theme";i:8;s:7:"authors";i:9;s:17:"isCommentsAllowed";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:10:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:4:"type";i:3;s:6:"images";i:4;s:10:"sourceName";i:5;s:5:"theme";i:6;s:5:"links";i:7;s:7:"authors";i:8;s:17:"isCommentsAllowed";i:9;s:4:"tags";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:11:{i:0;s:13:"flags[tomain]";i:1;s:12:"flags[torss]";i:2;s:5:"isBaa";i:3;s:9:"mainPhoto";i:4;s:10:"sourceName";i:5;s:8:"category";i:6;s:5:"theme";i:7;s:5:"links";i:8;s:7:"authors";i:9;s:4:"tags";i:10;s:9:"copyright";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:12:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:12:"flags[torss]";i:3;s:4:"type";i:4;s:6:"images";i:5;s:10:"sourceName";i:6;s:5:"theme";i:7;s:5:"links";i:8;s:7:"authors";i:9;s:17:"isCommentsAllowed";i:10;s:4:"tags";i:11;s:8:"crmOrder";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:12:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:5:"isBaa";i:3;s:4:"type";i:4;s:9:"mainPhoto";i:5;s:21:"socialBackgroundImage";i:6;s:6:"images";i:7;s:10:"sourceName";i:8;s:5:"theme";i:9;s:7:"authors";i:10;s:17:"isCommentsAllowed";i:11;s:8:"crmOrder";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:12:{i:0;s:9:"subheader";i:1;s:12:"imagesAuthor";i:2;s:13:"flags[tomain]";i:3;s:9:"mainPhoto";i:4;s:6:"images";i:5;s:6:"videos";i:6;s:8:"category";i:7;s:8:"keywords";i:8;s:5:"polls";i:9;s:5:"links";i:10;s:7:"authors";i:11;s:17:"isCommentsAllowed";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:12:{i:0;s:9:"subheader";i:1;s:12:"imagesAuthor";i:2;s:13:"flags[tomain]";i:3;s:9:"mainPhoto";i:4;s:6:"images";i:5;s:8:"category";i:6;s:8:"keywords";i:7;s:5:"polls";i:8;s:5:"links";i:9;s:7:"authors";i:10;s:17:"isCommentsAllowed";i:11;s:8:"crmOrder";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:13:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:10:"htmlButton";i:3;s:5:"feeds";i:4;s:4:"type";i:5;s:6:"images";i:6;s:6:"videos";i:7;s:10:"sourceName";i:8;s:5:"theme";i:9;s:5:"links";i:10;s:7:"authors";i:11;s:17:"isCommentsAllowed";i:12;s:4:"tags";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:14:{i:0;s:16:"32.autoBlock.fri";i:1;s:16:"32.autoBlock.mon";i:2;s:16:"32.autoBlock.sat";i:3;s:16:"32.autoBlock.sun";i:4;s:16:"32.autoBlock.thu";i:5;s:16:"32.autoBlock.tue";i:6;s:16:"32.autoBlock.wed";i:7;s:15:"32.domBlock.fri";i:8;s:15:"32.domBlock.mon";i:9;s:15:"32.domBlock.sat";i:10;s:15:"32.domBlock.sun";i:11;s:15:"32.domBlock.thu";i:12;s:15:"32.domBlock.tue";i:13;s:15:"32.domBlock.wed";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:14:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:12:"flags[torss]";i:3;s:4:"type";i:4;s:6:"images";i:5;s:6:"videos";i:6;s:10:"sourceName";i:7;s:8:"category";i:8;s:5:"theme";i:9;s:5:"links";i:10;s:7:"authors";i:11;s:17:"isCommentsAllowed";i:12;s:4:"tags";i:13;s:8:"crmOrder";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:14:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:5:"feeds";i:3;s:13:"flags[tomain]";i:4;s:5:"isBaa";i:5;s:4:"type";i:6;s:9:"mainPhoto";i:7;s:21:"socialBackgroundImage";i:8;s:6:"author";i:9;s:8:"category";i:10;s:5:"theme";i:11;s:7:"authors";i:12;s:17:"isCommentsAllowed";i:13;s:9:"customUrl";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:15:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:10:"htmlButton";i:3;s:13:"flags[tomain]";i:4;s:12:"flags[totop]";i:5;s:12:"flags[torss]";i:6;s:4:"type";i:7;s:9:"mainPhoto";i:8;s:6:"images";i:9;s:8:"category";i:10;s:5:"theme";i:11;s:5:"polls";i:12;s:7:"authors";i:13;s:17:"isCommentsAllowed";i:14;s:15:"commercialLabel";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:15:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:13:"flags[tomain]";i:3;s:12:"flags[totop]";i:4;s:12:"flags[torss]";i:5;s:4:"type";i:6;s:9:"mainPhoto";i:7;s:6:"images";i:8;s:10:"sourceName";i:9;s:5:"theme";i:10;s:5:"links";i:11;s:7:"authors";i:12;s:17:"isCommentsAllowed";i:13;s:4:"tags";i:14;s:9:"copyright";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:16:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:10:"htmlButton";i:3;s:13:"flags[tomain]";i:4;s:15:"flags[yarabota]";i:5;s:5:"isBaa";i:6;s:4:"type";i:7;s:9:"mainPhoto";i:8;s:6:"images";i:9;s:6:"videos";i:10;s:10:"sourceName";i:11;s:5:"theme";i:12;s:7:"authors";i:13;s:17:"isCommentsAllowed";i:14;s:9:"customUrl";i:15;s:9:"copyright";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:17:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:10:"htmlButton";i:3;s:13:"flags[tomain]";i:4;s:5:"isBaa";i:5;s:4:"type";i:6;s:9:"mainPhoto";i:7;s:6:"images";i:8;s:6:"videos";i:9;s:10:"sourceName";i:10;s:5:"theme";i:11;s:7:"authors";i:12;s:17:"isCommentsAllowed";i:13;s:15:"commercialLabel";i:14;s:9:"customUrl";i:15;s:9:"copyright";i:16;s:8:"crmOrder";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:17:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:13:"flags[tomain]";i:3;s:12:"flags[torss]";i:4;s:13:"flags[n1spam]";i:5;s:4:"type";i:6;s:9:"mainPhoto";i:7;s:6:"images";i:8;s:10:"sourceName";i:9;s:5:"theme";i:10;s:14:"headerKeywords";i:11;s:5:"links";i:12;s:7:"authors";i:13;s:17:"isCommentsAllowed";i:14;s:4:"tags";i:15;s:9:"customUrl";i:16;s:9:"copyright";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:17:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:13:"flags[tomain]";i:3;s:12:"flags[tovip]";i:4;s:4:"type";i:5;s:9:"mainPhoto";i:6;s:9:"mainVideo";i:7;s:6:"images";i:8;s:6:"videos";i:9;s:10:"sourceName";i:10;s:8:"category";i:11;s:5:"theme";i:12;s:5:"polls";i:13;s:17:"isCommentsAllowed";i:14;s:15:"commercialLabel";i:15;s:9:"customUrl";i:16;s:9:"copyright";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:18:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:13:"flags[tomain]";i:3;s:12:"flags[tovip]";i:4;s:5:"isBaa";i:5;s:4:"type";i:6;s:9:"mainPhoto";i:7;s:9:"mainVideo";i:8;s:6:"images";i:9;s:10:"sourceName";i:10;s:6:"author";i:11;s:8:"category";i:12;s:5:"theme";i:13;s:5:"polls";i:14;s:15:"commercialLabel";i:15;s:9:"customUrl";i:16;s:9:"copyright";i:17;s:8:"crmOrder";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:1:{i:0;s:13:"flags[n1spam]";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:1:{i:0;s:16:"flags[comPlace7]";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:1:{i:0;s:4:"lead";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:1:{i:0;s:7:"authors";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:21:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:12:"imagesAuthor";i:3;s:13:"flags[tomain]";i:4;s:16:"flags[fixonmain]";i:5;s:15:"flags[toslider]";i:6;s:15:"flags[photorep]";i:7;s:5:"isBaa";i:8;s:4:"type";i:9;s:9:"mainPhoto";i:10;s:6:"images";i:11;s:10:"sourceName";i:12;s:8:"category";i:13;s:5:"theme";i:14;s:5:"links";i:15;s:7:"authors";i:16;s:17:"isCommentsAllowed";i:17;s:15:"commercialLabel";i:18;s:9:"customUrl";i:19;s:9:"copyright";i:20;s:8:"crmOrder";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:22:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:12:"imagesAuthor";i:3;s:13:"flags[tomain]";i:4;s:12:"flags[totop]";i:5;s:12:"flags[torss]";i:6;s:15:"flags[toslider]";i:7;s:15:"flags[photorep]";i:8;s:4:"type";i:9;s:9:"mainPhoto";i:10;s:21:"socialBackgroundImage";i:11;s:27:"socialBackgroundImageHeader";i:12;s:6:"images";i:13;s:10:"sourceName";i:14;s:8:"category";i:15;s:5:"theme";i:16;s:14:"headerKeywords";i:17;s:5:"links";i:18;s:7:"authors";i:19;s:17:"isCommentsAllowed";i:20;s:9:"customUrl";i:21;s:9:"copyright";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:22:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:13:"flags[tomain]";i:3;s:12:"flags[totop]";i:4;s:12:"flags[torss]";i:5;s:4:"type";i:6;s:9:"mainPhoto";i:7;s:9:"mainVideo";i:8;s:6:"images";i:9;s:6:"videos";i:10;s:10:"sourceName";i:11;s:8:"category";i:12;s:5:"theme";i:13;s:14:"headerKeywords";i:14;s:5:"links";i:15;s:7:"authors";i:16;s:12:"goodNewsText";i:17;s:17:"isCommentsAllowed";i:18;s:4:"tags";i:19;s:15:"commercialLabel";i:20;s:9:"customUrl";i:21;s:9:"copyright";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:24:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:10:"htmlButton";i:3;s:5:"feeds";i:4;s:13:"flags[tomain]";i:5;s:12:"flags[totop]";i:6;s:12:"flags[torss]";i:7;s:12:"flags[tovip]";i:8;s:4:"type";i:9;s:9:"mainPhoto";i:10;s:9:"mainVideo";i:11;s:6:"images";i:12;s:6:"videos";i:13;s:10:"sourceName";i:14;s:8:"category";i:15;s:5:"theme";i:16;s:14:"headerKeywords";i:17;s:5:"links";i:18;s:7:"authors";i:19;s:12:"goodNewsText";i:20;s:17:"isCommentsAllowed";i:21;s:4:"tags";i:22;s:15:"commercialLabel";i:23;s:9:"copyright";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:25:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:10:"htmlButton";i:3;s:13:"flags[tomain]";i:4;s:12:"flags[totop]";i:5;s:12:"flags[torss]";i:6;s:12:"flags[tovip]";i:7;s:4:"type";i:8;s:9:"mainPhoto";i:9;s:21:"socialBackgroundImage";i:10;s:27:"socialBackgroundImageHeader";i:11;s:6:"images";i:12;s:6:"videos";i:13;s:10:"sourceName";i:14;s:8:"category";i:15;s:5:"theme";i:16;s:14:"headerKeywords";i:17;s:5:"links";i:18;s:7:"authors";i:19;s:12:"goodNewsText";i:20;s:17:"isCommentsAllowed";i:21;s:4:"tags";i:22;s:15:"commercialLabel";i:23;s:9:"customUrl";i:24;s:9:"copyright";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:27:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:13:"flags[tomain]";i:3;s:12:"flags[totop]";i:4;s:17:"flags[tofixedtop]";i:5;s:12:"flags[torss]";i:6;s:15:"flags[goodNews]";i:7;s:4:"type";i:8;s:9:"mainPhoto";i:9;s:21:"socialBackgroundImage";i:10;s:27:"socialBackgroundImageHeader";i:11;s:9:"mainVideo";i:12;s:6:"images";i:13;s:10:"sourceName";i:14;s:6:"author";i:15;s:8:"category";i:16;s:5:"theme";i:17;s:14:"headerKeywords";i:18;s:5:"polls";i:19;s:5:"links";i:20;s:7:"authors";i:21;s:12:"goodNewsText";i:22;s:17:"isCommentsAllowed";i:23;s:4:"tags";i:24;s:15:"commercialLabel";i:25;s:9:"customUrl";i:26;s:9:"copyright";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:2:{i:0;s:4:"lead";i:1;s:9:"subheader";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:2:{i:0;s:8:"category";i:1;s:5:"theme";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:3:{i:0;s:21:"socialBackgroundImage";i:1;s:27:"socialBackgroundImageHeader";i:2;s:11:"socialAlign";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:3:{i:0;s:4:"lead";i:1;s:4:"type";i:2;s:9:"mainPhoto";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:3:{i:0;s:9:"subheader";i:1;s:10:"htmlButton";i:2;s:17:"isCommentsAllowed";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:5:{i:0;s:21:"socialBackgroundImage";i:1;s:27:"socialBackgroundImageHeader";i:2;s:8:"category";i:3;s:5:"theme";i:4;s:11:"socialAlign";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:5:{i:0;s:5:"feeds";i:1;s:13:"flags[tomain]";i:2;s:16:"flags[comPlace3]";i:3;s:16:"flags[comPlace7]";i:4;s:17:"isCommentsAllowed";}}',
'a:2:{s:9:"paramsKey";s:4:"name";s:4:"data";a:6:{i:0;s:4:"lead";i:1;s:9:"subheader";i:2;s:10:"sourceName";i:3;s:14:"headerKeywords";i:4;s:7:"authors";i:5;s:9:"copyright";}}',
'a:2:{s:9:"paramsKey";s:4:"type";s:4:"data";a:2:{i:0;s:16:"video_of_the_day";i:1;s:16:"photo_of_the_day";}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:1077;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:114160;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:123;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:124;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:138;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:142982;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:142;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:14;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:154;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:155;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:166;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:16;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:170;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:181490;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:182028;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:18;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:21;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:22;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:23;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:24;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:26;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:27;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:29;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:2;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:30;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:31;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:32;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:33;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:34;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:35;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:36;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:38;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:39;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:42;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:43403;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:43;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:44;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:45;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:46;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:47;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:48;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:51;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:52;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:53;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:54;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:55;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:56;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:57;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:58;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:59;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:60;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:61;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:62;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:63;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:64;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:66;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:67;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:68;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:69;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:70;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:71;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:72;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:73;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:74;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:75;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:76;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:86;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:89;}}',
'a:2:{s:9:"paramsKey";s:6:"region";s:4:"data";a:1:{i:0;i:93;}}',
'a:2:{s:9:"paramsKey";s:7:"project";s:4:"data";a:1:{i:0;i:10;}}',
'a:2:{s:9:"paramsKey";s:7:"project";s:4:"data";a:1:{i:0;i:11;}}',
'a:2:{s:9:"paramsKey";s:7:"project";s:4:"data";a:1:{i:0;i:1;}}',
'a:2:{s:9:"paramsKey";s:7:"project";s:4:"data";a:1:{i:0;i:28;}}',
'a:2:{s:9:"paramsKey";s:7:"project";s:4:"data";a:1:{i:0;i:2;}}',
'a:2:{s:9:"paramsKey";s:7:"project";s:4:"data";a:1:{i:0;i:3;}}',
'a:2:{s:9:"paramsKey";s:7:"project";s:4:"data";a:1:{i:0;i:48;}}',
'a:2:{s:9:"paramsKey";s:7:"project";s:4:"data";a:1:{i:0;i:4;}}',
'a:2:{s:9:"paramsKey";s:7:"project";s:4:"data";a:1:{i:0;i:6;}}',
'a:2:{s:9:"paramsKey";s:7:"project";s:4:"data";a:1:{i:0;i:7;}}',
'a:2:{s:9:"paramsKey";s:7:"project";s:4:"data";a:1:{i:0;i:8;}}',
'a:2:{s:9:"paramsKey";s:7:"project";s:4:"data";a:1:{i:0;i:9;}}',
'a:2:{s:9:"paramsKey";s:8:"template";s:4:"data";a:5:{i:0;i:1;i:1;i:2;i:2;i:3;i:3;i:4;i:4;i:5;}}',
'a:2:{s:9:"paramsKey";s:8:"template";s:4:"data";a:6:{i:0;i:1;i:1;i:2;i:2;i:3;i:3;i:4;i:4;i:5;i:5;i:6;}}'
 */
