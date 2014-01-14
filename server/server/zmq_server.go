package main

import (
	"fmt"
	"log"
	zmq "github.com/alecthomas/gozmq"
	ds "github.com/abhishekg16/kvstore/server/dataStore"
	"encoding/json"
	"os"
)


// Request Type
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



func ProcessResponse(result string, uid int, err error)  ([]byte){
	 res := reply{result,uid,err}
         msg, err := json.Marshal(res)
         log.Println("Encoded Response:" )
         os.Stdout.Write(msg)
         fmt.Println() 
         if err != nil {
                fmt.Println("error:", err)
                return nil
         }
         return msg
}


func handleClient(msg []byte) []byte{
		log.Println("Incoming Request")
	        os.Stdout.Write(msg)
		var req request
		err := json.Unmarshal(msg,&req)
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Println(req)
		log.Printf("%d : Requested: %s ", req.Uid, req.Command)
		log.Println()
		switch {
		case req.Command == "Initialize":
			uid := ds.CreateUser()
			log.Printf("New User Created with userId : %d",uid)
			fmt.Println()
			return ProcessResponse("New User created",uid,nil)
		
		case req.Command == "CreateTable":
			uid := req.Uid			
			if uid == 0  || len(req.Parameters) != 1 {
				fmt.Println("Invalid Arguments " )
				return ProcessResponse("Invalid Argument",uid,nil)
			}
			tableName := req.Parameters[0]
			ok, _ := ds.CreateTable(uid, tableName)   
			if (ok)	{		
				log.Println("Table Created")
				return ProcessResponse("Table Created",uid,nil)
			} else {
				fmt.Println("Table Not Created")
				return ProcessResponse("Table Created",uid,nil)
			}
			return ProcessResponse("Table Created",uid,nil)
	
		case req.Command == "Put":	
			uid := req.Uid			
			if uid == 0  || len(req.Parameters) != 3 {
				fmt.Println("Invalid Arguments " )
				return ProcessResponse("Invalid Argument",uid,nil)
			}
			tableName := req.Parameters[0]
			key:= req.Parameters[1]
			value:= req.Parameters[2]
	
			ok, err:= ds.Put(uid, tableName,key,value)   
			if (ok)	{		
				log.Println("Inserted Key-Value Pair")
				return ProcessResponse("Inserted Key-Value Pair",uid,err)
			} else {
				fmt.Println("Error occured")
				return ProcessResponse("Inserted Key-Value Pair",uid,err)
			}
			return ProcessResponse("Value insered",uid,nil)
		
		case req.Command == "Get": 
			uid := req.Uid			
			if uid == 0  || len(req.Parameters) != 2 {
				fmt.Println("Invalid Arguments " )
				return ProcessResponse("Invalid Argument",uid,nil)
			}
			tableName := req.Parameters[0]
			key:= req.Parameters[1]
	
			value, err:= ds.Get(uid, tableName,key)   
			if err != nil {		
				log.Println("Can not fetch the value")
				return ProcessResponse("Inserted Key-Value Pair",uid,err)
			}
			return ProcessResponse(value,uid,nil)
		}
	 	return ProcessResponse("Wrong command",0,nil)
}


func main() {
	context, _ := zmq.NewContext()
	socket, err := context.NewSocket(zmq.REP)
	checkError(err)
	socket.Bind("tcp://127.0.0.1:8080")
	log.Println("Server  started ")
	log.Println("Initailizing the data service") 
	ds.InitializeDataService()
	for {
		msg, err  := socket.Recv(0)
		if err != nil {	
			continue
		}
		rep := handleClient(msg)
		socket.Send(([]byte(rep)),0)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ",err.Error())
		os.Exit(1)
	}
}	


