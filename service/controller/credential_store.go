package controller

import (
	"encoding/json"
	"fmt"
	"time"

	"git.ont.io/ontid/otf/common/message"
)

const (
	CredentialKey        = "Credential"
	RequestCredentialKey = "RequestCredential"
	OfferCredentialKey   = "OfferCredential"
)

func (c *CredentialController) SaveOfferCredential(did, id string, propsal *message.OfferCredential) error {
	key := []byte(fmt.Sprintf("%s_%s_%s", OfferCredentialKey, did, id))
	b, err := c.store.Has(key)
	if err != nil {
		return err
	}
	if b {
		return fmt.Errorf("id:%s already exist\n", id)
	}

	data, err := json.Marshal(propsal)
	if err != nil {
		return err
	}
	return c.store.Put(key, data)
}

func (c *CredentialController) SaveCredential(did, id string, credential message.IssueCredential) error {
	key := []byte(fmt.Sprintf("%s_%s_%s", CredentialKey, did, id))
	b, err := c.store.Has(key)
	if err != nil {
		return err
	}
	if b {
		return fmt.Errorf("id:%s already exist\n", id)
	}

	rec := message.CredentialRec{
		OwnerDID:   credential.Connection.TheirDid,
		Credential: credential,
		Timestamp:  time.Now(),
	}
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return c.store.Put(key, data)
}

func (c *CredentialController) SaveRequestCredential(did, id string, requestCredential message.RequestCredential) error {
	key := []byte(fmt.Sprintf("%s_%s_%s", RequestCredentialKey, did, id))
	b, err := c.store.Has(key)
	if err != nil {
		return err
	}
	if b {
		return fmt.Errorf("id:%s already exist\n", id)
	}

	rec := message.RequestCredentialRec{
		RequesterDID:      requestCredential.Connection.MyDid,
		RequestCredential: requestCredential,
		State:             message.RequestCredentialReceived,
	}
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return c.store.Put(key, data)
}

func (c *CredentialController) QueryCredentialFromStore(did, id string) (message.IssueCredential, error) {
	key := []byte(fmt.Sprintf("%s_%s_%s", CredentialKey, did, id))

	data, err := c.store.Get(key)
	if err != nil {
		return message.IssueCredential{}, err
	}

	rec := new(message.CredentialRec)
	err = json.Unmarshal(data, rec)
	if err != nil {
		return message.IssueCredential{}, err
	}
	return rec.Credential, nil
}

func (c *CredentialController) UpdateRequestCredential(did, id string, state message.RequestCredentialState) error {
	key := []byte(fmt.Sprintf("%s_%s_%s", RequestCredentialKey, did, id))
	data, err := c.store.Get(key)
	if err != nil {
		return err
	}

	rec := new(message.RequestCredentialRec)
	err = json.Unmarshal(data, rec)
	if err != nil {
		return err
	}
	if rec.State >= state {
		return fmt.Errorf("UpdateRequestCredential id :%s state invalid\n", id)
	}
	rec.State = state
	data, err = json.Marshal(rec)
	if err != nil {
		return err
	}
	return c.store.Put(key, data)
}

func (c *CredentialController) DelRequestCredential(did, id string) error {
	key := []byte(fmt.Sprintf("%s_%s_%s", RequestCredentialKey, did, id))
	return c.store.Delete(key)
}

func (c *CredentialController) DelCredential(did, id string) error {
	key := []byte(fmt.Sprintf("%s_%s_%s", CredentialKey, did, id))
	return c.store.Delete(key)
}
