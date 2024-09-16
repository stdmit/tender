package controllers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"example.com/tender/internal/database"
	"example.com/tender/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

var Validate = validator.New()

func PingServer() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		queryparam := c.Request.URL.Query()

		if len(queryparam) != 0 {
			c.IndentedJSON(500, "query is not empty")
			c.Abort()
			return
		}else{
			log.Println("query is empty")
		}

		c.Header("Access-Control-Allow-Origin","*")
		c.IndentedJSON(200, "OK")
	}
}

func CreateTender() gin.HandlerFunc {
	return func(c *gin.Context) {
		status :="Created"
		ver    :=uint(1)

		queryparam := c.Request.URL.Query()

		if len(queryparam) != 0 {
			c.IndentedJSON(500, "query is not empty")
			c.Abort()
			return
		}else{
			log.Println("query is empty")
		}

		var tender models.Tender
		if err := c.BindJSON(&tender); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		validationErr := Validate.Struct(tender)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		_,ok := models.STypesAllow[*tender.ServiceType]
		if !ok {
			c.IndentedJSON(400, "wrong service type")
			c.Abort()
			return
		}

		orgRespNames,err := database.GetOrgResponsible(*tender.OrgId)
		if err != nil{
			c.IndentedJSON(500, "can't get organization responsible from db")
			c.Abort()
			return
		}

		_,ok = orgRespNames[*tender.CreatorName]
		if !ok {
			c.IndentedJSON(400, "wrong username. you are not in organization responsible")
			c.Abort()
			return
		}
		
		tender.Status 	   = &status
		tender.Ver 		   = &ver
		tender.CreatedAt,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		query := "INSERT INTO tender(ver,name,description,servicetype,status,organization_id,created_at,creator_name) VALUES ($1,$2,$3,$4,$5,$6,$7,$8);"
		_, err = database.DB.ExecContext(context.Background(),query, *tender.Ver,*tender.Name,*tender.Description,*tender.ServiceType,*tender.Status,*tender.OrgId,
		tender.CreatedAt,*tender.CreatorName)

		if err != nil {
			c.IndentedJSON(500, "imposible insert tender to db table")
			c.Abort()
			return
		}
		
		c.JSON(http.StatusCreated, "Successfully created")

	}
}

func ListTenders() gin.HandlerFunc {
	return func(c *gin.Context) {
		queryparam := c.Request.URL.Query()
		var qoffset,qlimit,qservType []string
		var servType = []string{}
		//var offset,limit int32
		qoffset = []string{"0"}
		qlimit  = []string{"5"}
		for key, values := range queryparam {
			//log.Printf("key = %v, value(s) = %v\n", key, values)
			switch key {
			case "offset":
				qoffset = values
			case "limit":
				qlimit = values
			case "service_type":
				qservType = values
			default:
				c.IndentedJSON(400, "wrong query key")
				c.Abort()
				return
			}
		}
		if len(qservType) > 3{
			c.IndentedJSON(400, "wrong number of service_type filters (must be <=3)")
			c.Abort()
			return
		}

		for _,t := range qservType{
			_,ok := models.STypesAllow[t]
			if !ok{
				c.IndentedJSON(400, "wrong service_type value: " + t)
				c.Abort()
				return
			}
			servType = append(servType,t )
		}

		offset,err := strconv.Atoi(qoffset[0])
		if err != nil{
			c.IndentedJSON(400, "wrong offset format")
			c.Abort()
			return
		}
		
		limit,err := strconv.Atoi(qlimit[0])
		if err != nil{
			c.IndentedJSON(400, "wrong limit format")
			c.Abort()
			return
		}
		
		if offset<0 {
			c.IndentedJSON(400, "wrong offset value (must be > 0)")
			c.Abort()
			return
		}

		if  limit<1 {
			c.IndentedJSON(400, "wrong limit value (must be > 1)")
			c.Abort()
			return
		}
		
		dbQuery  := "select * from tender"
		var rows *sql.Rows
		if len(servType)>0 {
			dbQuery = "select id,name,description,status,servicetype,ver,created_at from tender where servicetype = ANY($1) and status='Published' order by name"
			rows, err = database.DB.Query(dbQuery, pq.Array(servType))
			if err != nil {
				panic(err)
			}
		}

		if len(servType)==0 {
			dbQuery = "select id,name,description,status,servicetype,ver,created_at from tender where status='Published' order by name "
			rows, err = database.DB.Query(dbQuery)
			if err != nil {
				panic(err)
			}
		}
		defer rows.Close()
		
		tenders :=[]models.Tender{}
		for rows.Next(){
			t := models.Tender{}
			//err := rows.Scan(&t.ID,&t.Ver, &t.Name,&t.Description,&t.ServiceType,&t.Status,&t.OrgId,&t.CreatedAt,&t.CreatorName)
			err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Status, &t.ServiceType,&t.Ver, &t.CreatedAt)
			if err != nil{
				log.Println(err)
				continue
			}
			tenders = append(tenders, t)
		}
		

		if offset > len(tenders) {
			c.IndentedJSON(400, "wrong offset value (should be less or equal " + strconv.Itoa(len(tenders)) + ")" )
			c.Abort()
			return			
		}
		tenders = tenders[offset:min(offset+limit,len(tenders))]
		
		c.Header("Access-Control-Allow-Origin","*")
		c.IndentedJSON(200, tenders)
	}
}

