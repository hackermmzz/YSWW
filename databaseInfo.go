package main
import(
	"fmt"
	"strconv"
)

type void struct{}


type User struct {
    email   string
    passwd  string
	date 	string
	vip		int
}

type XNCYInfo struct{
	account	string
	person string
	clothes string
	generate  string
	date 	string
	status	bool
}

type WSTInfo struct{
	account	string
	description string
	generate string
	date 	string
	status	bool
}

type RLFGHInfo struct{
	account	string
	face string
	description string
	generate string
	date 	string
	status	bool
}

type PortraitInfo struct{
	account	string
	person string
	generate string
	date 	string
	status	bool
}

type KouTuInfo struct{
	account	string
	image string
	generate string
	date 	string
	status	bool
}

type RegistVerityCodeInfo struct{
	account string
	code	string
	date 	string
	status	bool
}

func (user*User) Print(){
	fmt.Println("email:"+user.email+" password:"+user.passwd+" date:"+user.date+" vip:"+strconv.Itoa(user.vip))
}