package controllers

import (
	"context"
	"time"
	"log"
	"net/http"
	"strconv"
	

	"example.com/tender/internal/database"
	"example.com/tender/internal/models"
	"github.com/gin-gonic/gin"
)
//var Validate = validator.New()

func CreateBid() gin.HandlerFunc {
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

		var b models.Bid
		if err := c.BindJSON(&b); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		validationErr := Validate.Struct(b)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}

		_,ok := models.AuthTypeAllow[*b.AuthType]
		if !ok {
			c.IndentedJSON(400, "wrong service type")
			c.Abort()
			return
		}

		t,err := database.GetTender(*b.TenderId)
		if err != nil {
			c.IndentedJSON(500, "cant find tender with this id")
			c.Abort()
			return
		}

		if *t.Status != "Published"{
			c.IndentedJSON(400, "tender unavailable")
			c.Abort()
			return
		}
		log.Println(status,ver)
		if *b.AuthType == "Organization"{
			if !database.UserIdInOrgResp(*b.AuthId){
				c.IndentedJSON(400, "user_id not in organization")
				c.Abort()
				return
			}
		}

		if *b.AuthType == "User"{
			if !database.UserIdInEmployee(*b.AuthId){
				c.IndentedJSON(400, "user_id not found")
				c.Abort()
				return
			}
		}

		b.Status 	   = &status
		b.Ver 		   = &ver
		b.CreatedAt,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		query := "INSERT INTO bid(ver,name,description,status,tender_id,authortype,authorid,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8);"
		_, err = database.DB.ExecContext(context.Background(),query, *b.Ver,*b.Name,*b.Description,*b.Status,*b.TenderId,
		*b.AuthType,*b.AuthId, b.CreatedAt)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "imposible insert bid to db table")
			c.Abort()
			return
		}
		
		c.JSON(http.StatusCreated, "Bid successfully created")

	}
}

func ListMyBids() gin.HandlerFunc {
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

		// return last version of bid
		
		dbQuery := "select b.id, b.name, b.status, b.authortype, b.authorid, b.ver, b.created_at from bid b join " + 
				   "employee e on b.authorid = e.id join " + 
				   "(select id, max(ver) as ver from bid group by id) as b2 on b.id = b2.id and b.ver=b2.ver " +
				   "where e.username=$1"
			rows, err := database.DB.Query(dbQuery, username)
			if err != nil {
				panic(err)
		}

		defer rows.Close()
		
		bids :=[]models.Bid{}
		for rows.Next(){
			b := models.Bid{}
			err := rows.Scan(&b.ID, &b.Name, &b.Status, &b.AuthType, &b.AuthId, &b.Ver, &b.CreatedAt)
			if err != nil{
				log.Println(err)
				continue
			}
			bids = append(bids, b)
		}
		

		if offset > len(bids) {
			c.IndentedJSON(400, "wrong offset value (should be less or equal " + strconv.Itoa(len(bids)) + ")" )
			c.Abort()
			return			
		}
		bids = bids[offset:min(offset+limit,len(bids))]
		
		c.Header("Access-Control-Allow-Origin","*")
		c.IndentedJSON(200, bids)
	}
}

func ListTenderBids() gin.HandlerFunc {
	return func(c *gin.Context) {
		queryparam := c.Request.URL.Query()
		var qoffset,qlimit,qusername []string
		

		qoffset = []string{"0"}
		qlimit  = []string{"5"}
		tenderID := c.Param("Id")

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

		// t,err:=database.GetTender(tenderID)
		// if err != nil{
		// 	c.IndentedJSON(400, "tender with this id not found")
		// 	c.Abort()
		// 	return
		// }

		orgResp,err := database.GetTenderResponsible(tenderID)
		if err != nil{
			c.IndentedJSON(400, "can't get organization responsible for tenderId")
			c.Abort()
			return
		}
		
		dbQuery := "select b.id, b.name, b.status, b.authortype, b.authorid, b.ver, b.created_at from bid b join " + 
				   "employee e on b.authorid = e.id join " + 
				   "(select id, max(ver) as ver from bid group by id) as b2 on b.id = b2.id and b.ver=b2.ver " +
				   "where b.tender_id=$1 and b.status='Published'"
		
		_,ok = orgResp[username]
		if ok {
			dbQuery = "select b.id, b.name, b.status, b.authortype, b.authorid, b.ver, b.created_at from bid b join " + 
				   "employee e on b.authorid = e.id join " + 
				   "(select id, max(ver) as ver from bid group by id) as b2 on b.id = b2.id and b.ver=b2.ver " +
				   "where b.tender_id=$1"

		}


		
		rows, err := database.DB.Query(dbQuery, tenderID)
		if err != nil {
			panic(err)
		}

		defer rows.Close()
		
		bids :=[]models.Bid{}
		for rows.Next(){
			b := models.Bid{}
			err := rows.Scan(&b.ID, &b.Name, &b.Status, &b.AuthType, &b.AuthId, &b.Ver, &b.CreatedAt)
			if err != nil{
				log.Println(err)
				continue
			}
			bids = append(bids, b)
		}
		

		if offset > len(bids) {
			c.IndentedJSON(400, "wrong offset value (should be less or equal " + strconv.Itoa(len(bids)) + ")" )
			c.Abort()
			return			
		}
		bids = bids[offset:min(offset+limit,len(bids))]
		
		c.Header("Access-Control-Allow-Origin","*")
		c.IndentedJSON(200, bids)
	}
}

