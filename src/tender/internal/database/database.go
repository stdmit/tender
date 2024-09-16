package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"example.com/tender/internal/models"
	_ "github.com/lib/pq"
)
var DB *sql.DB


func PsqlConnect() {
	host   	 := os.Getenv("POSTGRES_HOST")
	port 	 := os.Getenv("POSTGRES_PORT")
	user 	 := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname   := os.Getenv("POSTGRES_DATABASE")
	
	pInfo := fmt.Sprintf("host=%s port=%s user=%s " + "password=%s dbname=%s", host, port, user, password, dbname)
	
	dbcur, err := sql.Open("postgres", pInfo)
	if err != nil {
		panic(err)
	}
	DB = dbcur
	//defer db.Close()

	err = DB.Ping() 
	if err != nil {
		panic(err)
	}
	log.Println("successfully connected!")

}


func PsqlInfo() {
	host   	 := os.Getenv("POSTGRES_HOST")
	port 	 := os.Getenv("POSTGRES_PORT")
	user 	 := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname   := os.Getenv("POSTGRES_DATABASE")
	
	pInfo := fmt.Sprintf("host=%s port=%s user=%s " + "password=%s dbname=%s", host, port, user, password, dbname)
	
	dbcur, err := sql.Open("postgres", pInfo)
	if err != nil {
		panic(err)
	}
	DB = dbcur
	//defer db.Close()

	err = DB.Ping() 
	if err != nil {
		panic(err)
	}
	log.Println("successfully connected!")

}

func GetTender(tenderID string) (t models.Tender,err error) {
	sqlQuery := "select t.id, t.ver, t.name, t.description, t.servicetype,"  +
	            "t.status, t.organization_id, t.created_at, t.creator_name  from tender t " +
				"join (" +
				"select id, MAX(ver) as ver from tender where id=$1 group by id" +
				") t2 " +
				"on t.id=t2.id and t.ver=t2.ver"

    err = DB.QueryRow(sqlQuery,tenderID).Scan(&t.ID, &t.Ver, &t.Name, &t.Description, &t.ServiceType,
											   &t.Status, &t.OrgId, &t.CreatedAt, &t.CreatorName)
	if err != nil{
		if err !=sql.ErrNoRows{
			log.Println(err)
		}
	}
	return t, err
}

func GetTenderVer(tenderID string,tenderVer int) (t models.Tender,err error) {
	sqlQuery := "select t.id, t.ver, t.name, t.description, t.servicetype,"  +
	            "t.status, t.organization_id, t.created_at, t.creator_name  from tender t " +
				"where t.id=$1 and t.ver=$2"
				

    err = DB.QueryRow(sqlQuery,tenderID,tenderVer).Scan(&t.ID, &t.Ver, &t.Name, &t.Description, &t.ServiceType,
											   &t.Status, &t.OrgId, &t.CreatedAt, &t.CreatorName)
	if err != nil{
		if err !=sql.ErrNoRows{
			log.Println(err)
		}
	}
	return t, err
}

func GetTenderLastVerNum(tenderID string) (ver int,err error) {
	sqlQuery := "select MAX(ver) as ver from tender where id=$1 group by id"

    err = DB.QueryRow(sqlQuery,tenderID).Scan(&ver)
	if err != nil{
		if err !=sql.ErrNoRows{
			log.Println(err)
		}
	}
	return ver, err
}

func GetOrgResponsible(orgID string) (orgResp map[string]bool,err error) {
		orgResp = make(map[string]bool)
		sqlQuery := "select e.username from organization_responsible o join "+
				    "employee e on o.user_id=e.id where o.organization_id=$1"
		
		
		rows, err := DB.Query(sqlQuery,orgID)
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		
		for rows.Next(){
			orgRespName := ""
			err := rows.Scan(&orgRespName)
			if err != nil{
				log.Println(err)
				continue
			}
			orgResp[orgRespName]=true
		}
	return orgResp, err
}

func GetTenderResponsible(tenderID string) (orgResp map[string]bool,err error) {
	orgResp = make(map[string]bool)
	sqlQuery := "select e.username from organization_responsible o join "+
				"employee e on o.user_id=e.id where " +
				" o.organization_id in " + 
				"(" +
				"select organization_id from tender " +
				"where id=$1 and ver in (select max(ver) from tender where id=$1)" +
				")"
	
	
	rows, err := DB.Query(sqlQuery,tenderID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	
	for rows.Next(){
		orgRespName := ""
		err := rows.Scan(&orgRespName)
		if err != nil{
			log.Println(err)
			continue
		}
		orgResp[orgRespName]=true
	}
	return orgResp, err
}

func UserIdInEmployee(uid string) bool{
	sqlQuery := "select id from employee where id=$1"
	err := DB.QueryRow(sqlQuery,uid).Scan(&uid)
	if err != nil{
		if err!=sql.ErrNoRows{
			log.Println(err)
		}
		return false
	}
	return true

}

func UserIdInOrgResp(uid string) bool{
	sqlQuery := "select user_id from organization_responsible where user_id=$1"
	err := DB.QueryRow(sqlQuery,uid).Scan(&uid)
	if err != nil{
		if err!=sql.ErrNoRows{
			log.Println(err)
		}
		return false
	}
	return true

}


func GetBid(bidID string) (b models.Bid,err error) {
	sqlQuery := "select b.id, b.name, b.description, b.status, b.tender_id, b.authortype, b.authorid, b.ver, b.created_at from bid b join " + 
			   "(select id, max(ver) as ver from bid group by id) as b2 on b.id = b2.id and b.ver=b2.ver " +
			   "where b.id=$1"

    err = DB.QueryRow(sqlQuery,bidID).Scan(&b.ID, &b.Name, &b.Description, &b.Status, &b.TenderId, &b.AuthType, &b.AuthId, &b.Ver, &b.CreatedAt)
	if err != nil{
		if err !=sql.ErrNoRows{
			log.Println(err)
		}
	}

	return b, err
}

func GetBidVer(bidID string,bidVer int) (b models.Bid,err error) {
	sqlQuery := "select b.id, b.name, b.description, b.status, b.tender_id, b.authortype, b.authorid, b.ver, b.created_at from bid b " + 
			   "where b.id=$1 and b.ver=$2"

    err = DB.QueryRow(sqlQuery,bidID,bidVer).Scan(&b.ID, &b.Name, &b.Description, &b.Status, &b.TenderId, &b.AuthType, &b.AuthId, &b.Ver, &b.CreatedAt)
	if err != nil{
		if err !=sql.ErrNoRows{
			log.Println(err)
		}
	}

	return b, err
}

func GetBidLastVerNum(bidID string) (ver int,err error) {
	sqlQuery := "select MAX(ver) as ver from bid where id=$1 group by id"

    err = DB.QueryRow(sqlQuery,bidID).Scan(&ver)
	if err != nil{
		if err !=sql.ErrNoRows{
			log.Println(err)
		}
	}
	return ver, err
}

func UserIdByName(username string) (uid string, err error) {
	sqlQuery := "select id from employee where username=$1"

	err = DB.QueryRow(sqlQuery,username).Scan(&uid)
	if err != nil{
		if err !=sql.ErrNoRows{
			log.Println(err)
		}
	}
	return uid,err
}


func TestQuery(){
	
	rows, err := DB.Query("select username from employee")
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    
	uname:=""
	unames :=[]string{}
	
	for rows.Next(){
        //p := product{}
        err := rows.Scan(&uname)
        if err != nil{
            fmt.Println(err)
            continue
        }
        unames = append(unames, uname)
    }
	log.Println(unames)
}
