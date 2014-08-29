package tokens

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/testhelper"
)

// authTokenPost verifies that providing certain AuthOptions and Scope results in an expected JSON structure.
func authTokenPost(t *testing.T, options gophercloud.AuthOptions, scope *Scope, requestJSON string) {
	setup()
	defer teardown()

	client := gophercloud.ServiceClient{
		Endpoint: endpoint(),
		Options:  options,
		TokenID:  "12345abcdef",
	}

	mux.HandleFunc("/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "Content-Type", "application/json")
		testhelper.TestHeader(t, r, "Accept", "application/json")
		testhelper.TestJSONRequest(t, r, requestJSON)

		fmt.Fprintf(w, `{}`)
	})

	_, err := Create(&client, scope)
	if err != nil {
		t.Errorf("Create returned an error: %v", err)
	}
}

func TestCreateUserIDAndPassword(t *testing.T) {
	authTokenPost(t, gophercloud.AuthOptions{UserID: "me", Password: "squirrel!"}, nil, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": { "id": "me", "password": "squirrel!" }
					}
				}
			}
		}
	`)
}

func TestCreateUsernameDomainIDPassword(t *testing.T) {
	authTokenPost(t, gophercloud.AuthOptions{Username: "fakey", Password: "notpassword", DomainID: "abc123"}, nil, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"domain": {
								"id": "abc123"
							},
							"name": "fakey",
							"password": "notpassword"
						}
					}
				}
			}
		}
	`)
}

func TestCreateUsernameDomainNamePassword(t *testing.T) {
	authTokenPost(t, gophercloud.AuthOptions{Username: "frank", Password: "swordfish", DomainName: "spork.net"}, nil, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"domain": {
								"name": "spork.net"
							},
							"name": "frank",
							"password": "swordfish"
						}
					}
				}
			}
		}
	`)
}

func TestCreateTokenID(t *testing.T) {
	authTokenPost(t, gophercloud.AuthOptions{}, nil, `
		{
			"auth": {
				"identity": {
					"methods": ["token"],
					"token": {
						"id": "12345abcdef"
					}
				}
			}
		}
	`)
}

func TestCreateProjectIDScope(t *testing.T) {
	options := gophercloud.AuthOptions{UserID: "fenris", Password: "g0t0h311"}
	scope := &Scope{ProjectID: "123456"}
	authTokenPost(t, options, scope, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"id": "fenris",
							"password": "g0t0h311"
						}
					}
				},
				"scope": {
					"project": {
						"id": "123456"
					}
				}
			}
		}
	`)
}

func TestCreateDomainIDScope(t *testing.T) {
	options := gophercloud.AuthOptions{UserID: "fenris", Password: "g0t0h311"}
	scope := &Scope{DomainID: "1000"}
	authTokenPost(t, options, scope, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"id": "fenris",
							"password": "g0t0h311"
						}
					}
				},
				"scope": {
					"domain": {
						"id": "1000"
					}
				}
			}
		}
	`)
}

func TestCreateProjectNameAndDomainIDScope(t *testing.T) {
	options := gophercloud.AuthOptions{UserID: "fenris", Password: "g0t0h311"}
	scope := &Scope{ProjectName: "world-domination", DomainID: "1000"}
	authTokenPost(t, options, scope, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"id": "fenris",
							"password": "g0t0h311"
						}
					}
				},
				"scope": {
					"project": {
						"domain": {
							"id": "1000"
						},
						"name": "world-domination"
					}
				}
			}
		}
	`)
}

func TestCreateProjectNameAndDomainNameScope(t *testing.T) {
	options := gophercloud.AuthOptions{UserID: "fenris", Password: "g0t0h311"}
	scope := &Scope{ProjectName: "world-domination", DomainName: "evil-plans"}
	authTokenPost(t, options, scope, `
		{
			"auth": {
				"identity": {
					"methods": ["password"],
					"password": {
						"user": {
							"id": "fenris",
							"password": "g0t0h311"
						}
					}
				},
				"scope": {
					"project": {
						"domain": {
							"name": "evil-plans"
						},
						"name": "world-domination"
					}
				}
			}
		}
	`)
}
