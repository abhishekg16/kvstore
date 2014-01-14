package main

import "fmt"
import zmq "github.com/alecthomas/gozmq"
import "encoding/json"
import "os"
import "log"

type request struct {
        Command string
        Uid int
        Parameters []string
}

type reply struct {
        Result string
        Uid int
        Err error
}


func SendRequest(socket *zmq.Socket,req *request)  ([]byte , error){ 	
	 msg, err := json.Marshal(*req)
 	 fmt.Print("Encoded Request:" )
	 os.Stdout.Write(msg)
	 fmt.Println()
	 if err != nil {
		fmt.Println("error:", err)
		return nil, err
	 }
	 (*socket).Send(msg, 0)
	 println("Sending", msg)
	 return (*socket).Recv(0)
}

func Reset(req *request) {
	req.Command=""
	req.Uid=0
	req.Parameters=nil
}

func logReply(rep reply) {
	log.Println("logging Reply")
	log.Println(rep.Result)
	log.Println(rep.Uid)
	if rep.Err != nil {
		log.Println(rep.Err)
	}
}

var UserId int

func main() {
	 

	context, _ := zmq.NewContext()
	 socket, _ := context.NewSocket(zmq.REQ)
	 socket.Connect("tcp://127.0.0.1:8080")
	 var req request
	 var rep reply
	 UserId = 0

	 //Request to initialize service
	
	 Reset(&req)
	 req.Command =	"Initialize"
         log.Println("Requesting to initialize a service")
	 msg, err := SendRequest(socket,&req)
	 checkError(err)
	 err = json.Unmarshal(msg,&rep)
         if err != nil {
                    fmt.Println("error:", err)
         }
         logReply(rep)
         UserId = rep.Uid
	 fmt.Printf("Got User Id: %d",UserId)
	 fmt.Println();	

	 // Create a table
	 Reset(&req)
	 req.Command =	"CreateTable"
	 req.Uid = UserId
	 req.Parameters = []string{"table1"}
         log.Println("Creating a new Table")
	 msg, err= SendRequest(socket,&req)
	 err = json.Unmarshal(msg,&rep)
         if err != nil {
                    fmt.Println("error:", err)
         }
         logReply(rep)
	 checkError(err)
	
	 // Insert a key
	 Reset(&req)
	 req.Command =	"Put"
	 req.Uid = UserId
	 req.Parameters = []string{"table1","key1","value1"}
         log.Println("Inserting value in table 1")
	 msg, err= SendRequest(socket,&req)
	 err = json.Unmarshal(msg,&rep)
         if err != nil {
                    fmt.Println("error:", err)
         }
         logReply(rep)
	 checkError(err)
	
	// get a Key-Value Pair
	
	 Reset(&req)
	 req.Command =	"Get"
	 req.Uid = UserId
	 req.Parameters = []string{"table1","key1"}
         log.Println("Fetching value from table")
	 msg, err= SendRequest(socket,&req)
	 err = json.Unmarshal(msg,&rep)
         if err != nil {
                    fmt.Println("error:", err)
         }
         logReply(rep)
	 checkError(err)
	
}

func checkError(err error ) {
	if  err != nil {
		log.Fatal("Error : %s",err)
		os.Exit(1)
	}
}
