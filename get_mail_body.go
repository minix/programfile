package main

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "gtuiw"
	password = "vvupLNG(+hQh"
	dbname   = "db_test"
)

type DbField struct {
	Date      string
	Subject   string
	Content   string
	Attr_path string
}

func main() {
	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("imap.mail.ru.:993", nil)
	//c, err := client.DialTLS("imap-mail.outlook.com.:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	//if err := c.Login("mgtuiw@outlook.com", "messi123"); err != nil {
	if err := c.Login("mgtuiw@mail.ru", "xxxxxx"); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// List mailboxes
	//mailboxes := make(chan *imap.MailboxInfo, 10)
	//done := make(chan error, 1)
	//go func() {
	//	done <- c.List("", "*", mailboxes)
	//}()

	//log.Println("Mailboxes:")
	//for m := range mailboxes {
	//	log.Println("* " + m.Name)
	//}

	//if err := <-done; err != nil {
	//	log.Fatal(err)
	//}

	// Select INBOX
	//var c *client.Client

	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	if mbox.Messages == 0 {
		log.Fatal("No message in mailbox")
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(mbox.Messages, mbox.Messages)

	//attrs := []string{"BODY[]"}
	attrs := []imap.FetchItem{imap.FetchItem("BODY[]")}

	messages := make(chan *imap.Message, 1)
	go func() {
		if err := c.Fetch(seqset, attrs, messages); err != nil {
			log.Fatal(err)
		}
	}()

	msg := <-messages
	r := msg.GetBody("BODY[]")
	if r == nil {
		log.Fatal("Server didn't returned message body")
	}

	mr, err := mail.CreateReader(r)
	if err != nil {
		log.Fatal(err)
	}

	header := mr.Header
	//if date, err := header.Date(); err == nil {
	//	//		log.Println("Date:", date)
	//	fmt.Println("Date:", date)
	//}
	date, err := header.Date()
	checkErr(err)

	//if from, err := header.AddressList("From"); err == nil {
	//	log.Println("From:", from)
	//}

	//if to, err := header.AddressList("To"); err == nil {
	//	log.Println("To:", to)
	//}

	subject, err := header.Subject()
	checkErr(err)

	//db Connect
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		switch h := p.Header.(type) {
		case mail.TextHeader:
			mailType := p.Header.Get("Content-Type")
			mailContentType := strings.SplitAfter(mailType, ";")
			if mailContentType[0] == "text/plain;" {
				b, _ := ioutil.ReadAll(p.Body)
				//log.Println(subject, date, string(b))
				//fmt.Println(string(b))
				_, err := db.Exec("INSERT INTO mail(date, subject, content,created_on, updated_on) VALUES($1, $2, $3, $4, $5)", date.Format("2006-01-02 15:04:05"), subject, string(b), time.Now(), time.Now())
				checkErr(err)
			}
		case mail.AttachmentHeader:
			filename, _ := h.Filename()
			log.Println("Got attachment: %v", filename)
		}
	}
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
