package bus

import (
	"time"

	"github.com/numary/ledger/pkg/core"
)

type Payload interface {
	PayloadType() string
}

type Event[T Payload] struct {
	Date    time.Time `json:"date"`
	Type    string    `json:"type"`
	Payload T         `json:"payload"`
	Ledger  string    `json:"ledger"`
}

func NewEvent[T Payload](ledger string, payload T) Event[T] {
	return Event[T]{
		Date:    time.Now(),
		Type:    payload.PayloadType(),
		Payload: payload,
		Ledger:  ledger,
	}
}

type CommittedTransactions struct {
	Transactions []core.Transaction `json:"transactions"`
	// Deprecated (use postCommitVolumes)
	Volumes           core.AccountsAssetsVolumes `json:"volumes"`
	PostCommitVolumes core.AccountsAssetsVolumes `json:"postCommitVolumes"`
	PreCommitVolumes  core.AccountsAssetsVolumes `json:"preCommitVolumes"`
}

func (CommittedTransactions) PayloadType() string {
	return CommittedTransactionsLabel
}

type SavedMetadata struct {
	TargetType string        `json:"targetType"`
	TargetID   string        `json:"targetId"`
	Metadata   core.Metadata `json:"metadata"`
}

func (SavedMetadata) PayloadType() string {
	return SavedMetadataLabel
}

type RevertedTransaction struct {
	RevertedTransaction core.Transaction `json:"revertedTransaction"`
	RevertTransaction   core.Transaction `json:"revertTransaction"`
}

func (RevertedTransaction) PayloadType() string {
	return RevertedTransactionLabel
}

type UpdatedMapping struct {
	Mapping core.Mapping `json:"mapping"`
}

func (UpdatedMapping) PayloadType() string {
	return UpdatedMappingLabel
}
