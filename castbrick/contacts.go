package castbrick

import (
	"context"
	"fmt"
	"net/url"
)

// ContactsResource handles contacts and contact list operations.
type ContactsResource struct{ c *Client }

func (r *ContactsResource) List(ctx context.Context, page, pageSize int, search string) (*PagedResult[Contact], error) {
	path := fmt.Sprintf("/audience/contacts?pageNumber=%d&pageSize=%d", page, pageSize)
	if search != "" {
		path += "&search=" + url.QueryEscape(search)
	}
	var out PagedResult[Contact]
	return &out, r.c.get(ctx, path, &out)
}

func (r *ContactsResource) Get(ctx context.Context, id string) (*Contact, error) {
	var out Contact
	return &out, r.c.get(ctx, "/audience/contacts/"+id, &out)
}

// Create adds one or more contacts from comma/newline-separated phone numbers.
func (r *ContactsResource) Create(ctx context.Context, phoneNumbers string) error {
	body := map[string]any{}
	if phoneNumbers != "" {
		body["phoneNumbers"] = phoneNumbers
	}
	return r.c.post(ctx, "/audience/contacts", body, nil)
}

func (r *ContactsResource) Delete(ctx context.Context, id string) error {
	return r.c.delete(ctx, "/audience/contacts/"+id)
}

func (r *ContactsResource) ListLists(ctx context.Context, page, pageSize int) (*PagedResult[ContactList], error) {
	var out PagedResult[ContactList]
	return &out, r.c.get(ctx, fmt.Sprintf("/audience/lists?pageNumber=%d&pageSize=%d", page, pageSize), &out)
}

func (r *ContactsResource) GetList(ctx context.Context, id string) (*ContactList, error) {
	var out ContactList
	return &out, r.c.get(ctx, "/audience/lists/"+id, &out)
}

// CreateList creates a new contact list and returns the new list ID.
func (r *ContactsResource) CreateList(ctx context.Context, name string) (string, error) {
	var out string
	return out, r.c.post(ctx, "/audience/lists", map[string]any{"name": name}, &out)
}

func (r *ContactsResource) AddToList(ctx context.Context, listID, contactID string) error {
	return r.c.post(ctx, "/audience/lists/"+listID+"/contacts", map[string]any{"contactId": contactID}, nil)
}

func (r *ContactsResource) RemoveFromList(ctx context.Context, listID, contactID string) error {
	return r.c.delete(ctx, "/audience/lists/"+listID+"/contacts/"+contactID)
}
