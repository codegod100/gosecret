package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
)

type SecretService struct {
	conn *dbus.Conn
}

type SecretItem struct {
	Path       dbus.ObjectPath
	Label      string
	Attributes map[string]string
	Secret     string
	Created    time.Time
	Modified   time.Time
}

func NewSecretService() (*SecretService, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to session bus: %v", err)
	}

	return &SecretService{conn: conn}, nil
}

func (s *SecretService) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *SecretService) SetSecret(key, secret string) error {
	// First, open a session with the secret service
	serviceObj := s.conn.Object(DBUS_SERVICE_NAME, DBUS_PATH_SECRETS)
	sessionCall := serviceObj.Call("org.freedesktop.Secret.Service.OpenSession", 0, "plain", dbus.MakeVariant(""))
	if sessionCall.Err != nil {
		return fmt.Errorf("failed to open session: %v", sessionCall.Err)
	}

	var output dbus.Variant
	var sessionPath dbus.ObjectPath
	err := sessionCall.Store(&output, &sessionPath)
	if err != nil {
		return fmt.Errorf("failed to parse session response: %v", err)
	}

	// Get the default collection
	collectionPath := dbus.ObjectPath("/org/freedesktop/secrets/aliases/default")
	collectionObj := s.conn.Object(DBUS_SERVICE_NAME, collectionPath)

	// Use key as both label and as the primary attribute
	attributes := map[string]string{
		"gosecret-key": key,
		"application": "gosecret",
	}

	properties := map[string]dbus.Variant{
		"org.freedesktop.Secret.Item.Label":      dbus.MakeVariant(key),
		"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(attributes),
	}

	// Create the secret structure with proper format (oayays)
	secretStruct := struct {
		Session     dbus.ObjectPath
		Parameters  []byte
		Value       []byte
		ContentType string
	}{
		Session:     sessionPath,
		Parameters:  []byte{},
		Value:       []byte(secret),
		ContentType: "text/plain",
	}

	call := collectionObj.Call("org.freedesktop.Secret.Collection.CreateItem", 0, properties, secretStruct, true)
	if call.Err != nil {
		return fmt.Errorf("failed to store secret: %v", call.Err)
	}

	return nil
}

func (s *SecretService) GetSecret(key string) (string, error) {
	obj := s.conn.Object(DBUS_SERVICE_NAME, DBUS_PATH_SECRETS)
	
	// Search for items with our specific key
	attributes := map[string]string{
		"gosecret-key": key,
		"application": "gosecret",
	}

	call := obj.Call("org.freedesktop.Secret.Service.SearchItems", 0, attributes)
	if call.Err != nil {
		return "", fmt.Errorf("failed to search items: %v", call.Err)
	}

	var unlocked, locked []dbus.ObjectPath
	err := call.Store(&unlocked, &locked)
	if err != nil {
		return "", fmt.Errorf("failed to parse search results: %v", err)
	}

	if len(unlocked) == 0 && len(locked) == 0 {
		return "", nil
	}

	var itemPath dbus.ObjectPath
	if len(unlocked) > 0 {
		itemPath = unlocked[0]
	} else {
		itemPath = locked[0]
	}

	itemObj := s.conn.Object(DBUS_SERVICE_NAME, itemPath)
	
	sessionCall := obj.Call("org.freedesktop.Secret.Service.OpenSession", 0, "plain", dbus.MakeVariant(""))
	if sessionCall.Err != nil {
		return "", fmt.Errorf("failed to open session: %v", sessionCall.Err)
	}

	var output dbus.Variant
	var sessionPath dbus.ObjectPath
	err = sessionCall.Store(&output, &sessionPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse session response: %v", err)
	}

	secretCall := itemObj.Call("org.freedesktop.Secret.Item.GetSecret", 0, sessionPath)
	if secretCall.Err != nil {
		return "", fmt.Errorf("failed to get secret: %v", secretCall.Err)
	}

	var secretStruct []interface{}
	err = secretCall.Store(&secretStruct)
	if err != nil {
		return "", fmt.Errorf("failed to parse secret response: %v", err)
	}

	if len(secretStruct) >= 3 {
		if secretData, ok := secretStruct[2].([]byte); ok {
			return string(secretData), nil
		}
	}

	return "", fmt.Errorf("unexpected secret format")
}

