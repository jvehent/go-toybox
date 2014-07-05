package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Command struct {
	ID            uint64        `json:"id"`
	Action        Action        `json:"action"`
	AgentName     string        `json:"agentname"`
	AgentQueueLoc string        `json:"agentqueueloc"`
	Status        string        `json:"status"`
	Results       []interface{} `json:"results"`
	StartTime     time.Time     `json:"starttime"`
	FinishTime    time.Time     `json:"finishtime"`
}

// a MetaAction is a json object that extends an Action with
// additional parameters. It is used to track the completion
// of an action on agents.
type ExtendedAction struct {
	Action         Action    `json:"action"`
	Status         string    `json:"status"`
	StartTime      time.Time `json:"starttime"`
	FinishTime     time.Time `json:"finishtime"`
	LastUpdateTime time.Time `json:"lastupdatetime"`
	CommandIDs     []uint64  `json:"commandids"`
	Counters       counters  `json:"counters"`
}

// Some counters used to track the completion of an action
type counters struct {
	Sent      int `json:"sent"`
	Returned  int `json:"returned"`
	Done      int `json:"done"`
	Cancelled int `json:"cancelled"`
	Failed    int `json:"failed"`
	TimeOut   int `json:"timeout"`
}

// an Action is the json object that is created by an investigator
// and provided to the MIG platform. It must be PGP signed.
type Action struct {
	ID            uint64      `json:"id"`
	Name          string      `json:"name"`
	Target        string      `json:"target"`
	Description   Description `json:"description"`
	Threat        Threat      `json:"threat"`
	ValidFrom     time.Time   `json:"validfrom"`
	ExpireAfter   time.Time   `json:"expireafter"`
	Operations    []Operation `json:"operations"`
	PGPSignatures []string    `json:"pgpsignatures"`
	SyntaxVersion int         `json:"syntaxversion"`
}

// a description is a simple object that contains detail about the
// action's author, and it's revision.
type Description struct {
	Author   string `json:"author"`
	Email    string `json:"email"`
	URL      string `json:"url"`
	Revision int    `json:"revision"`
}

// a threat provides the investigator with an idea of how dangerous
// a the compromission might be, if the indicators return positive
type Threat struct {
	Level  string `json:"level"`
	Family string `json:"family"`
}

// an operation is an object that map to an agent module.
// the parameters of the operation are passed to the module as argument,
// and thus their format depend on the module itself.
type Operation struct {
	Module     string      `json:"module"`
	Parameters interface{} `json:"parameters"`
}

func main() {
	db, err := sql.Open("postgres", "user=someclient password=someclient dbname=testjson sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the tables

	// 1. Let's retrieve a document and print some fields
	rows, err := db.Query(`SELECT	CAST(document->>'id' AS numeric),
					CAST(document->'action'->>'id' AS numeric),
					document->'action'->>'name' as name,
					document->'action'->'description'->>'author',
					CAST(document->>'starttime' AS timestamp)
				FROM somejson
				WHERE document->>'status' = 'succeeded'
				AND CAST(document->>'starttime' AS timestamp) < NOW()`)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	defer rows.Close()
	if err == sql.ErrNoRows {
		log.Println("No row found")
	}
	var cmdid, aid uint64
	for rows.Next() {
		var name, author string
		var starttime time.Time
		err = rows.Scan(&cmdid, &aid, &name, &author, &starttime)
		if err != nil {
			panic(err)
		}
		fmt.Println("in", aid, "/", cmdid, ",", author, "ran", name, "on", starttime)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	// 2. update the document in the database
	rows, err = db.Query("SELECT document FROM somejson WHERE CAST(document->>'id' AS numeric) = $1 AND CAST(document->'action'->>'id' AS numeric) = $2", cmdid, aid)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var stuff []byte
		err = rows.Scan(&stuff)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", stuff)
		var cmd Command
		err = json.Unmarshal(stuff, &cmd)
		fmt.Println(cmd.ID)

	}

	// 3. write action in database
	var a Action
	err = json.Unmarshal([]byte(somejsonaction), &a)
	aDesc, err := json.Marshal(a.Description)
	if err != nil {
		panic(err)
	}
	aThreat, err := json.Marshal(a.Threat)
	if err != nil {
		panic(err)
	}
	aOperations, err := json.Marshal(a.Operations)
	if err != nil {
		panic(err)
	}
	aPGPSignatures, err := json.Marshal(a.PGPSignatures)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO actions
		(id, name, target, description, threat, operations,
		validfrom, expireafter, starttime, pgpsignatures, syntaxversion)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		a.ID, a.Name, a.Target, aDesc, aThreat, aOperations,
		a.ValidFrom, a.ExpireAfter, time.Now(), aPGPSignatures, a.SyntaxVersion)
	if err != nil {
		panic(err)
	}
}

