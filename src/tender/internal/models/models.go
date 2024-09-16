package models
import (
	"time"
)
var STypesAllow = map[string]bool{
	"Construction":true,
	"Delivery":true,
	"Manufacture":true,
}

var AuthTypeAllow = map[string]bool{
	"Organization":true,
	"User":true,
}

//var ServiceTypes = []string{"Construction", "Delivery", "Manufacture"}


// func InServiceTypesAllowed(a string) bool {
// 	for _, b := range ServiceTypes {
// 		if b == a {
// 			return true
// 		}
// 	}
// 	return false
// }


type Tender struct{
	ID	            *string	   `json:"uid"             						           	 "`	
	Ver             *uint      `json:"version"`
	Name 		    *string	   `json:"name"                      validate:"required,max=100"`
    Description     *string    `json:"description,omitempty"     validate:"required,max=500"`
	ServiceType		*string    `json:"serviceType"               validate:"required,max=100"`
	Status			*string    `json:"status"                                              "`  // 
	OrgId			*string    `json:"organizationId,omitempty"  validate:"required,max=100"`  //
	CreatedAt       time.Time  `json:"created_at""`  //RFC3339
	CreatorName     *string	   `json:"creatorUsername,omitempty" validate:"required,max=100"`
}

type Bid struct{
	ID	            *string	   `json:"uid"`
	Name 		    *string	   `json:"name"          validate:"required,max=100"`
	Description     *string    `json:"description,omitempty"   validate:"required,max=500"`
	Status			*string    `json:"status"        `  // 
	TenderId        *string    `json:"tenderId,omitempty"     `  //
	AuthType		*string    `json:"authorType"     `  //
	AuthId			*string    `json:"authorId"       `  //
	Ver             *uint      `json:"version"       `
	CreatedAt       time.Time  `json:"created_at""`
}

type Feedback struct{
	ID				*string	   `json:"uid,omitempty"`
	BidId	        *string	   `json:"bid_id"`
	BidVer 	  	    *string	   `json:"bid_ver"`
	BidAuthorId     *string    `json:"bidauthorid"`
	ReviewerId		*string    `json:"revid"        `  // 
	FBtext          *string    `json:"feedbacktext"     `  //
}