func ShowStatusBid() gin.HandlerFunc{
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin","*")

		bidID := c.Param("Id")		

		b,err :=database.GetBid(bidID)
		if err != nil{
			c.IndentedJSON(400,"Can't get bid with required id")
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
		
		orgResp,err := database.GetTenderResponsible(*b.TenderId)
		if err != nil {
			c.IndentedJSON(500,"Can't get organization responsible")
			c.Abort()
			return
		}
		
		uid, err := database.UserIdByName(quserName[0])
		if err != nil {
			c.IndentedJSON(500,"user_id for passed username")
			c.Abort()
			return
		}
		if uid == *b.AuthId{
			c.IndentedJSON(200,map[string]string{"status":*b.Status})
			return
		}

		 _,ok = orgResp[quserName[0]]
		if ok {
			c.IndentedJSON(200,map[string]string{"status":*b.Status})
			return
		}
				
		c.IndentedJSON(400, "this username don't have access to bid_id")
	}
	
}

func ChangeStatusBid() gin.HandlerFunc{
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin","*")

		bidID := c.Param("Id")		

		b,err :=database.GetBid(bidID)
		if err != nil{
			c.IndentedJSON(400,"Can't get bid with required id")
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
		
		orgResp,err := database.GetTenderResponsible(*b.TenderId)
		if err != nil {
			c.IndentedJSON(500,"Can't get organization responsible")
			c.Abort()
			return
		}
		
		

		 _,ok = orgResp[quserName[0]]
		if !ok {
			c.IndentedJSON(400, "this username don't have access to bid_id")
			c.Abort()
			return
		}

		if *b.Status == "Closed" {
			c.IndentedJSON(400, "can't change status. bid closed")
			c.Abort()
			return
		}

		query := "update bid set status=$1 where id=$2 and ver=$3"
		_, err = database.DB.ExecContext(context.Background(),query, qstatus[0], *b.ID,*b.Ver)
		if err != nil {
			c.IndentedJSON(500, "imposible update status in db table")
			c.Abort()
			return
		}

		b,err =database.GetBid(bidID)
		if err != nil {
			c.IndentedJSON(400,"Can't get bid with required id")
			c.Abort()
			return
		}
		
		c.IndentedJSON(200,map[string]string{"new_status":*b.Status})
		
	}
	
}

func EditBid() gin.HandlerFunc {
	return func(c *gin.Context) {
	
		bidId := c.Param("Id")
		queryparam := c.Request.URL.Query()
		
	
		quserName,ok := queryparam["username"]
		if !ok{
			c.IndentedJSON(400, "username not passed")
			c.Abort()
			return
		}
		qbidAuthorId,err := database.UserIdByName(quserName[0])
		if err != nil{
			c.IndentedJSON(400, "can't get author id for passed username")
			c.Abort()
			return
		}


		b,err :=database.GetBid(bidId)
		if err != nil{
			c.IndentedJSON(400,"Can't get bid with required id")
			c.Abort()
			return
		}

		if *b.AuthId != qbidAuthorId{
			c.IndentedJSON(400, "you are not creator")
			c.Abort()
			return
		}

		if err := c.BindJSON(&b); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		
		// tender.Status 	   = &status
		*b.Ver+=1
		b.CreatedAt,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		query := "INSERT INTO bid(id, ver,name,description,status,tender_id,authortype,authorid,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);"
		_, err = database.DB.ExecContext(context.Background(),query,*b.ID, *b.Ver,*b.Name,*b.Description,*b.Status,*b.TenderId,
		*b.AuthType,*b.AuthId, b.CreatedAt)
		if err != nil {
			c.IndentedJSON(500, "imposible insert tender to db table")
			c.Abort()
			return
		}
		
		c.JSON(http.StatusCreated, b)

	}
}

