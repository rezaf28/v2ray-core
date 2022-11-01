package vless

import (
	"strings"
	"sync"
        "fmt"
        "bufio"
        "os"
        "time"
	"github.com/v2fly/v2ray-core/v5/common/protocol"
	"github.com/v2fly/v2ray-core/v5/common/uuid"
	"github.com/Joker666/AsyncGoDemo/async"
)



type USER struct {
        UUID string
        IP   string
        Counter int
}

var userlist = []*USER{}

func GetUUIDIP(Num int) string{
	return userlist[Num].IP
}

func GetUUIDCount(ID string) int{
	for m:=0;m<len(userlist);m++ {
		if userlist[m].UUID == ID{
			return m
		}
	}
	return 0;
}

func CheckIP(ID string,IP string) bool{
	num := GetUUIDCount(ID)
	if IP == GetUUIDIP(num) {
		userlist[num].Counter = 60
	} else {
		if userlist[num].Counter == 0 {
			userlist[num].IP = IP
		} else {
			return false
		}
	}
	return true
}

func DoneAsync() int {
	for 
	{
		for m:=0;m<len(userlist);m++{
			if userlist[m].Counter != 0{
				userlist[m].Counter = userlist[m].Counter - 1;
			}
		}

		time.Sleep(1 * time.Second)
    		if (1 != 1){
        		break
    		}	
	}
	return 1
}



func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

// Validator stores valid VLESS users.
type Validator struct {
	// Considering email's usage here, map + sync.Mutex/RWMutex may have better performance.
	email sync.Map
	users sync.Map
}

// Add a VLESS user, Email must be empty or unique.
func (v *Validator) Add(u *protocol.MemoryUser) error {
	if u.Email != "" {
		_, loaded := v.email.LoadOrStore(strings.ToLower(u.Email), u)
		if loaded {
			return newError("User ", u.Email, " already exists.")
		}
	}
	v.users.Store(u.Account.(*MemoryAccount).ID.UUID(), u)
	fusr:= new(USER)
	fusr.UUID = u.Account.(*MemoryAccount).ID.String()
	fusr.IP = "0.0.0.0"
	fusr.Counter = 0
	userlist = append(userlist,fusr)
			
	list_uuid, err := readLines("UUID.txt")
	if err != nil {
        	fmt.Print(err)
		return newError("Can't Read User List , Create UUID.txt")
        }
        
	for i := 0; i < len(list_uuid); i++ {
		if(list_uuid[i] != "" && list_uuid[i] != " "){
			uid, uerror := uuid.ParseString(list_uuid[i])
			if(uerror != nil){
				return newError("Error: Check User List!")
			}
			usr:= new(USER)
			usr.UUID = uid.String()
			usr.IP = "0.0.0.0"
			usr.Counter = 0
			userlist = append(userlist,usr)
			v.users.Store(uid,u)
		}
	}
	future := async.Exec(func() interface{} {
		return DoneAsync()
	})
	_ = future
	return nil
}

// Del a VLESS user with a non-empty Email.
func (v *Validator) Del(e string) error {
	if e == "" {
		return newError("Email must not be empty.")
	}
	le := strings.ToLower(e)
	u, _ := v.email.Load(le)
	if u == nil {
		return newError("User ", e, " not found.")
	}
	v.email.Delete(le)
	v.users.Delete(u.(*protocol.MemoryUser).Account.(*MemoryAccount).ID.UUID())
	return nil
}

// Get a VLESS user with UUID, nil if user doesn't exist.
func (v *Validator) Get(id uuid.UUID, remoteip string) *protocol.MemoryUser {
	u, _ := v.users.Load(id)
	if u != nil {
		if CheckIP(id.String(), remoteip) {
			return u.(*protocol.MemoryUser)
		}
	}
	return nil
}
