package castbrick

import (
	"context"
	"net/url"
	"strconv"
)

// ContactsResource handles contacts and contact list operations.
type ContactsResource struct{ c *Client }

func (r *ContactsResource) List(ctx context.Context, page, pageSize int, search string) (*PagedResult[Contact], error) {
	q := url.Values{}
	q.Set("pageNumber", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	if search != "" {
		q.Set("search", search)
	}
	var out PagedResult[Contact]
	if err := r.c.getQuery(ctx, "/audience/contacts", q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *ContactsResource) Get(ctx context.Context, id string) (*Contact, error) {
	path, err := pathf("/audience/contacts/%s", id)
	if err != nil {
		return nil, err
	}
	var out Contact
	if err := r.c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
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
	path, err := pathf("/audience/contacts/%s", id)
	if err != nil {
		return err
	}
	return r.c.delete(ctx, path)
}

func (r *ContactsResource) ListLists(ctx context.Context, page, pageSize int) (*PagedResult[ContactList], error) {
	q := url.Values{}
	q.Set("pageNumber", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	var out PagedResult[ContactList]
	if err := r.c.getQuery(ctx, "/audience/lists", q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *ContactsResource) GetList(ctx context.Context, id string) (*ContactList, error) {
	path, err := pathf("/audience/lists/%s", id)
	if err != nil {
		return nil, err
	}
	var out ContactList
	if err := r.c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateList creates a new contact list and returns the new list ID.
func (r *ContactsResource) CreateList(ctx context.Context, name string) (string, error) {
	var out string
	if err := r.c.post(ctx, "/audience/lists", map[string]any{"name": name}, &out); err != nil {
		return "", err
	}
	return out, nil
}

func (r *ContactsResource) AddToList(ctx context.Context, listID, contactID string) error {
	path, err := pathf("/audience/lists/%s/contacts", listID)
	if err != nil {
		return err
	}
	return r.c.post(ctx, path, map[string]any{"contactId": contactID}, nil)
}

func (r *ContactsResource) RemoveFromList(ctx context.Context, listID, contactID string) error {
	path, err := pathf("/audience/lists/%s/contacts/%s", listID, contactID)
	if err != nil {
		return err
	}
	return r.c.delete(ctx, path)
}
