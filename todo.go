package main

// 1. Selections from auth_item, auth_item_child. In php implementation there is project condition. Is it actual now?
// 2. All rules (php serialized to json format map)
// 		Here add json format in new table field as a serialized:
// 		- News_Permission_Manager::assign,
// 		- News_Permission_Manager::saveAssignment,
// 		- News_Permission_Manager::createItem,
// 		- News_Permission_Manager::saveItem.
// 3. May be there is no reason to do bulk check instead of sync checking.
// 		Execute bulk checking for user with several actions for checking.
// 4. Don`t have to make an async checking. Only for several users (unlikely case).
// 		Check implementing time for both.
// 5. Check, what kind of data, may be serialized string with parameters.
// 6. May be don`t hav to use differrent logic for all types (boolean, integer, string, etc)