const somejsonaction string = `{
    "name": "botcode",
    "description": {
        "author": "Julien Vehent",
        "email": "jvehent@mozilla.com",
        "revision": 201402260532
    },
    "target": "linux",
    "threat": {
        "family": "backdoor",
        "level": "alert"
    },
    "operations": [
        {
            "module": "filechecker",
            "parameters": {
                "/bin": {
                    "sha1": {
                        "install.tar": [
                            "71e4602f80d4cb28cc9cc3ce8e91e013636d1f72"
                        ],
                        "b26": [
                            "8a2c86ff5c7583e7ef953a897a705a7b135e8de4"
                        ],
                        "cnet2": [
                            "a617e6fcfbfb55c60287d7066780b34778de3ca4"
                        ],
                        "fake.cfg": [
                            "b888d18d5083b5f558333b5d0fbd0d390228b394"
                        ],
                        "mysql515": [
                            "4d5e1c86e2353e28fd332262c262d0ccf53746df"
                        ],
                        "socket": [
                            "506f8270d6ff38be909a699492c10132c3f7ecfa"
                        ],
                        "taskgrm": [
                            "5c737f0b3858b94d1ccd352f17eca7ebd637b960"
                        ]
                    }
                },
                "/usr": {
                    "sha1": {
                        "install.tar": [
                            "71e4602f80d4cb28cc9cc3ce8e91e013636d1f72"
                        ],
                        "b26": [
                            "8a2c86ff5c7583e7ef953a897a705a7b135e8de4"
                        ],
                        "cnet2": [
                            "a617e6fcfbfb55c60287d7066780b34778de3ca4"
                        ],
                        "fake.cfg": [
                            "b888d18d5083b5f558333b5d0fbd0d390228b394"
                        ],
                        "mysql515": [
                            "4d5e1c86e2353e28fd332262c262d0ccf53746df"
                        ],
                        "socket": [
                            "506f8270d6ff38be909a699492c10132c3f7ecfa"
                        ],
                        "taskgrm": [
                            "5c737f0b3858b94d1ccd352f17eca7ebd637b960"
                        ]
                    }
                },
                "/sbin": {
                    "sha1": {
                        "install.tar": [
                            "71e4602f80d4cb28cc9cc3ce8e91e013636d1f72"
                        ],
                        "b26": [
                            "8a2c86ff5c7583e7ef953a897a705a7b135e8de4"
                        ],
                        "cnet2": [
                            "a617e6fcfbfb55c60287d7066780b34778de3ca4"
                        ],
                        "fake.cfg": [
                            "b888d18d5083b5f558333b5d0fbd0d390228b394"
                        ],
                        "mysql515": [
                            "4d5e1c86e2353e28fd332262c262d0ccf53746df"
                        ],
                        "socket": [
                            "506f8270d6ff38be909a699492c10132c3f7ecfa"
                        ],
                        "taskgrm": [
                            "5c737f0b3858b94d1ccd352f17eca7ebd637b960"
                        ]
                    }
                }
            }
        }
    ],
    "syntaxversion": 1
}`