func ListMyTenders() gin.HandlerFunc {
	return func(c *gin.Context) {
		queryparam := c.Request.URL.Query()
		var qoffset,qlimit,qusername []string
		
		
		qoffset = []string{"0"}
		qlimit  = []string{"5"}

		_,ok := queryparam["username"]
		if !ok {
			c.IndentedJSON(400, "no username in query")
			c.Abort()
			return
		}

		for key, values := range queryparam {
			//log.Printf("key = %v, value(s) = %v\n", key, values)
			switch key {
			case "offset":
				qoffset = values
			case "limit":
				qlimit = values
			case "username":
				qusername = values
			default:
				log.Println(key)
				c.IndentedJSON(400, "wrong query key")
				c.Abort()
				return
			}
		}
		username:=qusername[0]

		offset,err := strconv.Atoi(qoffset[0])
		if err != nil{
			c.IndentedJSON(400, "wrong offset format")
			c.Abort()
			return
		}
		
		limit,err := strconv.Atoi(qlimit[0])
		if err != nil{
			c.IndentedJSON(400, "wrong limit format")
			c.Abort()
			return
		}
		
		if offset<0 {
			c.IndentedJSON(400, "wrong offset value (must be > 0)")
			c.Abort()
			return
		}

		if  limit<1 {
			c.IndentedJSON(400, "wrong limit value (must be > 1)")
			c.Abort()
			return
		}
		
		dbQuery := "select id,name,description,status,servicetype,ver,created_at from tender where creator_name=$1 order by name"
			rows, err := database.DB.Query(dbQuery, username)
			if err != nil {
				panic(err)
		}

		defer rows.Close()
		
		tenders :=[]models.Tender{}
		for rows.Next(){
			t := models.Tender{}
			err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Status, &t.ServiceType,&t.Ver, &t.CreatedAt)
			if err != nil{
				log.Println(err)
				continue
			}
			tenders = append(tenders, t)
		}
		

		if offset > len(tenders) {
			c.IndentedJSON(400, "wrong offset value (should be less or equal " + strconv.Itoa(len(tenders)) + ")" )
			c.Abort()
			return			
		}
		tenders = tenders[offset:min(offset+limit,len(tenders))]
		
		c.Header("Access-Control-Allow-Origin","*")
		c.IndentedJSON(200, tenders)
	}
}

func ShowStatusTender() gin.HandlerFunc{
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin","*")

		tenderID := c.Param("tenderId")		

		t,err :=database.GetTender(tenderID)
		if err != nil{
			c.IndentedJSON(400,"Can't get tender with required id")
			c.Abort()
			return
		}

		queryparam := c.Request.URL.Query()

		quserName,ok := queryparam["username"]
		if !ok{
			c.IndentedJSON(400,"Access denied username do not passed !")
			c.Abort()
			return
		}
		
		orgResp,err := database.GetOrgResponsible(*t.OrgId)
		if err != nil {
			c.IndentedJSON(500,"Can't get organization responsible")
			c.Abort()
			return
		}

		if *t.CreatorName == quserName[0] {
			c.IndentedJSON(200,map[string]string{"status":*t.Status})
			return
		}
		 _,ok = orgResp[quserName[0]]
		if ok {
			c.IndentedJSON(200,map[string]string{"status":*t.Status})
			return
		}
				
		c.IndentedJSON(400, "this username don't have access")
	}
	
}