func (s *SecretService) DeleteSecret(key string) error {
	obj := s.conn.Object(DBUS_SERVICE_NAME, DBUS_PATH_SECRETS)
	
	// Search for items with our specific key
	attributes := map[string]string{
		"gosecret-key": key,
		"application": "gosecret",
	}

	call := obj.Call("org.freedesktop.Secret.Service.SearchItems", 0, attributes)
	if call.Err != nil {
		return fmt.Errorf("failed to search items: %v", call.Err)
	}

	var unlocked, locked []dbus.ObjectPath
	err := call.Store(&unlocked, &locked)
	if err != nil {
		return fmt.Errorf("failed to parse search results: %v", err)
	}

	allItems := append(unlocked, locked...)
	if len(allItems) == 0 {
		return fmt.Errorf("no secret found with key: %s", key)
	}

	for _, itemPath := range allItems {
		itemObj := s.conn.Object(DBUS_SERVICE_NAME, itemPath)
		deleteCall := itemObj.Call("org.freedesktop.Secret.Item.Delete", 0)
		if deleteCall.Err != nil {
			return fmt.Errorf("failed to delete item: %v", deleteCall.Err)
		}
	}

	return nil
}

func (s *SecretService) ListSecrets(pattern string) error {
	obj := s.conn.Object(DBUS_SERVICE_NAME, DBUS_PATH_SECRETS)
	
	// Search for all gosecret items
	attributes := map[string]string{
		"application": "gosecret",
	}

	call := obj.Call("org.freedesktop.Secret.Service.SearchItems", 0, attributes)
	if call.Err != nil {
		return fmt.Errorf("failed to search items: %v", call.Err)
	}

	var unlocked, locked []dbus.ObjectPath
	err := call.Store(&unlocked, &locked)
	if err != nil {
		return fmt.Errorf("failed to parse search results: %v", err)
	}

	allItems := append(unlocked, locked...)
	if len(allItems) == 0 {
		fmt.Println("No secrets found.")
		return nil
	}

	sessionCall := obj.Call("org.freedesktop.Secret.Service.OpenSession", 0, "plain", dbus.MakeVariant(""))
	if sessionCall.Err != nil {
		return fmt.Errorf("failed to open session: %v", sessionCall.Err)
	}

	var output dbus.Variant
	var sessionPath dbus.ObjectPath
	err = sessionCall.Store(&output, &sessionPath)
	if err != nil {
		return fmt.Errorf("failed to parse session response: %v", err)
	}

	for _, itemPath := range allItems {
		item, err := s.getItemDetails(itemPath, sessionPath)
		if err != nil {
			continue
		}

		// Apply pattern filtering if specified
		if pattern != "" && !strings.Contains(item.Label, pattern) {
			continue
		}

		s.printSimpleItem(item)
	}

	return nil
}

func (s *SecretService) getItemDetails(itemPath dbus.ObjectPath, sessionPath dbus.ObjectPath) (*SecretItem, error) {
	itemObj := s.conn.Object(DBUS_SERVICE_NAME, itemPath)

	labelVar, err := itemObj.GetProperty("org.freedesktop.Secret.Item.Label")
	if err != nil {
		return nil, err
	}
	label := labelVar.Value().(string)

	attributesVar, err := itemObj.GetProperty("org.freedesktop.Secret.Item.Attributes")
	if err != nil {
		return nil, err
	}
	attributes := attributesVar.Value().(map[string]string)

	createdVar, err := itemObj.GetProperty("org.freedesktop.Secret.Item.Created")
	if err != nil {
		return nil, err
	}
	created := time.Unix(int64(createdVar.Value().(uint64)), 0)

	modifiedVar, err := itemObj.GetProperty("org.freedesktop.Secret.Item.Modified")
	if err != nil {
		return nil, err
	}
	modified := time.Unix(int64(modifiedVar.Value().(uint64)), 0)

	secretCall := itemObj.Call("org.freedesktop.Secret.Item.GetSecret", 0, sessionPath)
	var secret string
	if secretCall.Err == nil {
		var secretStruct []interface{}
		if secretCall.Store(&secretStruct) == nil && len(secretStruct) >= 3 {
			if secretData, ok := secretStruct[2].([]byte); ok {
				secret = string(secretData)
			}
		}
	}

	return &SecretItem{
		Path:       itemPath,
		Label:      label,
		Attributes: attributes,
		Secret:     secret,
		Created:    created,
		Modified:   modified,
	}, nil
}

func (s *SecretService) printSimpleItem(item *SecretItem) {
	fmt.Printf("%-30s %s\n", item.Label, item.Created.Format("2006-01-02 15:04:05"))
}

