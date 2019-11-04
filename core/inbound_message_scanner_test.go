package core

import (
	"fmt"
	"sync"
	"testing"

	"github.com/OpenBazaar/openbazaar-go/pb"
	"github.com/OpenBazaar/openbazaar-go/repo"
	"github.com/OpenBazaar/openbazaar-go/repo/db"
	"github.com/OpenBazaar/openbazaar-go/schema"
	wi "github.com/OpenBazaar/wallet-interface"
	"github.com/op/go-logging"
)

// DISPUTE CASES
func TestPerformTaskInboundMessageScanner(t *testing.T) {
	var (
		orderMsgWithNoErr = repo.OrderMessage{
			MessageID:   "1",
			OrderID:     "1",
			MessageType: int32(pb.Message_ORDER),
			Message:     []byte("sample message"),
			MsgErr:      "",
			PeerID:      "sample",
			PeerPubkey:  []byte("sample"),
		}

		orderMsgWithErr = repo.OrderMessage{
			MessageID:   "2",
			OrderID:     "2",
			MessageType: int32(pb.Message_ORDER),
			Message:     []byte("sample message"),
			MsgErr:      ErrInsufficientFunds.Error(),
			PeerID:      "sample",
			PeerPubkey:  []byte("sample"),
		}

		existingRecords = []repo.OrderMessage{
			orderMsgWithNoErr,
			orderMsgWithErr,
		}

		appSchema = schema.MustNewCustomSchemaManager(schema.SchemaContext{
			DataPath:        schema.GenerateTempPath(),
			TestModeEnabled: true,
		})
	)

	if err := appSchema.BuildSchemaDirectories(); err != nil {
		t.Fatal(err)
	}
	defer appSchema.DestroySchemaDirectories()
	if err := appSchema.InitializeDatabase(); err != nil {
		t.Fatal(err)
	}

	database, err := appSchema.OpenDatabase()
	if err != nil {
		t.Fatal(err)
	}
	s, err := database.Prepare("insert into messages (messageID, orderID, message_type, message, peerID, err, pubkey) values (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		t.Fatal(err)
	}

	for _, r := range existingRecords {
		_, err = s.Exec(r.MessageID, r.OrderID, r.MessageType, r.Message, r.PeerID, r.MsgErr, r.PeerPubkey)
		if err != nil {
			t.Fatal(err)
		}
	}

	datastore := db.NewSQLiteDatastore(database, new(sync.Mutex), wi.Bitcoin)
	worker := &inboundMessageScanner{
		datastore: datastore,
		logger:    logging.MustGetLogger("testInboundMsgScanner"),
	}

	//worker.PerformTask()
	fmt.Println(worker)
	msgs, err := worker.datastore.Messages().GetAllErrored()

	fmt.Println(len(msgs))
	fmt.Println(err)

}