func ChangeStatusTender() gin.HandlerFunc{
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin","*")

		tenderID := c.Param("tenderId")		

		t,err :=database.GetTender(tenderID)
		if err != nil{
			c.IndentedJSON(400,"Can't get tender with required id")
			c.Abort()
			return
		}

		queryparam := c.Request.URL.Query()

		quserName,ok := queryparam["username"]
		if !ok{
			c.IndentedJSON(400,"Access denied username do not passed !")
			c.Abort()
			return
		}

		qstatus,ok := queryparam["status"]
		if !ok{
			c.IndentedJSON(400,"Access denied status do not passed !")
			c.Abort()
			return
		}
		
		orgResp,err := database.GetOrgResponsible(*t.OrgId)
		if err != nil {
			c.IndentedJSON(500,"Can't get organization responsible")
			c.Abort()
			return
		}

		 _,ok = orgResp[quserName[0]]
		if !ok {
			c.IndentedJSON(400, "this username don't have access")
			c.Abort()
			return	
		}
		// lastver,err := database.GetTenderLastVerNum(*t.ID)
		// if err != nil {
		// 	c.IndentedJSON(400, "cant get last version num")
		// 	c.Abort()
		// 	return
		// }

		query := "update tender set status=$1 where id=$2 and ver=$3"
		_, err = database.DB.ExecContext(context.Background(),query, qstatus[0], *t.ID,*t.Ver)
		if err != nil {
			c.IndentedJSON(500, "imposible update status in db table")
			c.Abort()
			return
		}

		t,err =database.GetTender(tenderID)
		if err != nil {
			c.IndentedJSON(400,"Can't get tender with required id")
			c.Abort()
			return
		}

		c.IndentedJSON(200,map[string]string{"status":*t.Status})
		
		//return
	}
	
}

func EditTender() gin.HandlerFunc {
	return func(c *gin.Context) {
	
		tendId := c.Param("tenderId")
		queryparam := c.Request.URL.Query()
		
	
		tendCreator,ok := queryparam["username"]
		if !ok{
			c.IndentedJSON(400, "username not passed")
			c.Abort()
			return
		}

		t,err :=database.GetTender(tendId)
		if err != nil{
			c.IndentedJSON(400,"Can't get tender with required id")
			c.Abort()
			return
		}

		if *t.CreatorName != tendCreator[0]{
			c.IndentedJSON(400, "you are not creator")
			c.Abort()
			return
		}

		if err := c.BindJSON(&t); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		
		// tender.Status 	   = &status
		*t.Ver+=1
		t.CreatedAt,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		query := "INSERT INTO tender(id,ver,name,description,servicetype,status,organization_id,created_at,creator_name) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);"
		_, err = database.DB.ExecContext(context.Background(),query, *t.ID,*t.Ver,*t.Name,*t.Description,*t.ServiceType,*t.Status,*t.OrgId,
		t.CreatedAt,*t.CreatorName)
		if err != nil {
			c.IndentedJSON(500, "imposible insert tender to db table")
			c.Abort()
			return
		}
		
		c.JSON(http.StatusCreated, t)

	}
}

func RollbackVerTender() gin.HandlerFunc {
	return func(c *gin.Context) {
	
		tendId := c.Param("tenderId")

		tendVer,err := strconv.Atoi(c.Param("ver"))
		if err != nil{
			c.IndentedJSON(400, "wrong format of tender version")
			c.Abort()
			return
		}
		
		queryparam := c.Request.URL.Query()
		tendCreator,ok := queryparam["username"]
		if !ok{
			c.IndentedJSON(400, "username not passed")
			c.Abort()
			return
		}

		t,err :=database.GetTenderVer(tendId,tendVer)
		if err != nil{
			c.IndentedJSON(400,"Can't get tender with required id and ver")
			c.Abort()
			return
		}

		if *t.CreatorName != tendCreator[0]{
			c.IndentedJSON(400, "you are not creator")
			c.Abort()
			return
		}
		lastVer, err := database.GetTenderLastVerNum(tendId)
		if err!=nil{
			c.IndentedJSON(500, "cant get last version number from db")
			c.Abort()
			return
		}
		// tender.Status 	   = &status
		*t.Ver=uint(lastVer)+1
		t.CreatedAt,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		query := "INSERT INTO tender(id,ver,name,description,servicetype,status,organization_id,created_at,creator_name) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);"
		_, err = database.DB.ExecContext(context.Background(),query,*t.ID, *t.Ver,*t.Name,*t.Description,*t.ServiceType,*t.Status,*t.OrgId,
		t.CreatedAt,*t.CreatorName)
		if err != nil {
			c.IndentedJSON(500, "imposible insert tender to db table")
			c.Abort()
			return
		}
		
		c.JSON(http.StatusCreated, t)

	}
}



func PingDatabase() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		queryparam := c.Request.URL.Query()

		if len(queryparam) != 0 {
			c.IndentedJSON(500, "query is not empty")
			c.Abort()
			return
		}else{
			database.PsqlInfo()
			database.TestQuery()
		}

		c.Header("Access-Control-Allow-Origin","*")
		c.IndentedJSON(200, "OK")
	}
}