func RollbackVerBid() gin.HandlerFunc {
	return func(c *gin.Context) {
	
		bidId := c.Param("Id")

		bidVer,err := strconv.Atoi(c.Param("ver"))
		if err != nil{
			c.IndentedJSON(400, "wrong format of bid version")
			c.Abort()
			return
		}
		
		queryparam := c.Request.URL.Query()
		quserName,ok := queryparam["username"]
		if !ok{
			c.IndentedJSON(400, "username not passed")
			c.Abort()
			return
		}
		qbidAuthorId,err := database.UserIdByName(quserName[0])
		if err != nil{
			c.IndentedJSON(400, "can't get author id for passed username")
			c.Abort()
			return
		}

		b,err :=database.GetBidVer(bidId,bidVer)
		if err != nil{
			c.IndentedJSON(400,"Can't get tender with required id and ver")
			c.Abort()
			return
		}

		if *b.AuthId != qbidAuthorId{
			c.IndentedJSON(400, "you are not creator")
			c.Abort()
			return
		}
		lastVer, err := database.GetBidLastVerNum(bidId)
		if err!=nil{
			c.IndentedJSON(500, "cant get last version number from db")
			c.Abort()
			return
		}
		// tender.Status 	   = &status
		*b.Ver=uint(lastVer)+1
		b.CreatedAt,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		query := "INSERT INTO bid(id, ver,name,description,status,tender_id,authortype,authorid,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);"
		_, err = database.DB.ExecContext(context.Background(),query, *b.ID, *b.Ver,*b.Name,*b.Description,*b.Status,*b.TenderId,
		*b.AuthType,*b.AuthId, b.CreatedAt)
		if err != nil {
			c.IndentedJSON(500, "imposible insert bid to db table")
			log.Println(err)
			c.Abort()
			return
		}
		
		c.JSON(http.StatusCreated, b)

	}
}

func BidFeedback() gin.HandlerFunc {
	return func(c *gin.Context) {
		bidId := c.Param("Id")

		queryparam := c.Request.URL.Query()
		qbidFB,ok  := queryparam["bidFeedback"]
		if !ok{
			c.IndentedJSON(400, "feedback not passed")
			c.Abort()
			return
		}

		quserName,ok := queryparam["username"]
		if !ok{
			c.IndentedJSON(400, "username not passed")
			c.Abort()
			return
		}

		b,err :=database.GetBid(bidId)
		if err != nil{
			c.IndentedJSON(400,"Can't get bid with required id")
			c.Abort()
			return
		}

		orgResp,err := database.GetTenderResponsible(*b.TenderId)
		if err != nil {
			c.IndentedJSON(500,"Can't get organization responsible")
			c.Abort()
			return
		}
		
		 _,ok = orgResp[quserName[0]]
		if !ok {
			c.IndentedJSON(400, "this username don't have access to bid_id")
			c.Abort()
			return
		}
		uid,_ := database.UserIdByName(quserName[0])

		query := "INSERT INTO feedback ( bid_id, bid_ver,bidauthor_id, reviewer_id, fb_text) VALUES ($1,$2,$3,$4,$5);"
		_, err = database.DB.ExecContext(context.Background(),query, *b.ID, *b.Ver, *b.AuthId, uid, qbidFB[0])
		if err != nil {
			c.IndentedJSON(500, "imposible insert bid to db table")
			log.Println(err)
			c.Abort()
			return
		}

		c.JSON(http.StatusCreated, "feedbacke message succefully submitted")

	}

}

func BidReviews() gin.HandlerFunc {
	return func(c *gin.Context) {
		bidId := c.Param("Id")

		queryparam := c.Request.URL.Query()
		qAuthor,ok  := queryparam["authorUsername"]
		if !ok{
			c.IndentedJSON(400, "feedback not passed")
			c.Abort()
			return
		}

		qRequester,ok := queryparam["requesterUsername"]
		if !ok{
			c.IndentedJSON(400, "username not passed")
			c.Abort()
			return
		}

		b,err :=database.GetBid(bidId)
		if err != nil{
			c.IndentedJSON(400,"Can't get bid with required id")
			c.Abort()
			return
		}

		orgResp,err := database.GetTenderResponsible(*b.TenderId)
		if err != nil {
			c.IndentedJSON(500,"Can't get organization responsible")
			c.Abort()
			return
		}
		
		 _,ok = orgResp[qRequester[0]]
		if !ok {
			c.IndentedJSON(400, "this username don't have access to bid_id")
			c.Abort()
			return
		}

		//requsterId,_ := database.UserIdByName(qRequester[0])
		authorId,err   := database.UserIdByName(qAuthor[0])
		if err != nil{
			c.IndentedJSON(400, "this author not found")
			c.Abort()
			return
		}
		query := "select bid_id, bid_ver, bidauthor_id, reviewer_id, fb_text from  feedback where bidauthor_id=$1"
		rows, err := database.DB.Query(query, authorId)
			if err != nil {
				panic(err)
		}

		defer rows.Close()
		
		fbs :=[]models.Feedback{}
		for rows.Next(){
			fb := models.Feedback{}
			err := rows.Scan(&fb.BidId, &fb.BidVer, &fb.BidAuthorId, &fb.ReviewerId, &fb.FBtext)
			if err != nil{
				log.Println(err)
				continue
			}
			fbs = append(fbs, fb)
		}
		log.Println(len(fbs))
		c.IndentedJSON(200,fbs)

	}

}