package DataStore

import (
	"fmt"
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


// _userId Map maintain the a map from userId to there namespace. So in case there is not map entry available corresponding to a particular map id. Then it implies the user does not exist


// table type contians the key/value - in present implementation  both key values are string
type table map[string]string

// dir contains a set of tables which are alloted to a single user
type dir map[string](*table) 


var userCounter int; 

// userDir is a maps a particular user to dir to his specific directory 
var _userDir map[int](*dir) 

func InitializeDataService () bool {
	userCounter = 0
	_userDir = make( map[int](*dir) , 0)
	return true
}

func CreateUser() int {
	t_dir := make(dir,0)
	userCounter++
	t_userID := userCounter
	_userDir[t_userID] = &t_dir
	return t_userID
}

func getUserDir(uid int) (*dir, error ) {
	t_dir, ok := _userDir[uid]
	if ( ok ) { 
		return t_dir, nil
	}
	return nil, fmt.Errorf("user directory does not exit")
}

func getUserTable(uid int, tableName string) ( *table, error ) {
	u_dir, ok := _userDir[uid]	
	if ( !ok ) {
		return nil, fmt.Errorf( "%d : user id does not exist", uid )
	}
	u_table , ok :=	(*u_dir)[tableName]
	if ( !ok ) {
		return nil, fmt.Errorf( "%s : table does not exit", tableName )
	}
	return u_table,nil
}

func CreateTable( uid int, tableName string) (bool, error) {	
	u_dir, err := getUserDir(uid) 		
	if err != nil {
		fmt.Println(err);
		return false, err
	}
	_ , ok := (*u_dir)[tableName]
	if ( ok ) {
		return false, fmt.Errorf("%s table already exist")	
	}
	t_table := make(table,0)	
	(*u_dir)[tableName] = &t_table
	return true, nil
}

func getTable(uid int, tableName string) ( *table , error ) {
	u_dir , ok := _userDir[uid]
	if ( !ok ) {
		return nil, fmt.Errorf("%g :User Id does not exit ",uid) 
	}
	table, ok := (*u_dir)[tableName]   
	
	if ( !ok ) {
		return nil, fmt.Errorf("Table %s doen not exist",tableName ) 
	}
	return table , nil 
}

func Get(uid int, tableName string , key string ) ( string, error) {
	u_table, err :=	getTable(uid,tableName)
	if err != nil {
		return "" , err
	}
	value, ok := ( *u_table )[key]
	if ( !ok ) {
		return "", fmt.Errorf("key does not exist");
	} 
	return value, nil
} 

func Put( uid int, tableName string , key string , value string) ( bool, error ) {
	u_table , _ := getTable(uid, tableName)
	( *u_table )[key] = value
	return true , nil
}

func deleteKey( uid int, tableName string, key string) (bool, error) {
	u_table , err := getTable(uid, tableName)
	if u_table == nil {
		return false, err
	}
	_, ok := (*u_table)[key]
	if ( !ok ) {	
		return false , fmt.Errorf("The Key does not exit")			
	}
	delete((*u_table),key)
	return true, nil
}

/*
func getTableList( uid int) [] string {
	tableMap , ok = userIds[uid]
	if (!ok) 
		return nil, fmt.Errorf("%d : user id doen not exist" )
	tables := tableMap.Keys
	fmt.Print
}

*/

