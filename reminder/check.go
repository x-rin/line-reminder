package reminder

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func Check(c *gin.Context) {
	id := c.PostForm("id")
	target, pErr := GetProfile(id)
	if pErr != nil {
		log.Fatal(pErr.Error())
	}

	statusFlag, status, sErr := getStatus(id)
	if sErr != nil {
		log.Fatal(sErr.Error())
	}

	if !statusFlag {
		mErr := PostMessage(target + ": " + os.Getenv("CHECKED_MESSAGE"))
		if mErr != nil {
			log.Println(mErr.Error())
		}
	}

	Response(c, status)
}
