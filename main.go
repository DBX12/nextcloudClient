package main

import (
	"fmt"
	"nextcloudClient/nextcloudClient"
)

func main() {
	client := nextcloudClient.NewClient("https://nc.dbx-12.de", "Deniz", "XE66D-6Yp4A-L2LgW-rLHk6-FmPbf")
	//users, err := client.GetUsers()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(users)
	//
	//details, err := client.GetUserDetails("test")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(*details)
	//
	//success, err := client.PromoteToSubadmin("test", "testGroup")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("promote:", success)
	//
	//groups, err := client.GetSubadminGroups("test")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("group admin of:",groups)
	//
	//success, err = client.DemoteFromSubadmin("test", "testGroup")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("demote:", success)

	groups, err := client.GetGroups()
	if err != nil {
		panic(err)
	}
	fmt.Println(groups)

	groupMembers, err := client.GetGroupMembers("testGroup")
	if err != nil {
		panic(err)
	}
	fmt.Println(groupMembers)

	groupSubadmins, err := client.GetGroupSubadmins("testGroup")
	if err != nil {
		panic(err)
	}
	fmt.Println(groupSubadmins)
}
