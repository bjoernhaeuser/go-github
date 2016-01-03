// Copyright 2014 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tests

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/go-github/github"
)

func TestUsers_Get(t *testing.T) {
	// list all users
	users, _, err := client.Users.ListAll(nil)
	if err != nil {
		t.Fatalf("Users.ListAll returned error: %v", err)
	}

	if len(users) == 0 {
		t.Errorf("Users.ListAll returned no users")
	}

	// mojombo is user #1
	if want := "mojombo"; want != *users[0].Login {
		t.Errorf("user[0].Login was %q, wanted %q", *users[0].Login, want)
	}

	// get individual user
	u, _, err := client.Users.Get("octocat")
	if err != nil {
		t.Fatalf("Users.Get('octocat') returned error: %v", err)
	}

	if want := "octocat"; want != *u.Login {
		t.Errorf("user.Login was %q, wanted %q", *u.Login, want)
	}
	if want := "The Octocat"; want != *u.Name {
		t.Errorf("user.Name was %q, wanted %q", *u.Name, want)
	}
}

func TestUsers_Update(t *testing.T) {
	if !checkAuth("TestUsers_Get") {
		return
	}

	u, _, err := client.Users.Get("")
	if err != nil {
		t.Fatalf("Users.Get('') returned error: %v", err)
	}

	if *u.Login == "" {
		t.Errorf("wanted non-empty values for user.Login")
	}

	// save original location
	var location string
	if u.Location != nil {
		location = *u.Location
	}

	// update location to test value
	testLoc := fmt.Sprintf("test-%d", rand.Int())
	u.Location = &testLoc

	_, _, err = client.Users.Edit(u)
	if err != nil {
		t.Fatalf("Users.Update returned error: %v", err)
	}

	// refetch user and check location value
	u, _, err = client.Users.Get("")
	if err != nil {
		t.Fatalf("Users.Get('') returned error: %v", err)
	}

	if testLoc != *u.Location {
		t.Errorf("Users.Get('') has location: %v, want: %v", *u.Location, testLoc)
	}

	// set location back to the original value
	u.Location = &location
	_, _, err = client.Users.Edit(u)
}

func TestUsers_Emails(t *testing.T) {
	if !checkAuth("TestUsers_Emails") {
		return
	}

	emails, _, err := client.Users.ListEmails(nil)
	if err != nil {
		t.Fatalf("Users.ListEmails() returned error: %v", err)
	}

	// create random address not currently in user's emails
	var email string
EmailLoop:
	for {
		email = fmt.Sprintf("test-%d@example.com", rand.Int())
		for _, e := range emails {
			if e.Email != nil && *e.Email == email {
				continue EmailLoop
			}
		}
		break
	}

	// Add new address
	_, _, err = client.Users.AddEmails([]string{email})
	if err != nil {
		t.Fatalf("Users.AddEmails() returned error: %v", err)
	}

	// List emails again and verify new email is present
	emails, _, err = client.Users.ListEmails(nil)
	if err != nil {
		t.Fatalf("Users.ListEmails() returned error: %v", err)
	}

	var found bool
	for _, e := range emails {
		if e.Email != nil && *e.Email == email {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("Users.ListEmails() does not contain new addres: %v", email)
	}

	// Remove new address
	_, err = client.Users.DeleteEmails([]string{email})
	if err != nil {
		t.Fatalf("Users.DeleteEmails() returned error: %v", err)
	}

	// List emails again and verify new email was removed
	emails, _, err = client.Users.ListEmails(nil)
	if err != nil {
		t.Fatalf("Users.ListEmails() returned error: %v", err)
	}

	for _, e := range emails {
		if e.Email != nil && *e.Email == email {
			t.Fatalf("Users.ListEmails() still contains address %v after removing it", email)
		}
	}
}

func TestUsers_Keys(t *testing.T) {
	keys, _, err := client.Users.ListKeys("willnorris", nil)
	if err != nil {
		t.Fatalf("Users.ListKeys('willnorris') returned error: %v", err)
	}

	if len(keys) == 0 {
		t.Errorf("Users.ListKeys('willnorris') returned no keys")
	}

	// the rest of the tests requires auth
	if !checkAuth("TestUsers_Keys") {
		return
	}

	keys, _, err = client.Users.ListKeys("", nil)
	if err != nil {
		t.Fatalf("Users.ListKeys('') returned error: %v", err)
	}

	// ssh public key for testing (fingerprint: 16:c4:b0:2b:99:b4:eb:6d:06:a3:7d:7f:d6:55:e1:65)
	key := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDZ4khAJykHnJXKbv2mWgegSVyVKG6mAl6MwlMjGP9CsuBsbgGrMxCVVH1lhRIWNyGNPh3tqt5qmz0rSjaRqBO9FbTxTXbokAbktWgLvSH3PqSKXyS8NmsaMk/PfAm9BOwYpVQ6P2dg75MmKbXgq1xxe3PxEmgdEHzg4XSw+dkoO4Bika88BksFaXHsE1DIlqfDNtTIheq6vccJ2L9bEoXdhsEbe/OaZsyZHeaOg2Z7Z3K531kLfAYH/uD3kUOcA+4sX/tW29aYpVvca7v1sHFjz9RotF3PkxRH/zekLy8/WHxfxwOptBd0TpEiI0Vg6Yy+pNjkSEGwKNPO8Z2aEflT"
	for _, k := range keys {
		if k.Key != nil && *k.Key == key {
			t.Fatalf("Test key already exists for user.  Please manually remove it first.")
		}
	}

	// Add new key
	_, _, err = client.Users.CreateKey(&github.Key{
		Title: github.String("go-github test key"),
		Key:   github.String(key),
	})
	if err != nil {
		t.Fatalf("Users.CreateKey() returned error: %v", err)
	}

	// List keys again and verify new key is present
	keys, _, err = client.Users.ListKeys("", nil)
	if err != nil {
		t.Fatalf("Users.ListKeys('') returned error: %v", err)
	}

	var id int
	for _, k := range keys {
		if k.Key != nil && *k.Key == key {
			id = *k.ID
			break
		}
	}

	if id == 0 {
		t.Fatalf("Users.ListKeys('') does not contain added test key")
	}

	// Verify that fetching individual key works
	k, _, err := client.Users.GetKey(id)
	if err != nil {
		t.Fatalf("Users.GetKey(%q) returned error: %v", id, err)
	}
	if *k.Key != key {
		t.Fatalf("Users.GetKey(%q) returned key %v, want %v", id, *k.Key, key)
	}

	// Remove test key
	_, err = client.Users.DeleteKey(id)
	if err != nil {
		t.Fatalf("Users.DeleteKey(%d) returned error: %v", id, err)
	}

	// List keys again and verify test key was removed
	keys, _, err = client.Users.ListKeys("", nil)
	if err != nil {
		t.Fatalf("Users.ListKeys('') returned error: %v", err)
	}

	for _, k := range keys {
		if k.Key != nil && *k.Key == key {
			t.Fatalf("Users.ListKeys('') still contains test key after removing it")
		}
	}
}
