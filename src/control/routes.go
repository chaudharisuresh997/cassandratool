package control

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var Session *gocql.Session

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"EmpSave",
		"POST",
		"/empCreate",
		EmpCreate,
	},
	Route{
		"CreateKeySpacePath",
		"POST",
		"/cassandra/createKeyspace",
		CreateKeySpacePath,
	},
	Route{
		"EmpGET",
		"GET",
		"/Get",
		Get,
	},
}

func CreateKeySpacePath(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside EmpCreate")
	var emp Emp
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1068487))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &emp); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	name := &emp.Name
	n1 := *name
	Createkeyspace(n1)

}
func EmpCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside EmpCreate")
	var emp Emp
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1068487))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &emp); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	SaveEmpToDb(&emp)

}

//get all the records
func Get(w http.ResponseWriter, r *http.Request) {
	var empList []Emp

	m := map[string]interface{}{}
	type AllEmpsResponse struct {
		Emps []Emp `json:"messages"`
	}
	//create query
	//query1 := "SELECT name FROM shivapreals.emptable"
	query1 := "SELECT * FROM system_schema.keyspaces"
	fmt.Println(query1)
	var ses = GetSession()
	defer ses.Close()
	iterable := ses.Query(query1).Iter()
	for iterable.MapScan(m) {

		fmt.Print("**")
		empList = append(empList, Emp{
			//Id:   m["id"].(gocql.UUID),
			Name: m["keyspace_name"].(string),
		})
		m = map[string]interface{}{}
	}
	fmt.Println(empList)
	json.NewEncoder(w).Encode(AllEmpsResponse{Emps: empList})
}
func GetSession() *gocql.Session {
	fmt.Println("Sesis")
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "system"
	Session, err := cluster.CreateSession()
	if err != nil {
		fmt.Println("Session NULL")
	}

	return Session
}
func SaveEmpToDb(emp *Emp) {
	var gocqlUuid gocql.UUID
	gocqlUuid = gocql.TimeUUID()

	var ses = GetSession()
	defer ses.Close()
	fmt.Println("cassandra init done")
	// writing data to Cassandra
	if err := ses.Query(`
      INSERT INTO emptable (id, Name) VALUES (?, ?)`,
		gocqlUuid, &emp.Name).Exec(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("yes")
	}
}

//describe keyspace
func ListKeyspace() {
	var empList []Emp

	m := map[string]interface{}{}
	type AllEmpsResponse struct {
		Emps []Emp `json:"messages"`
	}
	//create query
	//query1 := "SELECT name FROM shivapreals.emptable"
	query1 := "SELECT * FROM system_schema.keyspaces"
	//createkeyspace := "CREATE KEYSPACE i2Tutorials	WITH replication = {'class': 'NetworkTopologyStrategy', 'DC1' : 1, 'DC2' : 3} AND durable_writes = false;"
	//dropkeyspace := "drop keyspace i2tutorials"
	fmt.Println(query1)
	var ses = GetSession()
	defer ses.Close()
	iterable := ses.Query(query1).Iter()
	for iterable.MapScan(m) {

		fmt.Print("**")
		empList = append(empList, Emp{
			//Id:   m["id"].(gocql.UUID),
			Name: m["keyspace_name"].(string),
		})
		m = map[string]interface{}{}
	}
	fmt.Println(empList[3:5])

}

//drop keyspace
func Dropkeyspace(keyspace_name string) {
	//
	fmt.Println("Dropping keyspace.....")

	var empList []Emp

	m := map[string]interface{}{}
	type AllEmpsResponse struct {
		Emps []Emp `json:"messages"`
	}

	dropkeyspace := "drop keyspace " + keyspace_name
	fmt.Println(dropkeyspace)
	var ses = GetSession()
	defer ses.Close()
	iterable := ses.Query(dropkeyspace).Iter()
	for iterable.MapScan(m) {

		fmt.Print("**")
		empList = append(empList, Emp{

			Name: m["keyspace_name"].(string),
		})
		m = map[string]interface{}{}
	}
	fmt.Println(empList)

}

//create keyspace
func Createkeyspace(keyspace_name string) {
	//
	fmt.Println("create keyspace.....")

	var empList []Emp

	m := map[string]interface{}{}
	type AllEmpsResponse struct {
		Emps []Emp `json:"messages"`
	}
	createkeyspace := "CREATE KEYSPACE " + keyspace_name + " WITH replication = {'class': 'NetworkTopologyStrategy', 'DC1' : 1, 'DC2' : 3} AND durable_writes = false;"

	fmt.Println(createkeyspace)
	var ses = GetSession()
	defer ses.Close()
	iterable := ses.Query(createkeyspace).Iter()
	for iterable.MapScan(m) {

		fmt.Print("**")
		empList = append(empList, Emp{

			Name: m["keyspace_name"].(string),
		})
		m = map[string]interface{}{}
	}
	fmt.Println(empList)

}

//create keyspace
func Createtable(table string) {
	var buffer bytes.Buffer
	buffer.WriteString("CREATE TABLE shivapreals.users (")
	buffer.WriteString("userid text PRIMARY KEY,")
	buffer.WriteString("first_name text,last_name text,emails set<text>,top_scores list<int>,todo map<timestamp, text>);")
	//
	fmt.Println("create table.....")

	var empList []Emp

	m := map[string]interface{}{}
	type AllEmpsResponse struct {
		Emps []Emp `json:"messages"`
	}
	createTable := buffer.String()
	fmt.Println(createTable)
	var ses = GetSession()
	defer ses.Close()
	iterable := ses.Query(createTable).Iter()
	for iterable.MapScan(m) {

		empList = append(empList, Emp{

			Name: m["keyspace_name"].(string),
		})
		m = map[string]interface{}{}
	}
	fmt.Println(empList)

}

//create keyspace
func DescTables(table string) {

	fmt.Println("DescTables table.....")

	var empList []Emp

	m := map[string]interface{}{}
	type AllEmpsResponse struct {
		Emps []Emp `json:"messages"`
	}
	DescTables := "Desc tables"
	fmt.Println(DescTables)
	var ses = GetSession()
	defer ses.Close()
	iterable := ses.Query(DescTables).Iter()
	for iterable.MapScan(m) {
		fmt.Println(m)
		for k, v := range m {
			fmt.Println("k:", k, "v:", v)
		}
		/*empList = append(empList, Emp{

			Name: m["keyspace_name"].(string),
		})*/
		m = map[string]interface{}{}
	}

	fmt.Println(empList)

}
