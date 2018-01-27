package reminder

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func Check(c *gin.Context) {
	id := c.PostForm("id")
	config := NewLineConfig()
	target, pErr := config.GetProfile(id)
	if pErr != nil {
		log.Fatal(pErr.Error())
	}

	statusFlag, status, sErr := GetStatus(id)
	if sErr != nil {
		log.Fatal(sErr.Error())
	}

	if !statusFlag {
		mErr := config.PostMessage(target + ": " + os.Getenv("CHECKED_MESSAGE"))
		if mErr != nil {
			log.Println(mErr.Error())
		}
	}

	Response(c, status)
}
